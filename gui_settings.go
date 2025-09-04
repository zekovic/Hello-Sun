package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gen2brain/iup-go/iup"
)

var wndSettings iup.Ihandle
var wndInput WndInput

var btnLocation iup.Ihandle
var btnLocationAir iup.Ihandle
var labelLocation iup.Ihandle
//var changed_location string
var labelLocationAir iup.Ihandle
var listLocations iup.Ihandle
var btnRefreshLocations iup.Ihandle
var imgLogoLocations iup.Ihandle
var comboCFK iup.Ihandle
var comboWind iup.Ihandle
var comboPressure iup.Ihandle
var comboRain iup.Ihandle
var comboTime iup.Ihandle
var comboRefresh iup.Ihandle
var labelOpacity iup.Ihandle
var sliderOpacity iup.Ihandle
var checkSystray iup.Ihandle
var checkBoot iup.Ihandle
var checkOnTop iup.Ihandle
var labelFont iup.Ihandle
var labelFontVal iup.Ihandle
var btnFont iup.Ihandle
var dlgFont iup.Ihandle
var checkPartsWrap iup.Ihandle

var weatherTmp WeatherResult
var aqiResultTmp AirResult

func createSettingsDialog() {
	
	btnAttrs := `USERSIZE="75x24", PADDING="12x8"`
	inputAttrs := `USERSIZE="120x24"`
	
	btnSave := iup.Button("&OK").SetCallback("ACTION", iup.ActionFunc(onClickSettingsSave))
	btnCancel := iup.Button("&Cancel").SetCallback("ACTION", iup.ActionFunc(onClickSettingsClose))
	btnApply := iup.Button("&Apply").SetCallback("ACTION", iup.ActionFunc(onClickSettingsApply))
	btnLocation = iup.Button("Change").SetCallback("ACTION", iup.ActionFunc(onClickSettingsLocation))
	btnLocationAir = iup.Button("Change").SetCallback("ACTION", iup.ActionFunc(onClickSettingsLocationAir))
	btnRefreshLocations = iup.Button("Refresh").SetCallback("ACTION", iup.ActionFunc(onClickRefreshLocations))
	labelLocation = iup.Label("")//.SetAttributes(`USERSIZE="150x24", PADDING="12x8"`)
	labelLocationAir = iup.Label("")//.SetAttributes(`USERSIZE="150x24", PADDING="12x8"`)
	listLocations = iup.FlatList()
	comboCFK = iup.List().SetAttributes(`DROPDOWN=YES, 1="Celsius", 2="Fahrenheit", 3="Kelvin", X1="C", X2="F", X3="K"`).SetAttributes(inputAttrs)
	comboWind = iup.List().SetAttributes(`DROPDOWN=YES, 1="m/s", 2="km/h", 3="mph", 4="knots", 5="Bf", X1="m/s", X2="km/h", X3="mph", X4="knots", X5="Bf"`).SetAttributes(inputAttrs)
	comboPressure = iup.List().SetAttributes(`DROPDOWN=YES, 1="mBar (mBar=hPa)", 2="hPa (mBar=hPa)", 3="inHg", 4="mmHg", 5="psi", X1="mBar", X2="hPa", X3="inHg", X4="mmHg", X5="psi"`).SetAttributes(inputAttrs)
	comboRain = iup.List().SetAttributes(`DROPDOWN=YES, 1="mm", 2="in", X1="mm", X2="in"`).SetAttributes(inputAttrs)
	comboTime = iup.List().SetAttributes(`DROPDOWN=YES, 1="24 hour", 2="AM/PM", X1="24h", X2="AMPM"`).SetAttributes(inputAttrs)
	comboRefresh = iup.List().SetAttributes(`DROPDOWN=YES, 1="10 min", 2="20 min", 3="1 hour", X1="10", X2="20", X3="60"`).SetAttributes(inputAttrs)
	labelOpacity = iup.Label("0%").SetAttributes(`USERSIZE="35x24"`)
	sliderOpacity = iup.FlatVal("HORIZONTAL").SetAttributes("USERSIZE=120x24,MIN=30,MAX=100")
	checkSystray = iup.Toggle("Leave on systray on close")
	checkBoot = iup.Toggle("Start with sistem boot")
	checkOnTop = iup.Toggle("Always on top")
	
	labelFont = iup.Label("Font:").SetAttributes(`USERSIZE="60x24"`)
	labelFontVal = iup.Label("(choose font...)").SetAttributes(`USERSIZE="200x24"`)
	btnFont = iup.Button("Change...").SetAttributes(`USERSIZE="60x24"`).SetCallback("ACTION", iup.ActionFunc(onClickChangeFont))
	dlgFont = iup.FontDlg()
	
	hAttrs := `ExXPANDCHILDREN=YES, UxSERSIZE="150x20"`
	checkPartsWrap = iup.Vbox(
		iup.Hbox(newUpBtn(), newDownBtn(), iup.Label("  "), 
			iup.Toggle("Daily Forecast").SetHandle("daily"), ).SetAttributes(hAttrs),
		iup.Hbox(newUpBtn(), newDownBtn(), iup.Label("  "), 
			iup.Toggle("Hourly Forecast").SetHandle("hourly"), ).SetAttributes(hAttrs),
		iup.Hbox(newUpBtn(), newDownBtn(), iup.Label("  "), 
			iup.Toggle("Moon Phase").SetHandle("moon"), ).SetAttributes(hAttrs),
		iup.Hbox(newUpBtn(), newDownBtn(), iup.Label("  "), 
			iup.Toggle("Air Quality").SetHandle("air"), ).SetAttributes(hAttrs),
		iup.Hbox(newUpBtn(), newDownBtn(), iup.Label("  "), 
			iup.Toggle("Details").SetHandle("details"), ).SetAttributes(hAttrs),
		iup.Hbox(newUpBtn(), newDownBtn(), iup.Label("  "), 
			iup.Toggle("Sunrise and Sunset").SetHandle("sunrise"), ).SetAttributes(hAttrs),
		//iup.Fill(),
	)
	
	btnSave.SetAttributes(btnAttrs)
	btnCancel.SetAttributes(btnAttrs)
	btnApply.SetAttributes(btnAttrs)
	btnLocation.SetAttributes(`USERSIZE="60x24"`)
	btnRefreshLocations.SetAttributes(`USERSIZE="60x24"`)
	btnLocationAir.SetAttributes(`USERSIZE="60x24"`)
	
	sliderOpacity.SetCallback("VALUECHANGED_CB", iup.ValueChangedFunc(onOpacityChange))
	
	listLocations.SetAttributes(`USERSIZE="350x210"`)
	listLocations.SetCallback("VALUECHANGED_CB", iup.ValueChangedFunc(onListLocationsChanged))
	btnRefreshLocations.SetAttributes(``)
	
	wndSettings = iup.Dialog(
		iup.Vbox(
			iup.Tabs(
				iup.Vbox(
						iup.Hbox(
							iup.Label("Weather:").SetAttributes(`USERSIZE="60x24"`), 
							labelLocation.SetAttributes(`USERSIZE="200x24"`),
							//iup.Fill(),
							btnLocation,
						).SetAttributes(`CGAP=10`),
						iup.Hbox(
							iup.Label("Air quality locations:").SetAttributes(`USERSIZE="282x24"`), 
							iup.Hbox(btnRefreshLocations, imgLogoLocations),
						).SetAttributes(`CGAP=2`),
					iup.Hbox(
						listLocations,
					),
					iup.Vbox(
						iup.Hbox(
							iup.Label("Selected:").SetAttributes(`USERSIZE="60x24"`), 
							labelLocationAir.SetAttributes(`USERSIZE="180x24"`),
							//iup.Fill(),
							//btn_location_air,
						).SetAttributes(`CGAP=2`),
						
					),
					iup.Fill(),
					
				).SetAttributes(`MARGIN="5x5", TABTITLE="Location"`),
				iup.Vbox(
					iup.Hbox(iup.Label("Temperature unit:"), iup.Fill(), comboCFK,),
					iup.Hbox(iup.Label("Wind speed unit:"), iup.Fill(), comboWind, ),
					iup.Hbox(iup.Label("Pressure unit:"), iup.Fill(), comboPressure, ),
					iup.Hbox(iup.Label("Rain unit:"), iup.Fill(), comboRain, ),
					iup.Hbox(iup.Label("Time format:"), iup.Fill(), comboTime, ),
					iup.Fill(),
				).SetAttributes(`MARGIN="5x5", TABTITLE="Units"`),
				iup.Vbox(
					iup.Hbox(iup.Label("Refresh rate:"), iup.Fill(), comboRefresh, ),
					iup.Hbox(iup.Label("Opacity:"), iup.Fill(), labelOpacity, sliderOpacity,),
					iup.Hbox(checkSystray, ),
					iup.Hbox(checkBoot, ),
					iup.Hbox(checkOnTop, ),
					iup.Hbox(labelFont, labelFontVal, iup.Fill(), btnFont, ),
					iup.Fill(),
				).SetAttributes(`MARGIN="5x5", TABTITLE="Display"`),
				checkPartsWrap.SetAttributes(`MARGIN="5x5", TABTITLE="Parts"`),
				iup.Vbox(
					iup.Hbox(iup.Label("About this..."), ),
					iup.Fill(),
				).SetAttributes(`MARGIN="5x5", TABTITLE="About"`),
				//iup.Fill(),
			),
			//iup.Fill(),
			iup.Hbox(iup.Fill(), btnSave, btnCancel, btnApply,).SetAttributes(`MARGIN="5x5", CGAP=3`),
		),
	//).SetAttributes(`USERSIZE="500x300",TITLE="Settings",RESIZE=NO,SIMULATEMODAL=YES`)
	//).SetAttributes(`USERSIZE="500x300",TITLE="Settings",RESIZE=NO,HIDETASKBAR=YES,TASKBARBUTTON=HIDE,TOPMOST=YES,MENUBOX=NO`)
	).SetAttributes(`USERSIZE="410x500",TITLE="Settings",RESIZE=NO,MENUBOX=NO`)
	wndSettings.SetHandle("wnd_settings")
	
	wndInput.Create()
}


