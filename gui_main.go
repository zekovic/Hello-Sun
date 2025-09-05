package main

import (
	"embed"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"time"
	"unicode/utf8"

	"github.com/gen2brain/iup-go/iup"
)

var dlg iup.Ihandle
var cv iup.Ihandle
var trayMenu iup.Ihandle
var dlgW = 280
var dlgH = 500
var dlgOpacity = 255
var dlgShape *image.RGBA
var mouseHere bool
var mouseDown bool = false
var mouseX int
var mouseY int
var currentStateInt = -1

var theTimer iup.Ihandle

//go:embed res/*
var res_folder embed.FS

func showGui() {
	iup.SetGlobal("UTF8MODE", "YES")
	iup.Open()
	defer iup.Close()
	
	getBackgroundImage("1_1_img_sun.jpg")
	
	file, err := res_folder.Open("res/menu.png")
	if err != nil {
		fmt.Println(err)
	}
	menuImg, err := png.Decode(file)
	if err != nil {
		fmt.Println(err)
	}
	menuImgClip := menuImg.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(0, 0, dlgW, 30))
	iup.ImageFromImage(menuImgClip).SetHandle("menu_img_clip")
	
	
	
	file, err = res_folder.Open("res/waqi.info_aqicn.org_logo.png")
	if err != nil {
		fmt.Println(err)
	}
	aqicnLogo, err := png.Decode(file)
	if err != nil {
		fmt.Println(err)
	}
	iup.ImageFromImage(aqicnLogo).SetHandle("aqicn_logo")
	
	file, err = res_folder.Open("res/wttr.in_logo.jpg")
	if err != nil {
		fmt.Println(err)
	}
	wttrLogo, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println(err)
	}
	iup.ImageFromImage(wttrLogo).SetHandle("wttr_logo")
	
	file.Close()
	
	
	loadSubImageArrays()
	
	cv = iup.Canvas().SetAttribute("USERSIZE", fmt.Sprintf("%vx%v", dlgW, dlgH))
	cv.SetCallback("ACTION", iup.CanvasActionFunc(actionCb))
	
	dlg = iup.Dialog(cv,);
	
	dlg.SetAttributes(`TITLE="Hello Sun",RXESIZE=NO, CUSTOMFRAMESIMULATE=YES, 
		TRAY=YES, TRAYIMAGE=state_1, TRAYTIP="Hello Sun"`)
	dlg.SetAttribute("TOPMOST", config.AlwaysOnTop)
	dlg.SetAttribute("TASKBARBUTTON", "HIDE")
	//dlg_opacity = config.getOpacity()
	dlg.SetAttribute("OPACITY", config.Opacity)
	fmt.Printf("OPACITY....... [%v]", config.Opacity)
	
	cv.SetCallback("MOTION_CB", iup.MotionFunc(motionCb))
	cv.SetCallback("BUTTON_CB", iup.ButtonFunc(onClick))
	cv.SetCallback("ENTERWINDOW_CB", iup.EnterWindowFunc(onMyMouseEnter))
	cv.SetCallback("LEAVEWINDOW_CB", iup.LeaveWindowFunc(onMyMouseLeave))
	
	dlg.SetCallback("SHOW_CB", iup.ShowFunc(onMainDialogShow))
	
	dlg.SetCallback("TRAYCLICK_CB", iup.TrayClickFunc(onMainDialogTrayClick))
	
	// dlg.SetAttribute("SHAPEIMAGE", "shape")
	updateWindowShape(dlgW, dlgH)
	
	theTimer = iup.Timer()
	theTimer.SetAttribute("TIME", config.Refresh * 60 * 1000)
	// the_timer.SetAttribute("TIME", 10000)
	theTimer.SetCallback("ACTION_CB", iup.TimerActionFunc(timerTick))
	
	createSettingsDialog()
	
	trayMenu = iup.Menu(
		iup.Item("&Restore window").SetCallback("ACTION", iup.ActionFunc(restoreWindow)),
		iup.Item("&Hide window").SetCallback("ACTION", iup.ActionFunc(hideWindow)),
		iup.Item("&Settings").SetCallback("ACTION", iup.ActionFunc(openSettings)),
		iup.Item("&Close").SetCallback("ACTION", iup.ActionFunc(closeApp)),
	)
	
	initDrawFunctions()
	
	xDlg, yDlg := loadWindowPosition()
	iup.ShowXY(dlg, xDlg, yDlg)
	
	timerTick(theTimer)
	theTimer.SetAttribute("RUN", "YES")
	
	iup.MainLoop()
}

