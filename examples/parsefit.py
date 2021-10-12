import fitdecode
import datetime
import pandas as pd
import json
from dotenv import load_dotenv
import os
from influxdb_client import InfluxDBClient, Point, WriteOptions
from influxdb_client.client.write_api import SYNCHRONOUS

load_dotenv()

INFLUXDB_TOKEN = os.getenv("INFLUXDB_TOKEN")

recordFields = []
lapFields = []
sessionFields = []
lengthFields = []
deviceFields = []


def parse_fit_file(fName):
    records = []
    laps = []
    lengths = []
    deviceInfo = []
    with fitdecode.FitReader(fName) as ff:
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

    records_df = pd.DataFrame(records, columns=recordFields).dropna(axis=1, how="all")
    records_df["timestamp"] = pd.to_datetime(records_df["timestamp"], format="%Y-%m-%d %H:%M:%S%z")
    records_df.set_index("timestamp", inplace=True)

    print(records_df)

    return records_df


def write_to_influxdb(dataframe):
    with InfluxDBClient(url="http://localhost:8086", token=INFLUXDB_TOKEN, org="user") as _client:
        with _client.write_api(write_options=WriteOptions(batch_size=500,
                                                            flush_interval=10_000,
                                                            jitter_interval=2_000,
                                                            retry_interval=5_000,
                                                            max_retries=5,
                                                            max_retry_delay=30_000,
                                                            exponential_base=2)) as _write_client:
            _write_client.write("activities", "user", record=dataframe, data_frame_measurement_name="bike-outdoors")
    

def get_from_influx():
    q = '''
        from(bucket: "activities") 
            |> range(start: time(v: "2021-09-24T22:45:25Z"), stop: time(v: "2021-09-25T02:17:58Z")) 
            |> filter(fn: (r) => r["_measurement"] == "bike-outdoors") 
            |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
            |> sort(columns: ["_time"], desc: false)
    '''
    with InfluxDBClient(url="http://localhost:8086", token=INFLUXDB_TOKEN, org="user") as _client:
        query_api = _client.query_api()
        tables = query_api.query_data_frame(q)
        df = pd.DataFrame(tables)
        print(df.head())
        print(df.tail())




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
    df = parse_fit_file("bike-outdoors.fit")
    write_to_influxdb(df)

    #get_from_influx()


if __name__ == "__main__":
    main()