func newUpBtn() iup.Ihandle {
	btn := iup.Button("▲").SetAttributes(`USERSIZE="16x16", BTN_TYPE="up"`)
	btn.SetCallback("ACTION", iup.ActionFunc(btnUpClick))
	return btn
}
func newDownBtn() iup.Ihandle {
	btn := iup.Button("▼").SetAttributes(`USERSIZE="16x16", BTN_TYPE="down"`)
	btn.SetCallback("ACTION", iup.ActionFunc(btnDownClick))
	return btn
}


func onOpacityChange(ih iup.Ihandle) int {
	opacity, _ := strconv.ParseFloat(ih.GetAttribute("VALUE"), 32)
	labelOpacity.SetAttribute("TITLE", fmt.Sprintf("%.0f%%", opacity))
	return iup.DEFAULT
}

func onListLocationsChanged(ih iup.Ihandle) int {
	selected := listLocations.GetAttribute("VALUE")
	locationName := listLocations.GetAttribute(fmt.Sprintf("name_%v", selected))
	labelLocationAir.SetAttribute("TITLE", locationName)
	return iup.DEFAULT
}

func btnUpClick(ih iup.Ihandle) int {
	thisHb := iup.GetParent(ih)
	thisIndex := iup.GetChildPos(checkPartsWrap, thisHb)
	if thisIndex == 0 {
		return iup.DEFAULT
	}
	prevHb := iup.GetChild(checkPartsWrap, thisIndex - 1)
	iup.Reparent(thisHb, checkPartsWrap, prevHb)
	iup.Refresh(checkPartsWrap)
	return iup.DEFAULT
}
func btnDownClick(ih iup.Ihandle) int {
	thisHb := iup.GetParent(ih)
	thisIndex := iup.GetChildPos(checkPartsWrap, thisHb)
	if thisIndex == iup.GetChildCount(checkPartsWrap) - 1 {
		return iup.DEFAULT
	}
	nextHb := iup.GetChild(checkPartsWrap, thisIndex + 1)
	iup.Reparent(nextHb, checkPartsWrap, thisHb)
	iup.Refresh(checkPartsWrap)
	return iup.DEFAULT
}