func restoreWindow(ih iup.Ihandle) int {
	if dlg.GetAttribute("IS_MINIMIZED") == "YES" {
		dlg.SetAttribute("IS_MINIMIZED", "NO")
		dlg.SetAttribute("PLACEMENT", "NORMAL")
		iup.Show(dlg)
	}
	return iup.DEFAULT
}
func hideWindow(ih iup.Ihandle) int {
	dlg.SetAttribute("PLACEMENT", "MINIMIZED")
	dlg.SetAttribute("IS_MINIMIZED", "YES")
	iup.Show(dlg)
	return iup.MINIMIZE
}
func openSettings(ih iup.Ihandle) int {
	settingsWindow(true)
	return iup.DEFAULT
}
func closeApp(ih iup.Ihandle) int {
	setConfigWindowPosition()
	config.save()
	iup.ExitLoop()
	return iup.DEFAULT
}

func updateWindowShape(w, h int) {
	dlgShape = createRoundedImage(w, h, 20)
	iup.ImageFromImage(dlgShape).SetHandle("shape").SetAttribute("RESIZE", fmt.Sprintf("%vx%v", w, h))
	dlg.SetAttribute("SHAPEIMAGE", "shape")
	dlg.SetAttribute("USERSIZE", fmt.Sprintf("%vx%v", w, h))
	cv.SetAttribute("USERSIZE", fmt.Sprintf("%vx%v", w, h))
	
	fmt.Printf("SIZE......: [%v], TOTAL_Y:  [%v]", dlg.GetAttribute("USERSIZE"), h)
}

func timerTick(ih iup.Ihandle) int {
	go fetchData()
	fmt.Println("TICK ")
	return iup.DEFAULT
}

func loadSubImage(img image.Image, x, y, w, h int, nameId string) {
	subImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(x, y, x + w, y + h))
	iup.ImageFromImage(subImg).SetHandle(nameId)
}

func loadSubImageArrays() {
	file, err := res_folder.Open("res/weather_icons.png")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	pngImage, err := png.Decode(file)
	if err != nil {
		fmt.Println(err)
	}
	
	for i, img := range windArr {
		loadSubImage(pngImage, img.IconX, img.IconY, img.IconW, img.IconH, fmt.Sprintf("wind_%v", i))
	}
	for i, img := range moonArr {
		loadSubImage(pngImage, img.IconX, img.IconY, img.IconW, img.IconH, fmt.Sprintf("moon_%v", i))
	}
	for i, img := range statesArr {
		loadSubImage(pngImage, img.IconX, img.IconY, img.IconW, img.IconH, fmt.Sprintf("state_%v", i))
	}
	
	sun := statesArr[statesTxtMap["o"]]
	loadSubImage(pngImage, sun.IconX, sun.IconY, sun.IconW, int(float32(sun.IconH) / 2.5 ), "sunrise_img")
}


func onMainDialogShow(ih iup.Ihandle, state int) int {
	
	//fmt.Printf("Show STATE: [%v]\n", state) // 0 open, 1 restore, 2 minimize
	if state == 0 && !mouseDown {
		
		//dlg_opacity = config.getOpacity()
		dlg.SetAttribute("OPACITY", config.Opacity)
		fmt.Printf("OPACITY... [%v]", config.Opacity)
		
		dlg.SetAttribute("IS_MINIMIZED", "NO")
		
	}
	return iup.DEFAULT
}

func fetchData() {
	weatherResult.getAndParse("")
	if config.Location.Weather == "" {
		weatherResult.saveLocationToConfig()
		config.save()
	}
	getInTextFormat()
	
	err := aqiResult.getAirQuality()
	fmt.Printf("AIR: %v", aqiResult)
	
	if len(weatherResult.CurrentCondition) > 0 && err == nil {
		iup.Update(cv)
	}
}

var paintMenuStart int = -1
var paintMenuEnd int = -1

