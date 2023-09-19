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

// GenVideoWithSetting 配置生成Video
func GenVideoWithSetting(essay []common.PageData, outPath string, setting *common.Setting) {
	fmt.Printf("\033[1;32;42m%s\n", "开始生成视频!")
	path := fmt.Sprintf("%s/%d", outPath, time.Now().Unix())
	if err := os.MkdirAll(path, 0444); err != nil {
		panic(err)
	}

	counter := 0
	allFpsCount := 0

	for _, e := range essay {
		conf := common.GetConfig(setting, e)
		allFpsCount += conf.FpsCount
	}

	for i, e := range essay {
		conf := common.GetConfig(setting, e)
		fmt.Printf("\033[1;32;42m%s%d%s\n", "正在生成第 ", i+1, " 幕视频帧......")
		img.GenImage(path, e, counter, allFpsCount, setting)
		counter += conf.FpsCount
	}

	fpsRate := 8.0
	if setting.FpsRate != 0 {
		fpsRate = setting.FpsRate
	}

	maxTime := float64(counter)/fpsRate + 1
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