func onClickChangeFont(ih iup.Ihandle) int {
	dlgFont.SetAttribute("VALUE", labelFontVal.GetAttribute("TITLE") + ", 12")
	iup.Popup(dlgFont, iup.CURRENT, iup.CURRENT)
	if dlgFont.GetAttribute("STATUS") == "1" {
		labelFontVal.SetAttribute("TITLE", strings.Split(dlgFont.GetAttribute("VALUE"), ",")[0])
	}
	return iup.DEFAULT
}

func comboGetValue(ih iup.Ihandle) string {
	return ih.GetAttribute("X"+ih.GetAttribute("VALUE"))
}

func fillSettingsWindow() {
	labelLocation.SetAttribute("TITLE", config.Location.Weather)
	//changed_location = config.Location.Weather
	
	labelLocationAir.SetAttribute("TITLE", config.Location.AirQuality)
	listLocations.SetAttribute("1", nil)
	listLocations.SetAttribute("1", fmt.Sprintf("%v [%v]", config.Location.AirQuality, config.Location.AirUID))
	listLocations.SetAttribute("x1", fmt.Sprintf("%v", config.Location.AirUID))
	listLocations.SetAttribute("lat_1", fmt.Sprintf("%v", config.Location.Lat_aqi))
	listLocations.SetAttribute("lon_1", fmt.Sprintf("%v", config.Location.Lon_aqi))
	listLocations.SetAttribute("name_1", fmt.Sprintf("%v", config.Location.AirQuality))
	listLocations.SetAttribute("VALUE", "1")
	
	comboCFK.SetAttribute("VALUE", configMap["temp_"+config.Temperature])
	comboWind.SetAttribute("VALUE", configMap["wind_"+config.Wind])
	comboTime.SetAttribute("VALUE", configMap["time_"+config.Time])
	comboPressure.SetAttribute("VALUE", configMap["pressure_"+config.Pressure])
	comboRain.SetAttribute("VALUE", configMap["rain_"+config.Rain])
	
	sliderOpacity.SetAttribute("VALUE", config.getOpacityPercentage())
	labelOpacity.SetAttribute("TITLE", fmt.Sprintf("%v%%", config.getOpacityPercentage()))
	comboRefresh.SetAttribute("VALUE", configMap[fmt.Sprintf("update_%v",config.Refresh)])
	if config.AlwaysOnTop == "YES" {
		checkOnTop.SetAttribute("VALUE", "ON")
	} else {
		checkOnTop.SetAttribute("VALUE", "OFF")
	}
	// checkBoot.SetAttribute("VALUE", config.StartOnBoot)
	if (config.getRegistryBootValue()) {
		checkBoot.SetAttribute("VALUE", "ON")
	} else {
		checkBoot.SetAttribute("VALUE", "OFF")
	}
	checkSystray.SetAttribute("VALUE", config.Systray)
	labelFontVal.SetAttribute("TITLE", config.Font)
	
	for _, uiPart := range config.ShowParts {
		foundPart := iup.GetHandle(uiPart)
		if foundPart != 0 {
			foundPart.SetAttribute("VALUE", "ON")
			iup.Reparent(iup.GetParent(iup.GetHandle(uiPart)), checkPartsWrap, 0)
			iup.Refresh(checkPartsWrap)
		}
	}
	
}

