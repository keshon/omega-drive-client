set clientPath=win-client
set serverPath=rcd-server

cd %~dp0%clientPath% && call build-release.bat
cd %~dp0%serverPath% && call build-release.bat
