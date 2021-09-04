package main

import (
	//"fmt"
)

func GetBikeZones() map[string]int{
	zones := map[string]int{
		"OneMin": 0, 
		"OneMax": 145,
		"TwoMin": 146,
		"TwoMax": 200,
		"ThreeMin": 201,
		"ThreeMax": 261, // FTP
		"FourMin": 262,
		"FourMax": 2000,
	}
	
	return zones
}


type ZoneDistribution struct {
	ZoneOne   float64
	ZoneTwo   float64
	ZoneThree float64
	ZoneFour  float64
}
func TimeInZones(zones map[string]int, power []int) ZoneDistribution {
	var zoneDist ZoneDistribution

	return zoneDist
}
