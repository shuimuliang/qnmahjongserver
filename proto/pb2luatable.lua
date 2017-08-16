
local function findParms(s,tag,parms)
	local i,j,_type,_name = 0,0,""
	while true do
		i,j,_,_type,_,_name,_ = string.find(s,tag .. "(%s+)(%w+)(%s+)(%w+)(%s+)=",j+1)
		if i == nil then break end
		table.insert(parms, {_type = _type,_name = _name, _tag = tag})
	end
	return parms
end

local function dealSignleParm(parms)
	local str = ""
	for k,tt in pairs(parms) do
		str = str .. [[
		{"]] .. tt._tag .. [[", "]] .. tt._type ..[[", "]] .. tt._name .. [["},
]]
		-- if tt._type == "string"
		-- or tt._type == "double"
		-- or tt._type == "int32"
		-- or tt._type == "int64"
		-- or tt._type == "float"
		-- or tt._type == "bool" then
		--
		-- 	str = str .. [[			self:printItem("]] .. tt._name .. "\")\n"
		-- else
		-- 	str = str .."\n".. [[		local ]]
		-- 		.. tt._name .. [[= LCPB_]]
		-- 		.. tt._type ..".new()\n		"
		-- 		.. tt._name .."._msg = self._msg."
		-- 		.. tt._name .. "\n		"
		-- 		.. tt._name .. ":printTitle()\n		"
		-- 		.. tt._name .. ":print()\n\n"
		-- end
	end

	return str
end

local function dealMulitParm(parms)
	local str = ""

	for k,tt in pairs(parms) do
		str = str .. [[
		{"]] .. tt._tag .. [[", "]] .. tt._type ..[[", "]] .. tt._name .. [["},
]]
		-- str = str .."\n"..  [[			GDebug("]] .. tt._name .. [[ len %s",#self._msg.]] .. tt._name .. [[)]] .. "\n"
		--
		-- if tt._type == "string"
		-- or tt._type == "double"
		-- or tt._type == "int32"
		-- or tt._type == "int64"
		-- or tt._type == "bool" then
		--
		-- 	str = str .. [[			GDebug("]]..tt._name..[[ %s",table.concat(self._msg.]]..tt._name..[[," | "))]].. "\n\n"
		-- else
		--
		-- 	str = str .. [[			for i=1,#self._msg.]]..tt._name..[[ do]] .. "\n"
		-- 			.. [[				local data = LCPB_]].. tt._type .. [[.new()]] .. "\n"
		-- 			.. [[				data._msg = self._msg.]]..tt._name.."[i]" .. "\n"
		-- 			.. [[				data:printTitle()]] .. "\n"
		-- 			.. [[				data:print()]] .. "\n"
		-- 			.. [[			end]] .. "\n\n"
		-- end
	end

	return str
end

function getPF(s)
	local _,_,pfName = string.find(s,"message (%w+)")
	if not pfName then
		print("can not find pfName %s",s)
		return nil
	else
		print("deal %s",pfName)
	end
	local parms = {}

	-- find sigle parm
	local signleparms = {}
	local mulitparms = {}
	signleparms = findParms(s,"required",signleparms)
	signleparms = findParms(s,"optional",signleparms)

	local mulitparms = {}
	mulitparms = findParms(s,"repeated",mulitparms)


	local resPF = [[	]] .. pfName ..[[ = {
]] .. dealSignleParm(signleparms) .. dealMulitParm(mulitparms) ..[[
	},
]]
	--print(resPF)
	return resPF

end
