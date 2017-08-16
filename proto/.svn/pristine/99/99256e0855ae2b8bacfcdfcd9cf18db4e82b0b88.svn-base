@echo off

if not exist pb md pb

setlocal enabledelayedexpansion

for %%i in (*.proto) do (
	set prtofile=%%i
	set pbfile=!prtofile:proto=pb!
	protoc -o pb/!pbfile! !prtofile!
	echo convert [!prtofile!] ok
)

xcopy /y pb\* ..\cocos2d-x-2.2.6\projects\Casino\Resources\pb\
pause