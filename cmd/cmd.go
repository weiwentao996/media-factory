package cmd

import (
	"fmt"
	"github.com/weiwentao996/media-factory/lib/common"
	"github.com/weiwentao996/media-factory/lib/img"
	"github.com/weiwentao996/media-factory/lib/video"
	"github.com/weiwentao996/media-factory/sources"
	"os"
	"time"
)

// GenVideo 生成默认Video
func GenVideo(essay []img.ImageData, outPath string) {
	fmt.Printf("\033[1;32;42m%s\n", "开始生成视频!")
	path := fmt.Sprintf("%s/%d", outPath, time.Now().Unix())
	if err := os.MkdirAll(path, 0444); err != nil {
		panic(err)
	}

	for i, e := range essay {
		fmt.Printf("\033[1;32;42m%s%d%s\n", "正在生成第 ", i+1, " 幕视频帧......")
		img.GenImage(path, e, i, len(essay), nil)
	}

	fpsRate := 8.0

	maxTime := float64(img.FpsCount*len(essay))/fpsRate + 1
	fmt.Printf("\033[1;32;42m%s\n", "正在合成视频......")
	if err := video.MultiImageToVideo(path+"/%05d.png", sources.Path+"/mp3/Winter.mp3", path, fpsRate, maxTime); err != nil {
		panic(err)
	}
	fmt.Printf("\033[1;32;42m%s\n", "已生成视频!")
}

// GenVideoWithSetting 配置生成Video
func GenVideoWithSetting(essay []img.ImageData, outPath string, setting *common.Setting) {
	fmt.Printf("\033[1;32;42m%s\n", "开始生成视频!")
	path := fmt.Sprintf("%s/%d", outPath, time.Now().Unix())
	if err := os.MkdirAll(path, 0444); err != nil {
		panic(err)
	}

	for i, e := range essay {
		fmt.Printf("\033[1;32;42m%s%d%s\n", "正在生成第 ", i+1, " 幕视频帧......")
		img.GenImage(path, e, i, len(essay), setting)
	}

	fpsRate := 8.0
	if setting.FpsRate != 0 {
		fpsRate = setting.FpsRate
	}

	maxTime := float64(img.FpsCount*len(essay))/fpsRate + 1
	if setting.MaxTime != 0 {
		maxTime = setting.MaxTime
	}

	fmt.Printf("\033[1;32;42m%s\n", "正在合成视频......")
	bgmPath := sources.Path + "/mp3/Winter.mp3"

	if setting != nil && setting.MusicRule != "" {
		bgmPath = setting.MusicRule
	}

	if err := video.MultiImageToVideo(path+"/%05d.png", bgmPath, path, fpsRate, maxTime); err != nil {
		panic(err)
	}

	fmt.Printf("\033[1;32;42m%s\n", "已生成视频!")
}
