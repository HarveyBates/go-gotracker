import fitdecode
import datetime
import pandas as pd
import json
from dotenv import load_dotenv
import os
from influxdb_client import InfluxDBClient, Point, WriteOptions
from influxdb_client.client.write_api import SYNCHRONOUS
import psycopg2
import psycopg2.extras
import sys
import re

load_dotenv()
INFLUXDB_TOKEN = os.getenv("INFLUXDB_TOKEN")

recordFields = []
lapFields = []
sessionFields = []
lengthFields = []
deviceFields = []


def parse_fit_file(conn, fName):
    records = []
    laps = []
    lengths = []
    deviceInfo = []
    with fitdecode.FitReader(f"activities/{fName}") as ff:
        for row in ff:
            if isinstance(row, fitdecode.records.FitDataMessage):
                if row.name == "session":
                    session = parse_row(row, row.name)
                elif row.name == "record":
                    records.append(parse_row(row, row.name))
                elif row.name == "lap":
                    laps.append(parse_row(row, row.name))
                elif row.name == "length":
                    lengths.append(parse_row(row, row.name))
                elif row.name == "device_info":
                    deviceInfo.append(parse_row(row, row.name))

    #Postgres
    # - Session

    #Influx:
    # - Record
    # - Laps
    # - Length
    # - Device Information

    sport_type = "unknown"
    sub_sport = "unknown"

    # Session - Information at the end of an activity
    session_df = pd.DataFrame([session], columns=sessionFields).dropna(axis=1, how="all")
    if not session_df.empty:
        #print(session_df)
        if "sport" in session_df:
            sport_type = session_df["sport"].values[0]
            #print(sport_type)
        if "sub_sport" in session_df:
            sub_sport = session_df["sub_sport"].values[0]
            #print(sub_sport)

    activity_id = fName[:-4]
    activity_name = f"{sport_type}_{sub_sport}_{activity_id}"
    print(activity_name)

    if not session_df.empty:
        write_df_to_postgres(conn, session_df, sport_type, activity_name, activity_id)
    else:
        print("Session information was empty")


    # Time series data covering the entire activity
    records_df = pd.DataFrame(records, columns=recordFields).dropna(axis=1, how="all")
    if "timestamp" in records_df:
        records_df["timestamp"] = pd.to_datetime(records_df["timestamp"], format="%Y-%m-%dT%H:%M:%SZ")
        records_df.set_index("timestamp", inplace=True)
    if not records_df.empty:
        #print(records_df)
        # Write into records bucket
        write_df_to_influxdb(records_df, "records", activity_name)

    # Splits for each lap inluding rest time
    lap_df = pd.DataFrame(laps, columns=lapFields).dropna(axis=1, how="all")
    if "timestamp" in lap_df:
        lap_df["timestamp"] = pd.to_datetime(lap_df["timestamp"], format="%Y-%m-%dT%H:%M:%SZ")
        lap_df.set_index("timestamp", inplace=True)
    if not lap_df.empty:
        write_df_to_influxdb(lap_df, "laps", activity_name)
        #print(lap_df)

    # Length decribes splits for each length in the pool (e.g. 25 meter splits)
    length_df = pd.DataFrame(lengths, columns=lengthFields).dropna(axis=1, how="all")
    if "timestamp" in length_df:
        length_df["timestamp"] = pd.to_datetime(length_df["timestamp"], format="%Y-%m-%dT%H:%M:%SZ")
        length_df.set_index("timestamp", inplace=True)
    if not length_df.empty:
        write_df_to_influxdb(length_df, "lengths", activity_name)
        #print(length_df)

    # Device information e.g. battery over an activity 
    device_df = pd.DataFrame(deviceInfo, columns=deviceFields).dropna(axis=1, how="all")
    if "timestamp" in device_df:
        device_df["timestamp"] = pd.to_datetime(device_df["timestamp"], format="%Y-%m-%dT%H:%M:%SZ")
        device_df.set_index("timestamp", inplace=True)
    if not device_df.empty:
        write_df_to_influxdb(device_df, "device", activity_name)
        #print(device_df)

def connect_to_postgres():
    return psycopg2.connect(user="postgres",
                            password="admin",
                            host="localhost",
                            port="27222",
                            dbname="gogotracker")

     
