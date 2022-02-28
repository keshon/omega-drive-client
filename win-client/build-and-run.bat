set "appName=samuno"

rem Create rsrc file
..\3rd-party-tools\rsrc_windows_amd64.exe -manifest %appName%.manifest -o rsrc.syso

rem Compile as Console app (for better debugging)
go build -o app.exe src/main.go && %appName%.exe

rem Compile as Window app
rem go build -ldflags="-H=windowsgui -s -w" -o %appName%.exe src/main.go

rem Set Icon to binary file
..\3rd-party-tools\rcedit-x64.exe "%~dp0%appName%.exe" --set-icon "%~dp0assets\default.ico"