func motionCb(ih iup.Ihandle, x, y int, status string) int {
	
	if iup.IsButton1(status) {
		diffX := x - mouseX
		diffY := y - mouseY
		wndX, wndY := getWindowXY(dlg)
		iup.ShowXY(dlg, wndX + diffX, wndY + diffY)
		// dlg.SetAttributes(`OPACITY=200`)
		dlg.SetAttribute("OPACITY", config.Opacity)
		//dlg.SetAttribute("ACTIVEWINDOW", "YES")
		dlg.SetAttribute("BRINGFRONT", "YES")
	} else {
		dlg.SetAttributes(`OPACITY=255`)
		mouseX = x
		mouseY = y
		paintMenuStart = -1
		paintMenuEnd = -1
		if y >= dlgH - 30 {
			//if x < 265 {
			if x < dlgW {
				paintMenuStart = 206
				// paint_menu_end = 265
				paintMenuEnd = dlgW
			}
			if x < 205 {
				paintMenuStart = 111
				paintMenuEnd = 210
			}
			if x < 111 {
				paintMenuStart = 0
				paintMenuEnd = 110
			}
		}
	}
	
	return iup.DEFAULT
}

func onMyMouseEnter(ih iup.Ihandle) int {
	mouseHere = true
	mouseDown = false
	dlg.SetAttributes(`OPACITY=255`)
	//fmt.Print("a")
	iup.Update(ih)
	return iup.DEFAULT
}
func onMyMouseLeave(ih iup.Ihandle) int {
	mouseHere = false
	mouseDown = false
	dlg.SetAttribute("OPACITY", config.Opacity)
	//dlg.SetAttributes(`OPACITY=200`)
	//fmt.Print("b")
	iup.Update(ih)
	return iup.DEFAULT
}

func loadWindowPosition() (int, int) {
	xScr, yScr := getScreenSize()
	x := int(config.Position.X / 100.0 * float64(xScr))
	y := int(config.Position.Y / 100.0 * float64(yScr))
	return x, y
}

func setConfigWindowPosition() {
	xScr, yScr := getScreenSize()
	xDlg, yDlg := getWindowXY(dlg)
	if config.Position.X == 0.0 { config.Position.X = 30.0}
	if config.Position.Y == 0.0 { config.Position.Y = 30.0}
	// config.Position.X = int(float32(x_dlg) / float32(x_scr) * 100.0)
	// config.Position.Y = int(float32(y_dlg) / float32(y_scr) * 100.0)
	config.Position.X = float64(xDlg) / float64(xScr) * 100.0
	config.Position.Y = float64(yDlg) / float64(yScr) * 100.0
}

func onClick(ih iup.Ihandle, button, pressed, x, y int, status string) int {
	
	//btn_int, err := strconv.Atoi(strings.ReplaceAll(status, " ", ""))
	/*if status == iup.BUTTON1 { }*/
	if iup.IsButton1(status) {
		mouseDown = true
	} else {
		mouseDown = false
	}
	if y > dlgH - 30 && pressed == 1 && iup.IsButton1(status) {
		
		if x < 111 {
			openSettings(ih)
		}
		if x >= 111 && x < 205 {
			hideWindow(ih)
		}
		if x >= 205 /*&& x < 265*/ {
			if config.Systray == "ON" {
				hideWindow(ih)
			} else {
				closeApp(ih)
			}
		}
		
	}
	return iup.DEFAULT
}

func onMainDialogTrayClick(ih iup.Ihandle, button, pressed, dblclick int) int {
	fmt.Printf("BTN: [%v] RESTORE..... [%v] \n", button, dlg.GetAttribute("IS_MINIMIZED"))
	if button == 3 {
		xScr, yScr := getMousePosition()
		iup.Popup(trayMenu, xScr - 10, yScr - 15)
	} else {
		restoreWindow(ih)
	}
	return iup.DEFAULT
}

var totalY int
var totalYOld int


func getBackgroundImage(fileName string) {
	file, err := res_folder.Open(fmt.Sprintf("res/%v", fileName))
	if err != nil {
		fmt.Println(err)
	}

	jpgImage, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println(err)
	}
	
	mySubImage := jpgImage.(interface {
		SubImage(r image.Rectangle) image.Image
	// }).SubImage(image.Rect(180, 20, 180 + dlgW, 20 + dlgH))
	}).SubImage(image.Rect(1, 1, 1 + dlgW, 1 + dlgH))
	iup.ImageFromImage(mySubImage).SetHandle("background_image")
}

