package main

import (
	"math/rand/v2"
)


var statesMap map [int]int
var statesTxtMap map [string]int

type ImageData struct {
	Description string
	Image string
	IconX int
	IconY int
	IconW int
	IconH int
}

var statesArr [19]ImageData
var moonArr [8]ImageData
var moonTxtMap map [string]int
var windArr [8]ImageData
var windTxtMap map [string]int
var imagesArr [200]string

func initImagesInfo() {
	statesMap = make(map[int]int)
	statesTxtMap = make(map[string]int)
	windTxtMap = make(map[string]int)
	moonTxtMap = make(map[string]int)
	
	statesArr[0] =  ImageData{ "Unknown",				"", 180,178,70,	70 }
	statesArr[1] =  ImageData{ "Sunny",					"", 1,	1,	70,	70 }
	statesArr[2] =  ImageData{ "Partly Cloudy",			"", 91,	1,	70,	70 }
	statesArr[3] =  ImageData{ "Cloudy",				"", 180,1,	70,	70 }
	statesArr[4] =  ImageData{ "Very Cloudy",			"", 180,1,	70,	70 }
	statesArr[5] =  ImageData{ "Fog",					"", 270,1,	70,	70 }
	statesArr[6] =  ImageData{ "Light Showers",			"", 2,	89,	70,	70 }
	statesArr[7] =  ImageData{ "Light Sleet Showers",	"", 90,	90,	70,	70 }
	statesArr[8] =  ImageData{ "Light Sleet",			"", 90,	90,	70,	70 }
	statesArr[9] =  ImageData{ "Thundery Showers",		"", 269,90,	70,	70 }
	statesArr[10] = ImageData{ "Light Snow",			"", 3,	177,70,	70 }
	statesArr[11] = ImageData{ "Heavy Snow",			"", 92,	176,70,	70 }
	statesArr[12] = ImageData{ "Light Rain",			"", 2,	89,	70,	70 }
	statesArr[13] = ImageData{ "Heavy Showers",			"", 90,	90,	70,	70 }
	statesArr[14] = ImageData{ "Heavy Rain",			"", 90,	90,	70,	70 }
	statesArr[15] = ImageData{ "Light Snow Showers",	"", 3,	177,70,	70 }
	statesArr[16] = ImageData{ "Heavy Snow Showers",	"", 92,	176,70,	70 }
	statesArr[17] = ImageData{ "Thundery Heavy Rain",	"", 180,89,	70,	70 }
	statesArr[18] = ImageData{ "Thundery Snow Showers",	"", 269,90,	70,	70 }
	
	statesTxtMap["?"]		= 0		//  ✨
	statesTxtMap["mm"]		= 3		//  ☁️
	statesTxtMap["="]		= 5		//  🌫
	statesTxtMap["///"]		= 14	//  🌧
	statesTxtMap["//"]		= 13	//  🌧
	statesTxtMap["**"]		= 11	//  ❄️
	statesTxtMap["*/*"]		= 16	//  ❄️
	statesTxtMap["/"]		= 12	//  🌦
	statesTxtMap["."]		= 6		//  🌦
	statesTxtMap["x"]		= 8		//  🌧
	statesTxtMap["x/"]		= 7		//  🌧
	statesTxtMap["*"]		= 10	//  🌨
	statesTxtMap["*/"]		= 15	//  🌨
	statesTxtMap["m"]		= 2		//  ⛅️
	statesTxtMap["o"]		= 1		//  ☀️
	statesTxtMap["/!/"]		= 17	//  🌩
	statesTxtMap["!/"]		= 9		//  ⛈
	statesTxtMap["*!*"]		= 18	//  ⛈
	statesTxtMap["mmm"]		= 4		//  ☁️
	
	statesMap[113] = 1	//"Sunny",
	statesMap[116] = 2	//"PartlyCloudy",
	statesMap[119] = 3	//"Cloudy",
	statesMap[122] = 4	//"VeryCloudy",
	statesMap[143] = 5	//"Fog",
	statesMap[248] = 5	//"Fog",
	statesMap[260] = 5	//"Fog",
	statesMap[263] = 6	//"LightShowers",
	statesMap[353] = 6	//"LightShowers",
	statesMap[176] = 6	//"LightShowers",
	statesMap[179] = 7	//"LightSleetShowers",
	statesMap[374] = 7	//"LightSleetShowers",
	statesMap[362] = 7	//"LightSleetShowers",
	statesMap[365] = 7	//"LightSleetShowers",
	statesMap[281] = 8	//"LightSleet",
	statesMap[284] = 8	//"LightSleet",
	statesMap[311] = 8	//"LightSleet",
	statesMap[314] = 8	//"LightSleet",
	statesMap[317] = 8	//"LightSleet",
	statesMap[350] = 8	//"LightSleet",
	statesMap[377] = 8	//"LightSleet",
	statesMap[182] = 8	//"LightSleet",
	statesMap[185] = 8	//"LightSleet",
	statesMap[386] = 9	//"ThunderyShowers",
	statesMap[200] = 9	//"ThunderyShowers",
	statesMap[227] = 10	//"LightSnow",
	statesMap[320] = 10	//"LightSnow",
	statesMap[230] = 11	//"HeavySnow",
	statesMap[329] = 11	//"HeavySnow",
	statesMap[332] = 11	//"HeavySnow",
	statesMap[338] = 11	//"HeavySnow",
	statesMap[266] = 12	//"LightRain",
	statesMap[293] = 12	//"LightRain",
	statesMap[296] = 12	//"LightRain",
	statesMap[299] = 13	//"HeavyShowers",
	statesMap[305] = 13	//"HeavyShowers",
	statesMap[356] = 13	//"HeavyShowers",
	statesMap[302] = 14	//"HeavyRain",
	statesMap[308] = 14	//"HeavyRain",
	statesMap[359] = 14	//"HeavyRain",
	statesMap[323] = 15	//"LightSnowShowers",
	statesMap[326] = 15	//"LightSnowShowers",
	statesMap[368] = 15	//"LightSnowShowers",
	statesMap[335] = 16	//"HeavySnowShowers",
	statesMap[371] = 16	//"HeavySnowShowers",
	statesMap[395] = 16	//"HeavySnowShowers",
	statesMap[389] = 17	//"ThunderyHeavyRain",
	statesMap[392] = 18	//"ThunderySnowShowers",
	
	
	moonArr[0] =  ImageData{ "",	"", 2,	265,70,	70 }
	moonArr[1] =  ImageData{ "",	"", 90,	265,70,	70 }
	moonArr[2] =  ImageData{ "",	"", 180,266,70,	70 }
	moonArr[3] =  ImageData{ "",	"", 269,266,70,	70 }
	moonArr[4] =  ImageData{ "",	"", 2,	353,70,	70 }
	moonArr[5] =  ImageData{ "",	"", 90,	354,70,	70 }
	moonArr[6] =  ImageData{ "",	"", 179,353,70,	70 }
	moonArr[7] =  ImageData{ "",	"", 269,354,70,	70 }
	
	moonTxtMap["🌑"] = 0
	moonTxtMap["🌒"] = 1
	moonTxtMap["🌓"] = 2
	moonTxtMap["🌔"] = 3
	moonTxtMap["🌕"] = 4
	moonTxtMap["🌖"] = 5
	moonTxtMap["🌗"] = 6
	moonTxtMap["🌘"] = 7
	
	windArr[0] =  ImageData{ "",	"", 420,87,	70,	70 }
	windArr[1] =  ImageData{ "",	"", 423,177,70,	70 }
	windArr[2] =  ImageData{ "",	"", 350,1,	70,	70 }
	windArr[3] =  ImageData{ "",	"", 352,266,70,	70 }
	windArr[4] =  ImageData{ "",	"", 350,87,	70,	70 }
	windArr[5] =  ImageData{ "",	"", 422,266,70,	70 }
	windArr[6] =  ImageData{ "",	"", 420,1,	70,	70 }
	windArr[7] =  ImageData{ "",	"", 353,177,70,	70 }
	
	windTxtMap["↓"] = 0
	windTxtMap["↙"] = 1
	windTxtMap["←"] = 2
	windTxtMap["↖"] = 3
	windTxtMap["↑"] = 4
	windTxtMap["↗"] = 5
	windTxtMap["→"] = 6
	windTxtMap["↘"] = 7
	
	imagesArr[0] = "1_1_clouds-1117586_1280.jpg"
	// Sunny
	imagesArr[10] = "1_1_clouds-1117586_1280.jpg"
	imagesArr[11] = "1_1_field-3629120_1920.jpg"
	imagesArr[12] = "1_1_rapeseeds-474558_1280.jpg"
	imagesArr[13] = "1_1_trees-5033072_1280.jpg"
	imagesArr[14] = "1_2_sunset-789974_1280.jpg"
	// Partly Cloudy
	imagesArr[20] = "2_0_sunrise-1513802_1280.jpg"
	imagesArr[21] = "2_1_chris-von-krebs-cintorino-1RAiGhaaR1c-unsplash.jpg"
	imagesArr[22] = "2_1_clouds-2085112_1280.jpg"
	imagesArr[23] = "2_1_darling-7655568_1280.jpg"
	imagesArr[24] = "2_1_heaven-3395811_1280.jpg"
	imagesArr[25] = "2_1_sky-7456744_1280.jpg"
	imagesArr[26] = "2_2_purple-669046_1280.jpg"
	imagesArr[27] = "2_2_sunset-1661088_1280.jpg"
	// Cloudy
	imagesArr[30] = "3_1_sky-414199_1280.jpg"
	// Very Cloudy
	imagesArr[40] = "4_1_clouds-8029036_1280.jpg"
	imagesArr[41] = "4_1_sea-6811812_1280.jpg"
	// Fog
	imagesArr[50] = "5_1_clouds-4979558_1280.jpg"
	imagesArr[51] = "5_1_fog-4436636_1280.jpg"
	// Light Showers
	imagesArr[60] = "6_1_glass-window-1845534_1280.jpg"
	imagesArr[61] = "6_1_rain-122691_1280.jpg"
	// Light Sleet Showers
	imagesArr[70] = "8_1_ice-crystals-6939641_1280.jpg"
	imagesArr[71] = "8_1_winter-3183033_1280.jpg"
	// Light Sleet
	imagesArr[80] = "8_1_ice-crystals-6939641_1280.jpg"
	imagesArr[81] = "8_1_winter-3183033_1280.jpg"
	// Thundery Showers
	imagesArr[90] = "9_1_lightning-2702168_1280.jpg"
	imagesArr[91] = "9_1_lightning-4304449_1280.jpg"
	// Light Snow
	imagesArr[100] = "10_1_snow-1768544_1280.jpg"
	imagesArr[101] = "10_1_winter-7661769_1280.jpg"
	// Heavy Snow
	imagesArr[110] = "10_1_snow-1768544_1280.jpg"
	imagesArr[111] = "10_1_winter-7661769_1280.jpg"
	// Light Rain
	imagesArr[120] = "6_1_glass-window-1845534_1280.jpg"
	imagesArr[121] = "6_1_rain-122691_1280.jpg"
	// Heavy Showers
	imagesArr[130] = "6_1_glass-window-1845534_1280.jpg"
	imagesArr[131] = "6_1_rain-122691_1280.jpg"
	// Heavy Rain
	imagesArr[140] = "6_1_glass-window-1845534_1280.jpg"
	imagesArr[141] = "6_1_rain-122691_1280.jpg"
	// Light Snow Showers
	imagesArr[150] = "10_1_snow-1768544_1280.jpg"
	imagesArr[151] = "10_1_winter-7661769_1280.jpg"
	// Heavy Snow Showers
	imagesArr[160] = "10_1_snow-1768544_1280.jpg"
	imagesArr[161] = "10_1_winter-7661769_1280.jpg"
	// Thundery Heavy Rain
	imagesArr[170] = "9_1_lightning-2702168_1280.jpg"
	imagesArr[171] = "9_1_lightning-4304449_1280.jpg"
	// Thundery Snow Showers
	imagesArr[180] = "10_1_snow-1768544_1280.jpg"
	imagesArr[181] = "10_1_winter-7661769_1280.jpg"

}



func getImageOfWeather(state int) string {
	randomCount := 1
	if state == 1 { randomCount = 5 }
	if state == 2 { randomCount = 8 }
	if state == 3 { randomCount = 1 }
	if state == 4 { randomCount = 2 }
	if state == 5 { randomCount = 2 }
	if state == 6 { randomCount = 2 }
	if state == 7 { randomCount = 2 }
	if state == 8 { randomCount = 2 }
	if state == 9 { randomCount = 2 }
	if state == 10 { randomCount = 2 }
	if state == 11 { randomCount = 2 }
	if state == 12 { randomCount = 2 }
	if state == 13 { randomCount = 2 }
	if state == 14 { randomCount = 2 }
	if state == 15 { randomCount = 2 }
	if state == 16 { randomCount = 2 }
	if state == 17 { randomCount = 2 }
	if state == 18 { randomCount = 2 }
	
	randomItem := rand.IntN(randomCount) + (state * 10)
	imageName := imagesArr[randomItem]
	if imageName == "" {
		imageName = imagesArr[0]
	}
	//fmt.Printf("RANDOM: %v, %v, %v", randomCount, randomItem, imageName)
	return imageName
}


type Place struct {
	Name string
	Lat string
	Lon string
}



