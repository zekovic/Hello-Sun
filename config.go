package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/sys/windows/registry"
)

var configMap map [string]int

type MyConfig struct {
	Location struct{
		/*Search string		`json:"search"`
		Found string		`json:"found"`
		Url string			`json:"url"`*/
		Weather string		`json:"weather"`
		Lat string			`json:"lat"`
		Lon string			`json:"lon"`
		AirQuality string	`json:"air_quality"`
		Lat_aqi string		`json:"lat_aqi"`
		Lon_aqi string		`json:"lon_aqi"`
		AirUID string		`json:"air_uid"`
	} `json:"location"`
	//Location string			`json:"location"`
	Temperature string		`json:"temperature"`
	Wind string				`json:"wind"`
	Pressure string			`json:"pressure"`
	Rain string				`json:"rain"`
	Time string				`json:"time"`
	Refresh int				`json:"refresh"`
	Opacity string			`json:"opacity"`
	Systray string			`json:"systray"`
	//StartOnBoot string		`json:"start_on_boot"`
	AlwaysOnTop string		`json:"always_on_top"`
	Font string				`json:"font"`
	ShowParts []string		`json:"show_parts"`
	Position struct {
		X float64			`json:"x"`
		Y float64			`json:"y"`
	} `json:"position"`
}

func initConfigMap() {
	configMap = make(map[string]int)
	configMap["temp_C"] = 1
	configMap["temp_F"] = 2
	configMap["temp_K"] = 3
	configMap["wind_m/s"] = 1
	configMap["wind_km/h"] = 2
	configMap["wind_mph"] = 3
	configMap["wind_knots"] = 4
	configMap["wind_Bf"] = 5
	configMap["time_24h"] = 1
	configMap["time_AMPM"] = 2
	configMap["update_10"] = 1
	configMap["update_20"] = 2
	configMap["update_60"] = 3
	configMap["pressure_mBar"] = 1
	configMap["pressure_hPa"] = 2
	configMap["rain_mm"] = 1
	configMap["rain_in"] = 2
	configMap["parts_daily"] = 1
	configMap["parts_hourly"] = 2
	configMap["parts_details"] = 3
	configMap["parts_moon"] = 4
	configMap["parts_air"] = 5
	configMap["parts_sunrise"] = 6
	
}

func (c *MyConfig) init() {
	*c = MyConfig {
		Temperature: "C",
		Wind: "m/s",
		Pressure: "mBar",
		Rain: "mm",
		Time: "24h",
		Opacity: "230",
		Refresh: 20,
		Systray: "OFF",
		//StartOnBoot: "OFF",
		AlwaysOnTop: "YES",
		Font: "Tahoma",
		ShowParts: []string{
			"daily", /*"hourly",*/ "air", "details", "moon", "sunrise",
		},
	}
	c.Location.Weather = ""
	c.Location.AirQuality = ""
	c.Position.X = 30
	c.Position.Y = 30
}


func (c *MyConfig) read() {
	//time.Sleep(1 * time.Second)
	data, err := os.ReadFile("hello_sun_config.json")
	if err != nil {
		//saveConfig()
		c.save()
		return
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		fmt.Println("Error while reading config. ", err)
		//saveConfig()
		c.save()
	}
	
	if c.Refresh < 10 { c.Refresh = 10 }
	if c.getOpacityPercentage() < 30 { c.Opacity = "76" }
	if c.getOpacityPercentage() > 98 { c.Opacity = "255" }
	
	// parts - check if only existing, and not repeating:
	possibleParts := []string{}
	possiblePartsMap := make(map[string]int)
	for _, part := range c.ShowParts {
		_, ok := configMap["parts_" + part]
		if ok {
			_, notUnique := possiblePartsMap[part]
			if !notUnique {
				possibleParts = append(possibleParts, part)
				
			}
			possiblePartsMap[part] = 1
		}
	}
	c.ShowParts = possibleParts
}

func (c *MyConfig) save() {
	//time.Sleep(1 * time.Second)
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		fmt.Println("Error while making config json. ", err)
	}
	err = os.WriteFile("hello_sun_config.json", data, 0666)
	if err != nil {
		fmt.Println("Error while writing to config. ", err)
	}
}

func (c *MyConfig) setOpacity(percents string) {
	opacityVal, err := strconv.ParseFloat(percents, 64)
	if err != nil {
		opacityVal = 100.0
	}
	c.Opacity = strconv.Itoa(min(int(opacityVal * 2.55), 255))
}

func (c *MyConfig) getOpacityPercentage() int {
	opacityVal, err := strconv.ParseFloat(c.Opacity, 64)
	if err != nil {
		return 100
	}
	return min(int(opacityVal / 2.55), 100)
}

func (c *MyConfig) getRegistryBootValue() bool {
	
	regKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		fmt.Printf("Err (get - Reg path missing) %v", err)
		return false
	}
	defer regKey.Close()
	found_value, _, err := regKey.GetStringValue("HelloSun")
	if err != nil {
		fmt.Printf("Err (Start on boot is off) %v", err)
		return false
	}
	fmt.Printf("Boot val: %v", found_value)
	return true
}

func (c *MyConfig) setRegistryBootValue(bootStart bool) bool {
	regKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		fmt.Printf("Err (set - Reg path missing) %v", err)
		return false
	}
	defer regKey.Close()
	
	if bootStart {
		exe_path, err := os.Executable()
		if err != nil {
			fmt.Printf("Err while setting boot start (finding exe path)... %v", err)
			return false
		}
		err = regKey.SetStringValue("HelloSun", fmt.Sprintf("\"%v\"", exe_path))
		if err != nil {
			fmt.Printf("Err while setting boot start to true... %v", err)
		}
		return true
	} else {
		err := regKey.DeleteValue("HelloSun")
		if err != nil {
			fmt.Printf("Err while setting boot start to false... %v", err)
			return false
		}
		return true
	}
	
	return true
}