func saveSettings() {
	config.Location.Weather = labelLocation.GetAttribute("TITLE")
	
	config.Location.AirQuality = labelLocationAir.GetAttribute("TITLE")
	selected := listLocations.GetAttribute("VALUE")
	fmt.Printf("SELECTED: %v \n", selected)
	config.Location.AirUID = listLocations.GetAttribute("x"+selected)
	config.Location.Lat_aqi = listLocations.GetAttribute("lat_"+selected)
	config.Location.Lon_aqi = listLocations.GetAttribute("lon_"+selected)
	
	config.Temperature = comboGetValue(comboCFK)
	config.Wind = comboGetValue(comboWind)
	config.Pressure = comboGetValue(comboPressure)
	config.Rain = comboGetValue(comboRain)
	config.Time = comboGetValue(comboTime)
	
	config.Refresh, _ = strconv.Atoi(comboGetValue(comboRefresh))
	theTimer.SetAttribute("TIME", config.Refresh * 60 * 1000)
	config.setOpacity(sliderOpacity.GetAttribute("VALUE"))
	config.Systray = checkSystray.GetAttribute("VALUE")
	// config.StartOnBoot = checkBoot.GetAttribute("VALUE")
	config.setRegistryBootValue(checkBoot.GetAttribute("VALUE") == "ON")
	if checkOnTop.GetAttribute("VALUE") == "ON" {
		config.AlwaysOnTop = "YES"
	} else {
		config.AlwaysOnTop = "NO"
	}
	dlg.SetAttribute("TOPMOST", config.AlwaysOnTop)
	wndSettings.SetAttribute("TOPMOST", config.AlwaysOnTop)
	config.Font = labelFontVal.GetAttribute("TITLE")
	
	
	config.ShowParts = nil
	h_count := iup.GetChildCount(iup.GetChild(checkPartsWrap, 0))
	for i := 0; i < iup.GetChildCount(checkPartsWrap); i++ {
		part_h := iup.GetChild(checkPartsWrap, i)
		part_check := iup.GetChild(part_h, h_count - 1)
		if part_check.GetAttribute("VALUE") == "ON" {
			config.ShowParts = append(config.ShowParts, iup.GetName(part_check))
		}
		
	}
	
	fmt.Printf("\nCONFIG: [%v]\n", config)
	config.save()
	iup.Show(dlg)
}

