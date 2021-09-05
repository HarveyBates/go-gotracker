package main

import (
//	"fmt"
)



func main() {

	refreshToken := GetRefreshToken()

	distanceArr, wattsArr := GetWatts("5901981172", refreshToken)

	MakeChart(distanceArr, wattsArr, "Workout #1", "2 x 20 min @ 240 Watts")
}
