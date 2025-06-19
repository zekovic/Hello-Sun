package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var urlAir string = "http://localhost:8070"
//var urlAir string = "http://linki.co:8070"


type BoundsResponse struct {
	Status string   `json:"status"`
	Data []struct {
		Lat float64   `json:"lat"`
		Lon float64   `json:"lon"`
		UID int32     `json:"uid"`
		AQI string    `json:"aqi"`
		Station struct {
			Name string   `json:"name"`
			Time string   `json:"time"`
		}   `json:"station"`
	}   `json:"data"`
}

type AirResult struct {
	Status struct {
		Token string   `json:"token"`
		Error string   `json:"error"`
	}   `json:"status"`
	City struct {
		Lat     float64    `json:"lat"`
		Lon     float64    `json:"long"`
		Name     string    `json:"name"`
		URL      string    `json:"url"`
		Location string    `json:"location"`
		Idx       int32    `json:"idx"`
	} `json:"city"`
	Values struct {
		PM10   float64   `json:"pm10"`
		PM25   float64   `json:"pm25"`
	}   `json:"values"`
}


func (air *AirResult) getAirQuality() error {
	
	if air.City.Lat == 0 || air.City.Lon == 0 {
		err := air.readLocationFromConfig()
		if err != nil {
			fmt.Printf("Err %v", err)
		}
		
	}
	
	if (air.City.Lat == 0 || air.City.Lon == 0) && air.City.Idx == 0 {
		return errors.New("Error while getting AQI location")
	}
	
	url := fmt.Sprintf("%v/aqi_info?lat=%v&lon=%v&uid=%v", urlAir, air.City.Lat, air.City.Lon, air.City.Idx)
	
	
	//fmt.Printf("url: [%v]", url)
	data, err := getFromApi(url)
	if err != nil {
		fmt.Printf("Err on request: %v", err)
	}
	//fmt.Printf("data: [%v]", string(data))
	//aqi_response := AirResult{}
	//err = json.Unmarshal(data, &aqi_response)
	err = json.Unmarshal(data, &air)
	if err != nil {
		fmt.Printf("Err on unmarshal: %v", err)
		//return AirResult{}, err
		return err
	}
	// return aqi_response, nil
	return nil
}
func (air *AirResult) saveLocationToConfig() {
	config.Location.AirQuality = air.City.Name
	config.Location.Lat_aqi = fmt.Sprintf("%v", air.City.Lat)
	config.Location.Lon_aqi = fmt.Sprintf("%v", air.City.Lon)
	config.Location.AirUID = fmt.Sprintf("%v", air.City.Idx)
}
func (air *AirResult) readLocationFromConfig() error {
	air.City.Name = config.Location.AirQuality
	
	latVal, err := strconv.ParseFloat(config.Location.Lat_aqi, 64)
	if err != nil {
		fmt.Printf("Err parsing AQI latitude: %v", err)
		return err
	}
	lonVal, err := strconv.ParseFloat(config.Location.Lon_aqi, 64)
	if err != nil {
		fmt.Printf("Err parsing AQI longitude: %v", err)
		return err
	}
	air.City.Lat = latVal
	air.City.Lon = lonVal
	
	uidVal, err := strconv.ParseInt(config.Location.AirUID, 10, 32)
	if err != nil {
		fmt.Printf("Err parsing AQI UID: %v", err)
		return err
	}
	air.City.Idx = int32(uidVal)
	return nil
}

func getAqiLocations(lat, lon float64) (BoundsResponse, error) {
	
	url := fmt.Sprintf("%v/list_locations?lat=%v&lon=%v", urlAir, lat, lon)
	//fmt.Printf("url: [%v]", url)
	data, err := getFromApi(url)
	if err != nil {
		fmt.Printf("Err on list request: %v", err)
	}
	//fmt.Printf("data: [%v]", string(data))
	boundsResponse := BoundsResponse{}
	err = json.Unmarshal(data, &boundsResponse)
	if err != nil {
		return BoundsResponse{}, err
	}
	return boundsResponse, nil
	
}

func getAqiColor(pm25, pm10 float64) (string, string) {
	color25 := "0 153 0"
	color10 := "0 153 0"
	
	if pm25 > 30.0 { color25 = "51 255 51" }
	if pm25 > 60.0 { color25 = "255 255 102" }
	if pm25 > 90.0 { color25 = "255 153 0" }
	if pm25 > 120.0 { color25 = "255 0 0" }
	if pm25 > 250.0 { color25 = "153 0 0" }
	
	if pm10 > 50.0 { color10 = "51 255 51" }
	if pm10 > 100.0 { color10 ="255 255 102" }
	if pm10 > 250.0 { color10 ="255 153 0" }
	if pm10 > 350.0 { color10 = "255 0 0" }
	if pm10 > 430.0 { color10 = "153 0 0" }
	
	return color25, color10
}