func onClickSettingsSave(ih iup.Ihandle) int {
	fmt.Print("Saving...")
	toUpdate := isLocationChanged()
	saveSettings()
	settingsWindow(false)
	if toUpdate != "" {
		aqiResult.readLocationFromConfig()
		go fetchData()
	}
	return iup.DEFAULT
}
func onClickSettingsApply(ih iup.Ihandle) int {
	fmt.Print("Apply...")
	toUpdate := isLocationChanged()
	saveSettings()
	if toUpdate != "" {
		aqiResult.readLocationFromConfig()
		go fetchData()
	}
	return iup.DEFAULT
}


func isLocationChanged() string {
	changedWeather := ""
	changedAir := ""
	
	if labelLocation.GetAttribute("TITLE") != "" && config.Location.Weather != labelLocation.GetAttribute("TITLE") {
		changedWeather = fmt.Sprintf("Weather changed from %v to %v.\n", 
			config.Location.Weather, labelLocation.GetAttribute("TITLE"))
	}
	if labelLocationAir.GetAttribute("TITLE") != "" && config.Location.AirQuality != labelLocationAir.GetAttribute("TITLE") {
		changedAir = fmt.Sprintf("Air quality changed from %v to %v.\n", 
			config.Location.AirQuality, labelLocationAir.GetAttribute("TITLE"))
	}
	
	return changedWeather + changedAir
}

