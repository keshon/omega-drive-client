set clientPath=win-client
cd %~dp0%clientPath% && CALL build-release.bat

set serverPath=rcd-server
cd %~dp0%serverPath% && CALL build-release.bat
