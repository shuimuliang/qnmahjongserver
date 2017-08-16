require"parsePB"

local function superFindProtocol(str)
	local t = {}      -- table to store the indices
	local m = 0
	local i = 0
	while true do
		i = string.find(str, "message", i+1)   -- find 'next' newline
		if i == nil then --[[print(string.sub(str,m))--]] table.insert(t, string.sub(str,m)) break end
		local num = string.sub(str,m,i-1)
		if m ~= 0 then table.insert(t, num) --[[print(num)--]] end
		m = i
	end
	return t
end

function analytical(path)
	local s = ""
	local T = {}
	local f = assert(io.open(path, "rb"))
	local linenum = 0
	
	while true do
		local buffer = f:read("*l")
		linenum = linenum + 1
		if not buffer then
			break
		end
		local start1, end1 = string.find(buffer, "//")
		if start1 then
			local num = string.sub(buffer,0,start1-1)
			s = s .. num
		else
			s = s .. buffer
		end
		
	end
	local T_newStre = superFindProtocol(s)
	return T_newStre
end

local outPath = arg[1] or "Proto_pf.lua"
local inPath = arg[2] or "pf.proto"

local function Newmain()	
	local t = analytical(inPath)
	local s = [[
local protoPack = "pf"

]]
	
	for k,v in pairs(t) do
		local res = getPF(v)
		if not res then
			print("")
		end
		s = s..getPF(v)
	end
	local file = io.open(outPath,"w")
	file:write(s)                                             --写入信息
	file:close()                                              --关闭文件
end

Newmain()