func actionCb(ih iup.Ihandle, posx, posy float64) int {
	iup.DrawBegin(ih)
	
	// w, h := iup.DrawGetSize(ih)
	w, h := dlgW, dlgH
	
	totalY = 0
	
	if len(weatherResult.CurrentCondition) > 0 {
		
		stateInt := statesTxtMap[weatherResult.CurrentCondition[0].WeatherCodeTxt]
		if stateInt != currentStateInt {
			currentStateInt = stateInt
			getBackgroundImage(getImageOfWeather(stateInt))
		}
		
		iup.DrawImage(ih, "background_image", 0, 0, w, h)
		// weather_info := weather_result.CurrentCondition[0]
		
		totalY += 25
		drawMap["main"]()
		for _, part := range config.ShowParts {
			drawMap[part]()
		}
		drawMap["logo"]()
		
		if totalYOld != totalY {
			fmt.Printf("Needs to be changed......., old:[%v], new:[%v]", totalYOld, totalY)
			dlgH = totalY // + 20
			updateWindowShape(dlgW, dlgH)
		}
		
		totalYOld = totalY
	}
	
	if mouseHere {
		//ih.SetAttributes(`DRAWCOLOR="50 50 50", DRAWFONT="Tahoma, Bold 12"`)
		//iup.DrawText(ih, "Settings   X", w-100, h-30, -1, -1)
		
		iup.DrawImage(ih, "menu_img_clip", 0, h-30, w, 30)
		if paintMenuStart >= 0 {
			ih.SetAttributes(`DRAWCOLOR="90 90 90 150", DRAWSTYLE=FILL`)
			iup.DrawRectangle(ih, paintMenuStart, h - 30, paintMenuEnd, h)
			// fmt.Printf("C")
		}
		// fmt.Printf("B")
	}
	// fmt.Printf("A")
	iup.DrawEnd(ih)
	return iup.DEFAULT
}

var drawMap map [string]func()

