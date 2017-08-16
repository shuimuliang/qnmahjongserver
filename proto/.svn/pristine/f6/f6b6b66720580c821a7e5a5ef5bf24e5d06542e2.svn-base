#!/bin/sh

outPath=Proto_t.lua
inPath=pf.proto

lua runMe_.lua $outPath $inPath

cp -f $outPath ../cocos2d-x-2.2.6/projects/Casino/Resources/lua/network/

echo “copy complete”

echo “delete —— ${outPath}”
rm $outPath
