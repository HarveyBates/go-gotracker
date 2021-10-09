import fitdecode
import datetime
import matplotlib.pyplot as plt
import json

def swim(fitFrame, indoors=True):
    if not indoors:
        # position_lat, position_long, lap, distance, speed, cadence 
        POSITION_FIELDS = ["latitude", "longitude", "lap_n", "distance", "speed", "stroke_rate"]

        for field in POSITION_FIELDS:
            if frame.has_field(field):
                print(frame.get_value(field))

    LAP_FIELDS = ["latitude", "longitude", "lap_n", "distance", "speed", "stroke_rate"]


def session_summary(row):
    summary = {}
    for val in row.fields:
        try:
            # Remove tuples with all None values and unknown data names
            if val.value != None and all(val.value) and not "unknown" in val.name:
                summary.update({val.name: val.value})
        except TypeError:
            # all() doesn't work for datetime values
            if val.value != None and not "unknown" in val.name:
                if isinstance(val.value, datetime.datetime):
                    # Convert datetime to string
                    summary.update({val.name: str(val.value)})
                elif "lat" in val.name or "long" in val.name:
                    if val.value == 0:
                        # Indoor activity or no position data
                        pass
                    else:
                        convertPos = val.value / ((2**32) / 360)
                        summary.update({val.name: convertPos})
                else:
                    summary.update({val.name: val.value})
    print(json.dumps(summary, indent=4))

    return summary 


def parse_fit_file(fName):
    with fitdecode.FitReader(fName) as ff:
        for row in ff:
            if isinstance(row, fitdecode.records.FitDataMessage):
                if row.name == "session":
                    session_summary(row)

def main():
    parse_fit_file("swim-lap.fit")


if __name__ == "__main__":
    main()
