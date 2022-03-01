@echo off

set appName=rcd
set CPATH=%~dp0winfsp\inc\fuse

echo *** Swap sample config with actual data ***
rem Rename conf and example files
ren %~dp0src\conf\main.go main
ren %~dp0src\conf\conf_mirball conf_mirball.go

echo *** Compile binaries ***
rem Compile as Console app
go build -ldflags "-s -w" -tags cmount -o %appName%.exe src/main.go

echo *** Swap actual data with samplpe config ***
rem Rename conf and example files
ren %~dp0src\conf\main main.go
ren %~dp0src\conf\conf_mirball.go conf_mirball

echo *** Copy to _release folder ***
rem Copy the necessary files to ../release folder
set rootPath=_release
set srcPath=rcd-server

rem Create release folder
cd.. && mkdir %rootPath%\

rem Delete old content
del %rootPath%\%appName%.exe

rem Copy new content
copy %srcPath%\%appName%.exe %rootPath%

rem Clearing
del %srcPath%\%appName%.exe
del %srcPath%\logs-server.txt