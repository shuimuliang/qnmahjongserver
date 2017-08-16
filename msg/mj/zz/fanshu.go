package zz

type Fanshu struct {
	// 加分牌型
	Qiduijiabei        bool // 七对加倍
	Gangshanghuajiabei bool // 杠上花加倍
	Fourhun            bool // 四个混牌
}

func (fanshu Fanshu) GetHuTypes() []int32 {
	var huTypes []int32
	if fanshu.Qiduijiabei {
		huTypes = append(huTypes, 205)
	}
	if fanshu.Gangshanghuajiabei {
		huTypes = append(huTypes, 206)
	}
	if fanshu.Fourhun {
		huTypes = append(huTypes, 212)
	}
	if len(huTypes) == 0 {
		huTypes = append(huTypes, 204)
	}
	return huTypes
}