func initDrawFunctions() {
	drawMap = make(map [string]func())
	
	drawMap["main"] = func() {
		weather_info := weatherResult.CurrentCondition[0]
		stateInt := statesTxtMap[weather_info.WeatherCodeTxt]
		
		if weatherResult.isNight() {
			moonInt := moonTxtMap[weather_info.MoonIcon]
			iup.DrawImage(cv, fmt.Sprintf("moon_%v", moonInt), dlgW - 110, totalY + 15, 70, 70)
			overMoonInt := getStateOverMoon(stateInt)
			if overMoonInt != -1 {
				iup.DrawImage(cv, fmt.Sprintf("state_%v", overMoonInt), dlgW - 95, totalY + 35, 70, 70)
			}
		} else {
			iup.DrawImage(cv, fmt.Sprintf("state_%v", stateInt), dlgW - 110, totalY + 15, 70, 70)
		}
		
		
		t := getTempByUnit()
		
		cityFont := 30
		cityLen := utf8.RuneCountInString(config.Location.Weather)
		//fmt.Printf("CITY RUNE COUNT: [%v]", cityLen)
		if cityLen > 10 {
			cityFont = 25
		}
		if cityLen > 15 {
			cityFont = 15
		}
		
		cv.SetAttribute("DRAWFONT", fmt.Sprintf("%v, Bold %v", config.Font, cityFont))
		myOutlineText(cv, fmt.Sprintf("%v", config.Location.Weather), 20, totalY, -1, -1, 2, "255 255 255", "25 25 25")
		cv.SetAttribute("DRAWFONT", config.Font + ", Bold 30")
		myOutlineText(cv, fmt.Sprintf("%v", t), 20, totalY+45, -1, -1, 2, "255 255 255", "25 25 25")
		totalY += 125
	}
	
	drawMap["logo"] = func() {
		logoX := 80
		cv.SetAttributes(`DRAWCOLOR="90 90 90 150", DRAWSTYLE=FILL`)
		DrawRoundedRect(cv, logoX, totalY, dlgW - 5, totalY + 25, 6)
		// iup.DrawRectangle(cv, logo_x, pos_y,   dlg_w-30,  pos_y+25)
		cv.SetAttribute("DRAWFONT", config.Font + ", 10")
		iup.DrawImage(cv, "aqicn_logo", logoX + 10, totalY+3, 20, 20)
		myOutlineText(cv, "aqicn.org", logoX + 35, totalY+5, -1, -1, 1, "255 255 255", "25 25 25")
		iup.DrawImage(cv, "wttr_logo", logoX + dlgW/2-20, totalY+3, 20, 20)
		myOutlineText(cv, "wttr.in", logoX + dlgW/2 + 5, totalY+5, -1, -1, 1, "255 255 255", "25 25 25")
		totalY += 45
	}
	
	drawMap["daily"] = func() {
		cv.SetAttribute("DRAWFONT", config.Font + ", 11")
		cv.SetAttributes(`DRAWCOLOR="90 90 90 150", DRAWSTYLE=FILL`)
		DrawRoundedRect(cv, 5, totalY, dlgW - 5, totalY + 80, 6)
		for i, data := range weatherResult.Weather {
			if i > 3 {
				break
			}
			tMin, tMax, tAvg, tUnit := TempDayInfo(data)
			cv.SetAttribute("DRAWFONT", config.Font + ", 11")
			time_day, _ := time.Parse("2006-01-02", data.Date)
			myOutlineText(cv, fmt.Sprintf("%v", time_day.Format("Mon 02")), 20 + (i*86), totalY+5, -1, -1, 1, "255 255 255", "25 25 25")
			myOutlineText(cv, fmt.Sprintf("%v - %v", tMin, tMax), 20 + (i*86), totalY+56, -1, -1, 1, "255 255 255", "25 25 25")
			
			cv.SetAttribute("DRAWFONT", config.Font + ", 24")
			myOutlineText(cv, fmt.Sprintf("%v%v", tAvg, tUnit), 
				20 + (i*86), totalY + 20, -1, -1, 1, "255 255 255", "25 25 25")
			
		}
		totalY += 90
	}
	
	drawMap["hourly"] = func() {
		
	}
	
	drawMap["air"] = func() {
		cv.SetAttribute("DRAWFONT", config.Font + ", Bold 14")
		aqiPM25 := aqiResult.Values.PM25
		aqiPM10 := aqiResult.Values.PM10
		color25, color10 := getAqiColor(aqiPM25, aqiPM10)
		cv.SetAttributes(`DRAWCOLOR="90 90 90 150", DRAWSTYLE=FILL`)
		cv.SetAttributes(fmt.Sprintf(`DRAWCOLOR="%v 150", DRAWSTYLE=FILL`, color25))
		DrawRoundedRect(cv, 5, totalY, dlgW/2 - 5, totalY + 32, 6)
		cv.SetAttributes(fmt.Sprintf(`DRAWCOLOR="%v 150", DRAWSTYLE=FILL`, color10))
		DrawRoundedRect(cv, dlgW/2 + 2, totalY, dlgW - 5, totalY + 32, 6)
		myOutlineText(cv, fmt.Sprintf("PM 2.5: %v", aqiPM25), 10, totalY+5, -1, -1, 1, "255 255 255", "25 25 25")
		myOutlineText(cv, fmt.Sprintf("PM 10: %v", aqiPM10), dlgW/2+10, totalY+5, -1, -1, 1, "255 255 255", "25 25 25")
		totalY += 50
	}
	
	drawMap["details"] = func() {
		cv.SetAttributes(`DRAWCOLOR="90 90 90 150", DRAWSTYLE=FILL`)
		DrawRoundedRect(cv, 5, totalY, dlgW - 5, totalY + 70, 6)
		cv.SetAttribute("DRAWCOLOR", "255 255 255")
		cv.SetAttribute("DRAWFONT", config.Font + ", 12")
		iup.DrawText(cv, "Pressure: "+getPressureByUnit(), 10, totalY+5, -1, -1)
		iup.DrawText(cv, "UV: "+getUV(), dlgW-65, totalY+5, -1, -1)
		iup.DrawText(cv, "Rain: "+getPrecipitationByUnit(), 10, totalY + 25, -1, -1)
		
		windValue, windArrow := getWindByUnit()
		wind_pos := dlgW/2-10
		iup.DrawText(cv, "Wind:", wind_pos, totalY + 25, -1, -1)
		wind_w, _ := iup.DrawGetTextSize(cv, "Wind:")
		
		if windArrow != "" {
			wind_int, ok := windTxtMap[weatherResult.CurrentCondition[0].WinddirArrow]
			if ok {
				iup.DrawImage(cv, fmt.Sprintf("wind_%v", wind_int), wind_pos + wind_w, totalY + 25, 20, 20)
				wind_w += 23
			}
		}
		
		iup.DrawText(cv, windValue, wind_pos + wind_w, totalY + 25, -1, -1)
		
		iup.DrawText(cv, "Humidity: " + weatherResult.CurrentCondition[0].Humidity, 10, totalY + 45, -1, -1)
		totalY += 75
	}
	
	drawMap["moon"] = func() {
		weather_info := weatherResult.CurrentCondition[0]
		cv.SetAttribute("DRAWFONT", config.Font + ", Bold 12")
		moon_int, ok := moonTxtMap[weather_info.MoonIcon]
		if ok {
			iup.DrawImage(cv, fmt.Sprintf("moon_%v", moon_int), 20, totalY, 20, 20)
		}
		myOutlineText(cv, fmt.Sprintf("Moon age: %v", weather_info.MoonDay), 50, totalY, -1, -1, 1, "255 255 255", "25 25 25")
		totalY += 30
	}
	
	drawMap["sunrise"] = func() {
		weather_info := weatherResult.CurrentCondition[0]
		cv.SetAttribute("DRAWFONT", config.Font + ", Bold 12")
		
		sunrise_str := weather_info.Sunrise
		sunset_str := weather_info.Sunset
		if config.Time == "AMPM" {
			sunrise_str = time24toAMPM(sunrise_str)
			sunset_str = time24toAMPM(sunset_str)
		}
		
		iup.DrawImage(cv, "sunrise_img", 18, totalY + 5, 28, 10)
		iup.DrawImage(cv, "wind_4", 35, totalY, 20, 20)
		myOutlineText(cv, sunrise_str, 58, totalY, -1, -1, 1, "255 255 255", "25 25 25")
		iup.DrawImage(cv, "sunrise_img", 140, totalY + 5, 28, 10)
		iup.DrawImage(cv, "wind_0", 157, totalY, 20, 20)
		myOutlineText(cv, sunset_str, 180, totalY, -1, -1, 1, "255 255 255", "25 25 25")
		totalY += 30
	}
	
}