def write_df_to_postgres(conn, dataframe, sport_type, activity_name, activity_id):

    # Add activity ID to dataframe
    dataframe["activity_id"] = activity_id 
    activity_id_col = dataframe.pop("activity_id")
    dataframe.insert(0, "activity_id", activity_id_col)

    # Add activity name to dataframe
    dataframe["activity_name"] = activity_name
    activity_col = dataframe.pop("activity_name")
    dataframe.insert(1, "activity_name", activity_col)

    # Calculate end_time and add to dataframe
    start_time = datetime.datetime.strptime(dataframe["start_time"][0], "%Y-%m-%dT%H:%M:%SZ")
    elapsed_time = datetime.timedelta(seconds=dataframe["total_elapsed_time"][0])
    end_time = (start_time + elapsed_time).strftime("%Y-%m-%dT%H:%M:%SZ")
    dataframe["end_time"] = end_time
    end_time_col = dataframe.pop("end_time")
    dataframe.insert(4, "end_time", end_time_col)

    # Check if table exists
    check_table_exists = f"CREATE TABLE IF NOT EXISTS {sport_type}_session()"
    cursor = conn.cursor()
    cursor.execute(check_table_exists)
    conn.commit()

    # Check if column in table
    df_column_names = list(dataframe)
    # Create columns with corresponding types
    for col, val in zip(df_column_names, dataframe.values[0]):
        if isinstance(val, int) or col == "activity_id":
            check_col_exists = f"ALTER TABLE {sport_type}_session ADD COLUMN IF NOT EXISTS {col} INTEGER"
            if col == "activity_id":
                # Use this to check if item already exists
                check_col_exists = f"ALTER TABLE {sport_type}_session ADD COLUMN IF NOT EXISTS {col} BIGINT UNIQUE"
                conn.commit()
        elif isinstance(val, float):
            check_col_exists = f"ALTER TABLE {sport_type}_session ADD COLUMN IF NOT EXISTS {col} REAL"
        else:
            check_col_exists = f"ALTER TABLE {sport_type}_session ADD COLUMN IF NOT EXISTS {col} TEXT"
        cursor.execute(check_col_exists)
        conn.commit()
        
    # Add data to column
    column_names = ",".join(df_column_names)
    values = "{}".format(",".join(["%s" for _ in df_column_names]))
    insert_stmt = f"INSERT INTO {sport_type}_session ({column_names}) VALUES ({values}) ON CONFLICT (activity_id) DO NOTHING"
    psycopg2.extras.execute_batch(cursor, insert_stmt, dataframe.values)
    conn.commit()
    cursor.close()


def write_df_to_influxdb(dataframe, bucket, activity_name):
    with InfluxDBClient(url="http://localhost:8086", token=INFLUXDB_TOKEN, org="user") as _client:
        with _client.write_api(write_options=WriteOptions(batch_size=500,
                                                            flush_interval=10_000,
                                                            jitter_interval=2_000,
                                                            retry_interval=5_000,
                                                            max_retries=5,
                                                            max_retry_delay=30_000,
                                                            exponential_base=2)) as _write_client:
            _write_client.write(bucket, "user", record=dataframe, data_frame_measurement_name=activity_name)
    


def parse_row(row, rowType):
    parsedRow = {}
    for field in row.fields:
        if "unknown" in field.name:
            # Skip unknown rows 
            pass
        else:
            if rowType == "record" and field.name not in recordFields:
                recordFields.append(field.name)
            elif rowType == "session" and field.name not in sessionFields:
                sessionFields.append(field.name)
            elif rowType == "length" and field.name not in lengthFields:
                lengthFields.append(field.name)
            elif rowType == "lap" and field.name not in lapFields:
                lapFields.append(field.name)
            elif rowType == "device_info" and field.name not in deviceFields:
                deviceFields.append(field.name)

            try:
                # Remove tuples with all None values and unknown data names
                if field.value != None and all(field.value):
                    parsedRow.update({field.name: field.value})

            except TypeError:
                # all() doesn't work for datetime values
                if field.value != None:
                    if isinstance(field.value, datetime.datetime):
                        # Convert datetime to string
                        parsedRow.update({field.name: field.value.strftime("%Y-%m-%dT%H:%M:%SZ")})

                    elif "lat" in field.name or "long" in field.name:
                        if field.value == 0:
                            # Indoor activity or no position data
                            pass
                        else:
                            convertPos = field.value / ((2**32) / 360)
                            parsedRow.update({field.name: convertPos})
                    else:
                        parsedRow.update({field.name: field.value})

    return parsedRow 



def main(): 
    #parse_fit_file("bike-outdoors.fit")
    conn = connect_to_postgres()
    for file in os.listdir("activities"):
        if ".fit" in file:
            parse_fit_file(conn, file)
    #parse_fit_file("swim-ocean.fit")
    #write_db_to_influxdb()


if __name__ == "__main__":
    main()
