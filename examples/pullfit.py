import fitparse
import matplotlib.pyplot as plt


fitfile = fitparse.FitFile("swim-lap.fit")

def parse_swim_indoors(fitfile):
    jsonLaps = []
    for record in fitfile.get_messages("session"):
        print(record.get_values())

   # print(jsonLaps[0])

   # lapTimes = []
   # for lap in jsonLaps:
   #     lapTimes.append(lap["total_elapsed_time"])

   # plt.plot(lapTimes)
   # plt.show()

def main():
    fitfile = fitparse.FitFile("swim-lap.fit")
    for sport in fitfile.get_messages("sport"):
        sportType = sport.get_values()
        if sportType["sport"] == "swimming":
            parse_swim_indoors(fitfile)
        break


if __name__ == "__main__":
    main()
