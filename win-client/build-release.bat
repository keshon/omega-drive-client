set "appName=samuno"

rem Create rsrc file
..\3rd-party-tools\rsrc_windows_amd64.exe -manifest %appName%.manifest -o rsrc.syso

rem Compile as Window app
go build -ldflags="-H=windowsgui -s -w" -o %appName%.exe src/main.go

rem Set Icon to binary file
..\3rd-party-tools\rcedit-x64.exe "%~dp0%appName%.exe" --set-icon "%~dp0assets\default.ico"

rem Copy the necessary files to ../release folder
set rootPath=_release
set clientPath=win-client

cd .. && @RD /S /Q %rootPath%

mkdir %rootPath%\
copy README.md %rootPath%
copy LICENSE %rootPath%
copy NUL %rootPath%\key
copy %clientPath%\rsrc.syso %rootPath%
copy %clientPath%\%appName%.exe %rootPath%
xcopy /e %clientPath%\assets\ %rootPath%\assets\