func onClickSettingsClose(ih iup.Ihandle) int {
	fmt.Print("Closing...")
	toUpdate := isLocationChanged()
	if toUpdate != "" {
		question := fmt.Sprintf("%v \nDo you want to save the changes?", toUpdate)
		locationMsg := iup.MessageDlg()
		locationMsg.SetAttributes(`PARENTDIALOG="wnd_settings", DIALOGTYPE=QUESTION, BUTTONS=YESNO, TITLE="Update location?"`)
		locationMsg.SetAttribute("VALUE", question)
		iup.Popup(locationMsg, iup.CURRENT, iup.CURRENT)
		to_save := locationMsg.GetAttribute("BUTTONRESPONSE")
		fmt.Printf("Answer: [%v]", to_save)
		
		if to_save == "1" {
			//config.Location.Weather = changed_location
			//config.Location.Weather = txt_location.GetAttribute("TITLE")
			if len(weatherTmp.NearestArea) > 0 {
				config.Location.Weather = weatherTmp.NearestArea[0].AreaName[0].Value
				config.Location.Lat = weatherTmp.NearestArea[0].Latitude
				config.Location.Lon = weatherTmp.NearestArea[0].Longitude
			}
			
			//aqi_result.City.Name = txt_location.GetAttribute("TITLE")
			
			aqiResult.City.Name = aqiResultTmp.City.Name
			aqiResult.City.Lat = aqiResultTmp.City.Lat
			aqiResult.City.Lon = aqiResultTmp.City.Lon
			aqiResult.City.Idx = aqiResultTmp.City.Idx
			
			//weather_result.saveLocationToConfig()
			selected := listLocations.GetAttribute("VALUE")
			fmt.Printf("SELECTED: %v \n", selected)
			config.Location.AirQuality = labelLocationAir.GetAttribute("TITLE")
			///// config.Location.AirQuality = listLocations.GetAttribute("name_"+selected)
			config.Location.AirUID = listLocations.GetAttribute("x"+selected)
			config.Location.Lat_aqi = listLocations.GetAttribute("lat_"+selected)
			config.Location.Lon_aqi = listLocations.GetAttribute("lon_"+selected)
			config.save()
			go fetchData()
		}
	}
	settingsWindow(false)
	return iup.DEFAULT
}

func onClickSettingsLocation(ih iup.Ihandle) int {
	
	oldLocation := labelLocation.GetAttribute("TITLE")
	newLocation := wndInput.GetInput(wndSettings, "Enter location", "Location")
	newLocation = strings.TrimSpace(newLocation)
	if newLocation != "" && newLocation != oldLocation {
		
		fmt.Print("Checking location...")
		
		go func() {
			labelLocation.SetAttribute("ACTIVE", "NO")
			btnLocation.SetAttribute("ACTIVE", "NO")
			
			defer labelLocation.SetAttribute("ACTIVE", "YES")
			defer btnLocation.SetAttribute("ACTIVE", "YES")
			
			weatherTmp.getAndParse(newLocation)
			
			found, err := weatherTmp.getCurrentLocation()
			if err != nil {
				iup.MessageError(wndSettings, "Error while getting entered location")
				fmt.Printf("Location err: %v \n", err)
				
			} else {
				
				labelLocation.SetAttribute("TITLE", found.Name)
				
				fmt.Printf("FOUND: %v \n", found)
				
				settingsAirQualityInfo(found.Lat, found.Lon)
				
			}
			
		}()
		
	}	
	return iup.DEFAULT
}

