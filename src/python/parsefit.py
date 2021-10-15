import fitdecode
import datetime
import pandas as pd
import json
from dotenv import load_dotenv
import os
from influxdb_client import InfluxDBClient, Point, WriteOptions
from influxdb_client.client.write_api import SYNCHRONOUS
import psycopg2
import sys

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

    file_id = fName[:-4]
    activity_name = f"{sport_type}_{sub_sport}_{file_id}"
    print(activity_name)

    if not session_df.empty:
        write_df_to_postgres(conn, session_df, sport_type, activity_name)
        sys.exit(1)


    # Time series data covering the entire activity
    records_df = pd.DataFrame(records, columns=recordFields).dropna(axis=1, how="all")
    if "timestamp" in records_df:
        records_df["timestamp"] = pd.to_datetime(records_df["timestamp"], format="%Y-%m-%d %H:%M:%S%z")
        records_df.set_index("timestamp", inplace=True)
    if not records_df.empty:
        #print(records_df)
        # Write into records bucket
        write_df_to_influxdb(records_df, "records", activity_name)

    # Splits for each lap inluding rest time
    lap_df = pd.DataFrame(laps, columns=lapFields).dropna(axis=1, how="all")
    if "timestamp" in lap_df:
        lap_df["timestamp"] = pd.to_datetime(lap_df["timestamp"], format="%Y-%m-%d %H:%M:%S%z")
        lap_df.set_index("timestamp", inplace=True)
    if not lap_df.empty:
        write_df_to_influxdb(lap_df, "laps", activity_name)
        #print(lap_df)

    # Length decribes splits for each length in the pool (e.g. 25 meter splits)
    length_df = pd.DataFrame(lengths, columns=lengthFields).dropna(axis=1, how="all")
    if "timestamp" in length_df:
        length_df["timestamp"] = pd.to_datetime(length_df["timestamp"], format="%Y-%m-%d %H:%M:%S%z")
        length_df.set_index("timestamp", inplace=True)
    if not length_df.empty:
        write_df_to_influxdb(length_df, "lengths", activity_name)
        #print(length_df)

    # Device information e.g. battery over an activity 
    device_df = pd.DataFrame(deviceInfo, columns=deviceFields).dropna(axis=1, how="all")
    if "timestamp" in device_df:
        device_df["timestamp"] = pd.to_datetime(device_df["timestamp"], format="%Y-%m-%d %H:%M:%S%z")
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

     
def write_df_to_postgres(conn, dataframe, sport_type, activity_name):
    # Check if table exists
    check_table_exists = f"CREATE TABLE IF NOT EXISTS {sport_type}_session()"
    cursor = conn.cursor()
    cursor.execute(check_table_exists)
    conn.commit()

    # Check if column in table
    # SOMETHING LIKE ALTER TABLE...

    # Add data to column
    df_column_names = list(dataframe)
    df_column_names.insert(0, "activity_name")
    column_names = ",".join(df_column_names)
    values = "VALUES({})".format(",".join(["%s" for _ in df_column_names]))
    insert_stmt = f"INSERT INTO {sport_type}_session ({column_names}) {values}"
    print(insert_stmt)


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
                        parsedRow.update({field.name: str(field.value)})
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
