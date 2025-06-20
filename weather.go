package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)


type WeatherResult struct {
	CurrentCondition []struct {
		FeelsLikeC string				`json:"FeelsLikeC"`
		FeelsLikeF string				`json:"FeelsLikeF"`
		Cloudcover string				`json:"cloudcover"`
		Humidity string					`json:"humidity"`
		LocalObsDateTime string			`json:"localObsDateTime"`
		CurrentTime string				`json:"currentTime"`
		Observation_time string			`json:"observation_time"`
		PrecipInches string				`json:"precipInches"`
		PrecipMM string					`json:"precipMM"`
		Pressure string					`json:"pressure"`
		PressureMbar string				`json:"pressureMbar"`
		PressureInches string			`json:"pressureInches"`
		PressureMmHg string				`json:"pressureMmHg"`
		PressurePsi string				`json:"pressurePsi"`
		Visibility string				`json:"visibility"`
		VisibilityMiles string			`json:"visibilityMiles"`
		Temp_C string					`json:"temp_C"`
		Temp_F string					`json:"temp_F"`
		Temp_K string					`json:"temp_K"`
		UvIndex string					`json:"uvIndex"`
		WeatherCode string				`json:"weatherCode"`
		WeatherCodeTxt string			`json:"weatherCodeTxt"`
		Winddir16Point string			`json:"winddir16Point"`
		WinddirDegree string			`json:"winddirDegree"`
		WinddirArrow string				`json:"winddirArrow"`
		WindspeedKmph string			`json:"windspeedKmph"`
		WindspeedMiles string			`json:"windspeedMiles"`
		WindspeedMps string				`json:"windspeedMps"`
		WindspeedKnots string			`json:"windspeedKnots"`
		WindspeedBf string				`json:"windspeedBf"`
		Sunrise string					`json:"sunrise"`
		Sunset string					`json:"sunset"`
		MoonIcon string					`json:"moonIcon"`
		MoonDay string					`json:"moonDay"`
	} `json:"current_condition"`
	NearestArea []struct {
		AreaName []struct {
			Value string				`json:"value"`
		}						`json:"areaName"`
		Latitude string					`json:"latitude"`
		Longitude string				`json:"longitude"`
	}							`json:"nearest_area"`
	Weather []WeatherDay		`json:"weather"`
}

type WeatherDay struct {
	Astronomy []struct {
		Moon_illumination string	`json:"moon_illumination"`
		Moon_phase string			`json:"moon_phase"`
		Moonrise string				`json:"moonrise"`
		Moonset string				`json:"moonset"`
		Sunrise string				`json:"sunrise"`
		Sunset string				`json:"sunset"`
	}						`json:"astronomy"`
	AvgtempC string					`json:"avgtempC"`
	AvgtempF string					`json:"avgtempF"`
	Date string						`json:"date"`
	MaxtempC string					`json:"maxtempC"`
	MaxtempF string					`json:"maxtempF"`
	MintempC string					`json:"mintempC"`
	MintempF string					`json:"mintempF"`
	SunHour string					`json:"sunHour"`
	TotalSnow_cm string				`json:"totalSnow_cm"`
	UvIndex string					`json:"uvIndex"`
}



func (wr *WeatherResult) getAndParse(location string) {
	if location == "" {
		location = config.Location.Weather
	}
	//defer wg.Done() // ??
	//fmt.Println("This.")
	timeStr := strconv.FormatInt(time.Now().Unix(), 10)
	randomStr := fmt.Sprintf("%x", md5.Sum([]byte(timeStr)  ))
	//fmt.Printf("Random: %v", random_str)
	//url := fmt.Sprintf("https://wttr.in/%v?format=j2&nonce=%v", config.Location.Weather, random_str)
	url := fmt.Sprintf("https://wttr.in/%v?format=j2&nonce=%v", url.QueryEscape(location), randomStr)
	// url := fmt.Sprintf("https://wttr.in/%v?format=j1&nonce=%v", url.QueryEscape(location), random_str)
	
	data, err := getFromApi(url)
	if err != nil {
		fmt.Printf("Err on request: %v", err)
	}
	
	*wr, err = getWeatherFromJson(data)
	if err != nil {
		fmt.Printf("Err on JSON: %v", err)
	}
	fmt.Printf("Result: %v \n", wr)
}



func (wr *WeatherResult) getCurrentLocation() (Place, error) {
	if len(wr.NearestArea) == 0 {
		return Place{}, errors.New("Area info is empty")
	}
	if len(wr.NearestArea[0].AreaName) == 0 {
		return Place{}, errors.New("Area name is empty")
	}
	if wr.NearestArea[0].AreaName[0].Value == "" {
		return Place{}, errors.New("Area name value is empty")
	}
	return Place{
		Name: wr.NearestArea[0].AreaName[0].Value,
		Lat: wr.NearestArea[0].Latitude,
		Lon: wr.NearestArea[0].Longitude,
	}, nil
}

func (wr *WeatherResult) saveLocationToConfig() {
	config.Location.Weather = wr.NearestArea[0].AreaName[0].Value
	config.Location.Lat = wr.NearestArea[0].Latitude
	config.Location.Lon = wr.NearestArea[0].Longitude
	//config.save()
}

func (wr *WeatherResult) isNight() bool {
	if len(wr.CurrentCondition) == 0 {
		return false
	}
	if wr.CurrentCondition[0].CurrentTime < wr.CurrentCondition[0].Sunrise ||
		wr.CurrentCondition[0].Sunset < wr.CurrentCondition[0].CurrentTime {
		return true
	}
	return false
}



func getWeatherFromJson(data []byte) (WeatherResult, error) {
	weatherResult := WeatherResult{}
	err := json.Unmarshal(data, &weatherResult)
	if err != nil {
		return WeatherResult{}, err
	}
	return weatherResult, nil
}