func settingsAirQualityInfo(lat, lon string) {
	latVal, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		fmt.Printf("Err parsing AQI latitude: %v", err)
		return
	}
	lonVal, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		fmt.Printf("Err parsing AQI longitude: %v", err)
		return
	}
	
	labelLocationAir.SetAttribute("ACTIVE", "NO")
	btnLocationAir.SetAttribute("ACTIVE", "NO")
	listLocations.SetAttribute("VISIBLE", "NO") // fix IUP crash on list fill too fast
	btnRefreshLocations.SetAttribute("ACTIVE", "NO")
	
	defer labelLocationAir.SetAttribute("ACTIVE", "YES")
	defer btnLocationAir.SetAttribute("ACTIVE", "YES")
	defer listLocations.SetAttribute("VISIBLE", "YES")
	defer btnRefreshLocations.SetAttribute("ACTIVE", "YES")
	
	aqiResultTmp.City.Lat = latVal
	aqiResultTmp.City.Lon = lonVal
	aqiResultTmp.City.Idx = 0
	err = aqiResultTmp.getAirQuality()
	
	if err != nil {
		fmt.Printf("Error while getting air quality: %v", err)
		iup.MessageError(wndSettings, "Err while getting air quality... ")
		return
	}
	fmt.Printf("AIR: %v", aqiResultTmp)
	
	/////txt_location_air.SetAttribute("TITLE", aqi_result_tmp.City.Name)
	
	locations, err := getAqiLocations(aqiResultTmp.City.Lat, aqiResultTmp.City.Lon)
	if err != nil {
		fmt.Printf("Error on getting location list: %v", err)
		iup.MessageError(wndSettings, "Error on getting location list... ")
		return
	}
	listLocations.SetAttribute("1", nil) // clear the list
	foundID := 0
	fmt.Printf("locations.Data : [%v]\n", len(locations.Data))
	for i, item := range locations.Data {
		if aqiResultTmp.City.Idx == item.UID {
			foundID = i + 1
		}
		//fmt.Printf("LIST i : [%v][%v][%v]\n", i, item.Station.Name, item.UID)
		listLocations.SetAttribute(fmt.Sprintf("%v", i+1), fmt.Sprintf("%v [%v]", item.Station.Name, item.UID))
		listLocations.SetAttribute(fmt.Sprintf("x%v", i+1), fmt.Sprintf("%v", item.UID))
		listLocations.SetAttribute(fmt.Sprintf("lat_%v", i+1), fmt.Sprintf("%v", item.Lat))
		listLocations.SetAttribute(fmt.Sprintf("lon_%v", i+1), fmt.Sprintf("%v", item.Lon))
		listLocations.SetAttribute(fmt.Sprintf("name_%v", i+1), fmt.Sprintf("%v", item.Station.Name))
	}
	if foundID != 0 {
		listLocations.SetAttribute("VALUE", fmt.Sprintf("%v", foundID))
		labelLocationAir.SetAttribute("TITLE", locations.Data[foundID - 1].Station.Name)
	}
	
}


func onClickSettingsLocationAir(ih iup.Ihandle) int {
	
	return iup.DEFAULT
}

func onClickRefreshLocations(ih iup.Ihandle) int {
	lat := config.Location.Lat_aqi
	lon := config.Location.Lon_aqi
	if lat == "" || lon == "" {
		found, err := weatherTmp.getCurrentLocation()
		if err == nil {
			lat = found.Lat
			lon = found.Lon
		}
	}
	if lat == "" || lon == "" {
		lat = config.Location.Lat
		lon = config.Location.Lon
	}
	go settingsAirQualityInfo(lat, lon)
	return iup.DEFAULT
}

func settingsWindow(show bool) {
	if show {
		fillSettingsWindow()
		wndSettings.SetAttributes("SIMULATEMODAL=YES, BRINGFRONT=YES")
		wndSettings.SetAttribute("TOPMOST", config.AlwaysOnTop)
		mainX, mainY := getWindowXY(dlg)
		iup.ShowXY(wndSettings, max(3, mainX - 70), max(3, mainY + 50))
	} else {
		wndSettings.SetAttribute("SIMULATEMODAL", "NO")
		iup.Hide(wndSettings)
	}
}
