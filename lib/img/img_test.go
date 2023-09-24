package img

import (
	"github.com/weiwentao996/media-factory/lib/common"
	"testing"
)

func TestGenBanFoImage(t *testing.T) {
	GenAdviceImage("./output", &common.VttContent{
		Content:      "你好，这里是主播播报欢迎收听",
		ContentImage: "https://img.iuhub.cn/doutu/%E4%BD%A0%E4%B8%AA%E5%B0%8F%E8%A5%BF%E7%93%9C/20151207451199_IvDkZY.jpg",
		Time:         [2]float64{0.1, 6.2},
	}, 16, 0, &common.AdviceFoSetting{
		FpsRate: 6,
	}, &common.AdviceFoStyle{
		Align: "center",
		Size:  80,
		Color: &common.Color{
			R: 255,
			G: 255,
			B: 255,
		},
		Background: "https://img.iuhub.cn/unsplash/nature/photo-1509316975850-ff9c5deb0cd9.jpg",
	})
}
