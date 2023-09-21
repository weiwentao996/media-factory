package cmd

import (
	"fmt"
	"github.com/weiwentao996/media-factory/lib/common"
	"github.com/weiwentao996/media-factory/lib/img"
	"github.com/weiwentao996/media-factory/lib/video"
	"github.com/weiwentao996/media-factory/lib/voice"
	"math"
	"os"
	"time"
)

// GenVideoWithSetting 配置生成Video
func GenVideoWithSetting(essay []common.PageData, outPath string, setting *common.Setting) {
	fmt.Printf("\033[1;32;42m%s\n", "开始生成视频!")
	path := fmt.Sprintf("%s/%d", outPath, time.Now().Unix())
	bgmPath := path + "voice.wav"
	if err := os.MkdirAll(path, 0444); err != nil {
		panic(err)
	}

	counter := 0
	allFpsCount := 0
	var voiceTime []int
	for i, e := range essay {
		t := voice.CalVoiceTime(e.Content, bgmPath)
		second := int(math.Ceil(t)) - 1
		essay[i].Style.LiveTime = second
		voiceTime = append(voiceTime, second)
		conf := common.GetConfig(setting, essay[i])
		allFpsCount += conf.FpsCount
	}

	var voices []string
	for i, e := range essay {
		fmt.Printf("\033[1;32;42m%s%d%s\n", "正在生成第 ", i+1, " 幕视频帧......")
		conf := common.GetConfig(setting, e)
		voices = append(voices, e.Content...)
		img.GenImage(path, e, counter, allFpsCount, setting)
		counter += conf.FpsCount
	}
	fmt.Printf("\033[1;32;42m%s\n", "正在生成音频......")
	voice.GetVoiceTTS(voices, bgmPath)
	fpsRate := 8.0
	if setting.FpsRate != 0 {
		fpsRate = setting.FpsRate
	}

	maxTime := float64(counter)/fpsRate + 1
	if setting.MaxTime != 0 {
		maxTime = setting.MaxTime
	}

	fmt.Printf("\033[1;32;42m%s\n", "正在合成视频......")

	if setting != nil && setting.MusicRule != "" {
		bgmPath = setting.MusicRule
	}

	if err := video.MultiImageToVideo(path+"/%05d.png", bgmPath, path, fpsRate, maxTime); err != nil {
		panic(err)
	}

	fmt.Printf("\033[1;32;42m%s\n", "已生成视频!")
}
