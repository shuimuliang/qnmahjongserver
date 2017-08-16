package tdh

type Fanshu struct {
	// 加分牌型
	Gangshanghuajiabei bool // 杠上花加倍
}

func (fanshu Fanshu) GetHuTypes() []int32 {
	var huTypes []int32
	if fanshu.Gangshanghuajiabei {
		huTypes = append(huTypes, 206)
	}
	if len(huTypes) == 0 {
		huTypes = append(huTypes, 204)
	}
	return huTypes
}
