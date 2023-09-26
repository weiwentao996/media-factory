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
	content.AddConfigPath("./")     //设置读取的文件路径
	content.SetConfigName("advice") //设置读取的文件名
	content.SetConfigType("yaml")   //设置文件的类型
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
	GenAdviceVideoWithSetting(advice.Page, "zh-CN-YunyangNeural", "../output", &common.AdviceFoSetting{
		FpsFix:  0.3,
		FpsRate: 6,
	}, &common.AdviceFoStyle{
		Align:      "center",
		Size:       80,
		Background: "https://img.iuhub.cn/unsplash/nature/photo-1509316975850-ff9c5deb0cd9.jpg",
		Color: &common.Color{
			R: 255,
			G: 255,
			B: 255,
		},
	}, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiIsImF2YXRhciI6Imh0dHBzOi8vaW1nLml1aHViLmNuL3Vuc3BsYXNoL25hdHVyZS9waG90by0xNjkyMzAzNjEwMDc1LWJhZTU2MzI2MDMxMC5qcGciLCJwYXNzd29yZCI6IjEyMzQ1NiIsImV4cCI6MTY5NTc3MDgzNS4yNjAyNjQsImlzcyI6Ind3dCJ9.ugnDIHroCob7SxqeMpID52bmvAiVvTxCBjqpTIfEFwM")
	fmt.Printf("按任意键结束 ...")
	endKey := make([]byte, 1)
	os.Stdin.Read(endKey)
}
