package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/weiwentao996/media-factory/cmd"
	"github.com/weiwentao996/media-factory/lib/common"
	"os"
)

type Essay struct {
	Page []common.PageData `mapstructure:"page" `
}

func main() {
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
	cmd.GenVideoWithSetting(essay.Page, "./output", &common.Setting{
		FpsRate:         6,
		HighPerformance: true,
	})
	fmt.Printf("按任意键结束 ...")
	endKey := make([]byte, 1)
	os.Stdin.Read(endKey)
}
