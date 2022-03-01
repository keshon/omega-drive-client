@echo off

set "appName=samuno"

echo *** Create rsrc.syso ***
rem Create rsrc file
..\3rd-party-tools\rsrc_windows_amd64.exe -manifest %appName%.manifest -o rsrc.syso

echo *** Swap sample config with actual data ***
rem Rename conf and example files
ren %~dp0src\conf\main.go main
ren %~dp0src\conf\conf_mirball conf_mirball.go

echo *** Compile binaries ***
rem Compile as Window app
go build -ldflags="-H=windowsgui -s -w" -o %appName%.exe src/main.go

echo *** Insert app icon"
rem Set Icon to binary file
..\3rd-party-tools\rcedit-x64.exe "%~dp0%appName%.exe" --set-icon "%~dp0assets\default.ico"

echo *** Swap actual data with samplpe config ***
rem Rename conf and example files
ren %~dp0src\conf\main main.go
ren %~dp0src\conf\conf_mirball.go conf_mirball

echo *** Copy to _release folder ***
rem Copy the necessary files to ../release folder
set rootPath=_release
set srcPath=win-client

rem Create release folder
cd.. && mkdir %rootPath%\

rem Delete old content
del %rootPath%\README.md
del %rootPath%\LICENSE
del %rootPath%\key
del %rootPath%\rsrc.syso
del %rootPath%\%appName%.manifest
del %rootPath%\%appName%.exe
del %rootPath%\logs.txt
del %rootPath%\rcd.exe
del %rootPath%\logs-server.txt
@RD /S /Q "%rootPath%\assets"

rem Copy new content
copy README.md %rootPath%
copy LICENSE %rootPath%
copy NUL %rootPath%\key
copy %srcPath%\rsrc.syso %rootPath%
copy %srcPath%\%appName%.manifest %rootPath%
copy %srcPath%\%appName%.exe %rootPath%
xcopy /e %srcPath%\assets\ %rootPath%\assets\

rem Clearing
del %srcPath%\%appName%.exe
del %srcPath%\logs.txt
