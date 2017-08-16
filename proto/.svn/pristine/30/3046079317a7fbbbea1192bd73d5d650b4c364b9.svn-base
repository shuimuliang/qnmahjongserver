#!/bin/sh

outPath=Proto_pf.lua
inPath=pf.proto

echo “delete —— ${outPath} ${inPath}”
rm $outPath

lua runMe.lua $outPath $inPath

echo “copy start”

cp -f $outPath ../cocos2d-x-2.2.6/projects/Casino/Resources/lua/network/
cp -f $outPath ../cocos2d-x-2.2.6/projects/Casino/Resources_ar/lua/network/
cp -f $outPath ../cocos2d-x-2.2.6/projects/Casino/Resources_en/lua/network/
cp -f $outPath ../cocos2d-x-2.2.6/projects/Casino/Resources_ma/lua/network/
cp -f $outPath ../cocos2d-x-2.2.6/projects/Casino/Resources_tw/lua/network/

echo “copy complete”