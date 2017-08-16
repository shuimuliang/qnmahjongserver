call build.bat

@set outPath=Proto_t.lua
@set inPath=pf.proto

@ECHO É¾³ý¾ÉµÄ
del %outPath%

lua runMe_.lua %outPath% %inPath%

xcopy /y Proto_t.lua ..\cocos2d-x-2.2.6\projects\Casino\Resources\lua\network\
xcopy /y Proto_t.lua ..\cocos2d-x-2.2.6\projects\Casino\Resources_ar\lua\network\
xcopy /y Proto_t.lua ..\cocos2d-x-2.2.6\projects\Casino\Resources_en\lua\network\
xcopy /y Proto_t.lua ..\cocos2d-x-2.2.6\projects\Casino\Resources_ma\lua\network\
xcopy /y Proto_t.lua ..\cocos2d-x-2.2.6\projects\Casino\Resources_tw\lua\network\



pause