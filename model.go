package main


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

func initImagesInfo() {
	statesMap = make(map[int]int)
	statesTxtMap = make(map[string]int)
	windTxtMap = make(map[string]int)
	moonTxtMap = make(map[string]int)
	
	statesArr[0] =  ImageData{ "Unknown",				"img_sun.jpg", 180,	178,70,	70 }
	statesArr[1] =  ImageData{ "Sunny",				"img_sun.jpg", 1,	1,	70,	70 }
	statesArr[2] =  ImageData{ "Partly Cloudy",		"img_sun.jpg", 91,	1,	70,	70 }
	statesArr[3] =  ImageData{ "Cloudy",				"img_sun.jpg", 180,	1,	70,	70 }
	statesArr[4] =  ImageData{ "Very Cloudy",			"img_sun.jpg", 180,	1,	70,	70 }
	statesArr[5] =  ImageData{ "Fog",					"img_sun.jpg", 270,	1,	70,	70 }
	statesArr[6] =  ImageData{ "Light Showers",		"img_sun.jpg", 2,	89,	70,	70 }
	statesArr[7] =  ImageData{ "Light Sleet Showers",	"img_sun.jpg", 90,	90,	70,	70 }
	statesArr[8] =  ImageData{ "Light Sleet",			"img_sun.jpg", 90,	90,	70,	70 }
	statesArr[9] =  ImageData{ "Thundery Showers",	"img_sun.jpg", 269,	90,	70,	70 }
	statesArr[10] = ImageData{ "Light Snow",			"img_sun.jpg", 3,	177,70,	70 }
	statesArr[11] = ImageData{ "Heavy Snow",			"img_sun.jpg", 92,	176,70,	70 }
	statesArr[12] = ImageData{ "Light Rain",			"img_sun.jpg", 2,	89,	70,	70 }
	statesArr[13] = ImageData{ "Heavy Showers",		"img_sun.jpg", 90,	90,	70,	70 }
	statesArr[14] = ImageData{ "Heavy Rain",			"img_sun.jpg", 90,	90,	70,	70 }
	statesArr[15] = ImageData{ "Light Snow Showers",	"img_sun.jpg", 3,	177,70,	70 }
	statesArr[16] = ImageData{ "Heavy Snow Showers",	"img_sun.jpg", 92,	176,70,	70 }
	statesArr[17] = ImageData{ "Thundery Heavy Rain",	"img_sun.jpg", 180,	89,	70,	70 }
	statesArr[18] = ImageData{ "Thundery Snow Showers","img_sun.jpg",269,	90,	70,	70 }
	
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
	
	statesMap[113] = 1		//"Sunny",
	statesMap[116] = 2		//"PartlyCloudy",
	statesMap[119] = 3		//"Cloudy",
	statesMap[122] = 4		//"VeryCloudy",
	statesMap[143] = 5		//"Fog",
	statesMap[248] = 5		//"Fog",
	statesMap[260] = 5		//"Fog",
	statesMap[263] = 6		//"LightShowers",
	statesMap[353] = 6		//"LightShowers",
	statesMap[176] = 6		//"LightShowers",
	statesMap[179] = 7		//"LightSleetShowers",
	statesMap[374] = 7		//"LightSleetShowers",
	statesMap[362] = 7		//"LightSleetShowers",
	statesMap[365] = 7		//"LightSleetShowers",
	statesMap[281] = 8		//"LightSleet",
	statesMap[284] = 8		//"LightSleet",
	statesMap[311] = 8		//"LightSleet",
	statesMap[314] = 8		//"LightSleet",
	statesMap[317] = 8		//"LightSleet",
	statesMap[350] = 8		//"LightSleet",
	statesMap[377] = 8		//"LightSleet",
	statesMap[182] = 8		//"LightSleet",
	statesMap[185] = 8		//"LightSleet",
	statesMap[386] = 9		//"ThunderyShowers",
	statesMap[200] = 9		//"ThunderyShowers",
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
	
}



type Place struct {
	Name string
	Lat string
	Lon string
}



