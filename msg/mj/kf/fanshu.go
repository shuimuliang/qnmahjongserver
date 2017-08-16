package kf

type Fanshu struct {
	// 加分牌型
	Zimofanbei         bool // 自摸翻倍
	Gangshanghuajiabei bool // 杠上花加倍
	Qiduifanbei        bool // 七对翻倍
}

func (fanshu Fanshu) GetHuTypes() []int32 {
	var huTypes []int32
	if fanshu.Qiduifanbei {
		huTypes = append(huTypes, 205)
	}
	if fanshu.Gangshanghuajiabei {
		huTypes = append(huTypes, 206)
	}
	return huTypes
}
