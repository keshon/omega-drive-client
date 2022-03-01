set appName=rcd
set CPATH=%~dp0winfsp\inc\fuse

rem Rename conf and example files
ren %~dp0src\conf\main.go main
ren %~dp0src\conf\conf_mirball conf_mirball.go

rem Compile as Console app
go build -ldflags "-s -w" -tags cmount -o %appName%.exe src/main.go

rem Rename conf and example files
ren %~dp0src\conf\main main.go
ren %~dp0src\conf\conf_mirball.go conf_mirball