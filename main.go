package main

import (
	"strings"
	"sync"

	//"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var weatherResult WeatherResult
var aqiResult AirResult
var config MyConfig
var wg sync.WaitGroup

func main() {
	
	initImagesInfo()
	initConfigMap()
	config.init()
	config.read()
	showGui()
	
}

func getFromApi(urlString string) ([]byte, error) {
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	
	client := &http.Client{Timeout: time.Second * 10}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	//fmt.Printf("Response: [%v], REQ: [%v]", string(body), req)
	
	//return string(body), nil
	return body, nil
}


func getTempByUnit() string {
	if len(weatherResult.CurrentCondition) == 0 {
		return "__"
	}
	if config.Temperature == "C" {
		return weatherResult.CurrentCondition[0].Temp_C + "°C"
	}
	if config.Temperature == "F" {
		return weatherResult.CurrentCondition[0].Temp_F + "°F"
	}
	if config.Temperature == "K" {
		return weatherResult.CurrentCondition[0].Temp_K + "K"
	}
	return "__"
}

func TempDayInfo(item WeatherDay) (string, string, string, string) {
	if config.Temperature == "C" {
		return item.MintempC, item.MaxtempC, item.AvgtempC, "°C"
	}
	if config.Temperature == "F" {
		return item.MintempF, item.MaxtempF, item.AvgtempF, "°F"
	}
	if config.Temperature == "K" {
		t_c_min_int, _ := strconv.Atoi(item.MintempC)
		t_c_max_int, _ := strconv.Atoi(item.MaxtempC)
		t_c_avg_int, _ := strconv.Atoi(item.AvgtempC)
		return strconv.Itoa(t_c_min_int + 273),
			strconv.Itoa(t_c_max_int + 273),
			strconv.Itoa(t_c_avg_int + 273),
			"K"
	}
	return "", "", "", ""
}

func getUV() string {
	if len(weatherResult.CurrentCondition) == 0 {
		return "_"
	}
	return weatherResult.CurrentCondition[0].UvIndex
}


func getInTextFormat() {
	//defer wg.Done() // ??
	
	// wttr.in/lib/parse_query.py   use_metric=true / use_imperial=true
	// ?u - imperial   ?m - metric km/h   ?M - metric m/s
	url := fmt.Sprintf("https://wttr.in/%v?m&format=%%x__%%t__%%f__%%w__%%m__%%M__%%S__%%s__%%P__%%u__%%p__%%h", 
		url.QueryEscape(config.Location.Weather))
	data, err := getFromApi(url)
	if err != nil {
		fmt.Printf("Err on request: %v", err)
	}
	data_str := string(data)
	//fmt.Printf("Result: %v \n", data_str)
	data_arr := strings.Split(data_str, "__")
	fmt.Printf("Result: %v \n", data_arr)
	
	if len(weatherResult.CurrentCondition) == 0 || len(data_arr) < 11 {
		return
	}
	
	weatherResult.CurrentCondition[0].WeatherCodeTxt = data_arr[0]
	
	t := getTempByCelsius(data_arr[1])
	weatherResult.CurrentCondition[0].Temp_C = t[0]
	weatherResult.CurrentCondition[0].Temp_F = t[1]
	weatherResult.CurrentCondition[0].Temp_K = t[2]
	
	wind := getWindByKmh(data_arr[3])
	fmt.Printf("Wind: %v \n", wind)
	weatherResult.CurrentCondition[0].WindspeedKmph = wind[0]
	weatherResult.CurrentCondition[0].WindspeedMiles = wind[1]
	weatherResult.CurrentCondition[0].WindspeedMps = wind[2]
	weatherResult.CurrentCondition[0].WindspeedKnots = wind[3]
	weatherResult.CurrentCondition[0].WindspeedBf = wind[4]
	weatherResult.CurrentCondition[0].WinddirArrow = wind[5]
	
	weatherResult.CurrentCondition[0].MoonIcon = data_arr[4]
	weatherResult.CurrentCondition[0].MoonDay = data_arr[5]
	
	weatherResult.CurrentCondition[0].Sunrise = data_arr[6][0:5]
	weatherResult.CurrentCondition[0].Sunset = data_arr[7][0:5]
	
	pr_hPa := strings.ReplaceAll(data_arr[8], "hPa", "");
	pr_hPa_i, _ := strconv.Atoi(pr_hPa)
	pr_hPa_f := float64(pr_hPa_i)
	weatherResult.CurrentCondition[0].Pressure = pr_hPa
	weatherResult.CurrentCondition[0].PressureMbar = pr_hPa
	weatherResult.CurrentCondition[0].PressureInches = strconv.Itoa(int(pr_hPa_f * 0.0295299833))
	weatherResult.CurrentCondition[0].PressureMmHg = strconv.Itoa(int(pr_hPa_f * 0.750062))
	weatherResult.CurrentCondition[0].PressurePsi = strconv.Itoa(int(pr_hPa_f * 0.0145038))
	
	pre_mm := strings.ReplaceAll(data_arr[9], "mm", "");
	pre_mm_i, _ := strconv.Atoi(pre_mm)
	weatherResult.CurrentCondition[0].PrecipMM = pre_mm
	weatherResult.CurrentCondition[0].PrecipInches = strconv.Itoa(int(float64(pre_mm_i) * 0.03937))
	
	weatherResult.CurrentCondition[0].Humidity = data_arr[10]
}

func getTempByCelsius(temp_c string) [3]string {
	result := [3]string {"", "", ""}
	
	t_c := strings.ReplaceAll(temp_c, "°C", "")
	t_c_int, _ := strconv.Atoi(t_c)
	t_f_int := int((float64(t_c_int) * 1.8) + 32)
	
	result[0] = strconv.Itoa(t_c_int)
	result[1] = strconv.Itoa(t_f_int)
	result[2] = strconv.Itoa(t_c_int + 273)
	
	return result
}
func getWindByKmh(wind_kmh string) [6]string {
	result := [6]string {"", "", "", "", "", ""}
	windStr := strings.ReplaceAll(wind_kmh, "km/h", "")
	windArrow := string([]rune(windStr)[0:1])
	windVal := string([]rune(windStr)[1:])
	//fmt.Printf("arrow: [%v], str: [%v]\n", wind_arrow, wind_val)
	windKmhInt, err := strconv.Atoi(windVal)
	if err != nil {
		fmt.Println("Err wind: ", err)
	}
	windMphInt := int(float64(windKmhInt) * 0.621371)
	windMsInt := int(float64(windKmhInt) * 0.27778)
	windKnotsInt := int(float64(windKmhInt) * 0.53996)
	windBfInt := 0
	
	if windKmhInt >= 1   { windBfInt = 1  }
	if windKmhInt >= 6   { windBfInt = 2  }
	if windKmhInt >= 12  { windBfInt = 3  }
	if windKmhInt >= 20  { windBfInt = 4  }
	if windKmhInt >= 29  { windBfInt = 5  }
	if windKmhInt >= 39  { windBfInt = 6  }
	if windKmhInt >= 50  { windBfInt = 7  }
	if windKmhInt >= 62  { windBfInt = 8  }
	if windKmhInt >= 75  { windBfInt = 9  }
	if windKmhInt >= 89  { windBfInt = 10 }
	if windKmhInt >= 103 { windBfInt = 11 }
	if windKmhInt >= 118 { windBfInt = 12 }
	if windKmhInt >= 133 { windBfInt = 13 }
	if windKmhInt >= 149 { windBfInt = 14 }
	if windKmhInt >= 166 { windBfInt = 15 }
	if windKmhInt >= 184 { windBfInt = 16 }
	if windKmhInt >= 200 { windBfInt = 17 }
	
	result[0] = strconv.Itoa(windKmhInt)
	result[1] = strconv.Itoa(windMphInt)
	result[2] = strconv.Itoa(windMsInt)
	result[3] = strconv.Itoa(windKnotsInt)
	result[4] = strconv.Itoa(windBfInt)
	result[5] = windArrow
	return result
}

func getWindByUnit() (string, string) {
	result := ""
	arrow := ""
	if len(weatherResult.CurrentCondition) == 0 {
		return result, arrow
	}
	arrow = weatherResult.CurrentCondition[0].WinddirArrow
	if config.Wind == "m/s" { result = weatherResult.CurrentCondition[0].WindspeedMps + " m/s" }
	if config.Wind == "km/h" { result = weatherResult.CurrentCondition[0].WindspeedKmph + " km/h" }
	if config.Wind == "mph" { result = weatherResult.CurrentCondition[0].WindspeedMiles + " mph" }
	if config.Wind == "knots" { result = weatherResult.CurrentCondition[0].WindspeedKnots + " knots" }
	if config.Wind == "Bf" { result = "Bf " + weatherResult.CurrentCondition[0].WindspeedBf }
	return result, arrow
}

func getPressureByUnit() string {
	result := ""
	if config.Pressure == "mBar" { result = weatherResult.CurrentCondition[0].PressureMbar + " mBar" }
	if config.Pressure == "hPa" { result = weatherResult.CurrentCondition[0].Pressure + " hPa" }
	if config.Pressure == "inHg" { result = weatherResult.CurrentCondition[0].PressureInches + " inHg" }
	if config.Pressure == "mmHg" { result = weatherResult.CurrentCondition[0].PressureMmHg + " mmHg" }
	if config.Pressure == "psi" { result = weatherResult.CurrentCondition[0].PressurePsi + " psi" }
	return result
}

func getPrecipitationByUnit() string {
	result := ""
	if config.Rain == "mm" { result = weatherResult.CurrentCondition[0].PrecipMM + " mm" }
	if config.Rain == "in" { result = weatherResult.CurrentCondition[0].PrecipInches + " in" }
	return result
}