package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/weiwentao996/media-factory/lib/common"
	"os"
	"testing"
)

type Essay struct {
	Page []common.PageData `mapstructure:"page" `
}

func TestMultiImageToVideo(t *testing.T) {
	essay := Essay{}
	content := viper.New()
	content.AddConfigPath("./")      //设置读取的文件路径
	content.SetConfigName("content") //设置读取的文件名
	content.SetConfigType("yaml")    //设置文件的类型
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("\033[1;31;42m%v\n", err)
			fmt.Printf("\033[1;31;42m%s\n", "生成视频失败！")
			fmt.Printf("按任意键结束 ...")
			endKey := make([]byte, 1)
			os.Stdin.Read(endKey)
		}

	}()
	if err := content.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := content.Unmarshal(&essay); err != nil {
		panic(err)
	}

	fmt.Printf("\033[1;32;42m%s\n", "读取文件成功!")
	//cmd.GenVideo(essay.Page, "./sources/output")
	GenPPTVideoWithSetting(essay.Page, "../output", &common.PPTSetting{
		FpsRate:         6,
		HighPerformance: true,
	})
	fmt.Printf("按任意键结束 ...")
	endKey := make([]byte, 1)
	os.Stdin.Read(endKey)
}

type VttContentItem struct {
	Content      string `json:"content" mapstructure:"content"`
	ContentImage string `json:"content_image" mapstructure:"content_image"`
}

type Advice struct {
	Page []common.VttContent `mapstructure:"page" `
}

func TestGenAdviceVideoWithSetting(t *testing.T) {
	advice := Advice{}
	content := viper.New()
	content.AddConfigPath("./")    //设置读取的文件路径
	content.SetConfigName("music") //设置读取的文件名
	content.SetConfigType("yaml")  //设置文件的类型
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("\033[1;31;42m%v\n", err)
			fmt.Printf("\033[1;31;42m%s\n", "生成视频失败！")
			fmt.Printf("按任意键结束 ...")
			endKey := make([]byte, 1)
			os.Stdin.Read(endKey)
		}

	}()
	if err := content.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := content.Unmarshal(&advice); err != nil {
		panic(err)
	}

	fmt.Printf("\033[1;32;42m%s\n", "读取文件成功!")
	//cmd.GenVideo(essay.Page, "./sources/output")
	GenAdviceVideoWithSetting(advice.Page, "zh-CN-YunxiNeural", "../output", &common.AdviceFoSetting{
		FpsFix:  0.3,
		FpsRate: 6,
		Model:   "music",
		BgmUrl:  "http://m7.music.126.net/20231011100310/190ccb3526417d0a8097b33427fe87de/ymusic/030b/545b/5308/4a2c7e7115c11526e4c4db18c347c03c.mp3",
	}, &common.AdviceFoStyle{
		Align:      "center",
		Size:       48,
		Background: "https://img.iuhub.cn/unsplash/wallpapers/photo-1451224222030-cee2f5dbcd10.jpg",
		Color: &common.Color{
			R: 0,
			G: 0,
			B: 0,
		},
	}, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiIsImF2YXRhciI6Imh0dHBzOi8vaW1nLml1aHViLmNuL3Vuc3BsYXNoL25hdHVyZS9waG90by0xNjkyMzAzNjEwMDc1LWJhZTU2MzI2MDMxMC5qcGciLCJwYXNzd29yZCI6IjEyMzQ1NiIsImV4cCI6MTY5NzAyMjMyOS43NDg5NTQsImlzcyI6Ind3dCJ9.llbAofN4i5IP_JiZEuxhx-pzXcaI5YAOttw3cApLlSw")
	fmt.Printf("按任意键结束 ...")
	endKey := make([]byte, 1)
	os.Stdin.Read(endKey)
}
