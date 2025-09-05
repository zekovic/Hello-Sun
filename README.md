# Hello Sun

GUI program for viewing local weather and air quality info. 
Created on IUP-GO package for creating multiplatform GUI applications. 

Build for Windows:
```
go build -o hello_sun.exe -ldflags "-w -s -extldflags='-Wl,--allow-multiple-definition' -H windowsgui"
```
##### Preview

![Preview]

[Preview]: screenshots/main.png "Preview"

##### Settings - location

![Settings - location]

[Settings - location]: screenshots/settings_location.png "Settings - location"

##### Settings - display

![Settings - display]

[Settings - display]: screenshots/settings_display.png "Settings - display"

##### Settings - parts

![Settings - parts]

[Settings - parts]: screenshots/settings_parts.png "Settings - parts"

##### Settings - units

![Settings - units]

[Settings - units]: screenshots/settings_units.png "Settings - units"

##### Tray menu

![Tray menu]

[Tray menu]: screenshots/tray_menu.png "Tray menu"

### Which tools are used in this project

##### Weather forecast API:

https://wttr.in

https://github.com/chubin/wttr.in

##### Air Quality API:

https://waqi.info

https://aqicn.org

##### Go bindings for IUP

https://github.com/gen2brain/iup-go

https://www.tecgraf.puc-rio.br/iup/