func myOutlineText(ih iup.Ihandle, txt string, x, y, w, h, outline int, color_in, color_out string) {
	ih.SetAttribute("DRAWCOLOR", color_out)
	iup.DrawText(ih, txt, x-outline, y-outline, w, h)
	iup.DrawText(ih, txt, x+outline, y+outline, w, h)
	iup.DrawText(ih, txt, x-outline, y+outline, w, h)
	iup.DrawText(ih, txt, x+outline, y-outline, w, h)
	ih.SetAttribute("DRAWCOLOR", color_in)
	iup.DrawText(ih, txt, x, y, w, h)
}

func DrawRoundedRect(ih iup.Ihandle, x1, y1, x2, y2, r int) {
	R := r * 2
	iup.DrawRectangle(ih,   x1,          y1 + r + 1,   x2,               y2 - r - 1)
	iup.DrawRectangle(ih,   x1 + r + 1,  y1,           x2 - r - 1,       y1 + r)
	iup.DrawRectangle(ih,   x1 + r + 1,  y2 - r,       x2 - r - 1,       y2)
	iup.DrawArc(ih,         x1,          y1,           x1 + (R + 1),     y1 + (R + 1),   90.0, 180.0)
	iup.DrawArc(ih,         x2,          y1,           x2 - (R + 1),     y1 + (R + 1),   0.0, 90.0)
	iup.DrawArc(ih,         x1,          y2,           x1 + (R + 1),     y2 - (R + 1),   180.0, 270.0)
	iup.DrawArc(ih,         x2,          y2,           x2 - (R + 1),     y2 - (R + 1),   270.0, 360.0)
}

