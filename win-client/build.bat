set "appName=samuno"

rem Create rsrc file
..\3rd-party-tools\rsrc_windows_amd64.exe -manifest %appName%.manifest -o rsrc.syso

rem Rename conf and example files
ren %~dp0src\conf\main.go main
ren %~dp0src\conf\conf_mirball conf_mirball.go

rem Compile as Window app
go build -ldflags="-H=windowsgui -s -w" -o %appName%.exe src/main.go

rem Rename conf and example files
ren %~dp0src\conf\main main.go
ren %~dp0src\conf\conf_mirball.go conf_mirball

rem Set Icon to binary file
..\3rd-party-tools\rcedit-x64.exe "%~dp0%appName%.exe" --set-icon "%~dp0assets\default.ico"