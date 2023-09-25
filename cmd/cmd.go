package cmd

import (
	"fmt"
	"github.com/weiwentao996/media-factory/lib/common"
	"github.com/weiwentao996/media-factory/lib/img"
	"github.com/weiwentao996/media-factory/lib/video"
	"github.com/weiwentao996/media-factory/lib/voice"
	"os"
	"strings"
	"time"
)

// GenPPTVideoWithSetting 配置生成Video
func GenPPTVideoWithSetting(essay []common.PageData, outPath string, setting *common.PPTSetting) {
	fmt.Printf("\033[1;32;42m%s\n", "开始生成视频!")
	path := fmt.Sprintf("%s/%d", outPath, time.Now().Unix())
	if err := os.MkdirAll(path, 0444); err != nil {
		panic(err)
	}

	counter := 0
	allFpsCount := 0
	var maxTime float64
	for i, e := range essay {
		bgmPath := fmt.Sprintf("%s/%d.wav", path, i)
		vtt, err := voice.GenEdgeVoice(e.Content, bgmPath)
		if err != nil {
			panic(err)
		}
		endTime := vtt[len(vtt)-1].Time[1]

		essay[i].Style.LiveTime = endTime
		maxTime += endTime
		conf := common.GetConfig(setting, essay[i])
		allFpsCount += conf.FpsCount
	}
	bgmPath := fmt.Sprintf("%s/voice.wav", path)
	var voices []string
	for i, e := range essay {
		fmt.Printf("\033[1;32;42m%s%d%s\n", "正在生成第 ", i+1, " 幕视频帧......")
		conf := common.GetConfig(setting, e)
		voices = append(voices, e.Content...)
		img.GenPPTImage(path, e, counter, allFpsCount, setting)
		counter += conf.FpsCount
	}
	fmt.Printf("\033[1;32;42m%s\n", "正在生成音频......")
	err := voice.MergeWAV(path+"/*.wav", bgmPath)
	if err != nil {
		panic(err)
	}

	fpsRate := 8.0
	if setting.FpsRate != 0 {
		fpsRate = setting.FpsRate
	}

	if setting.MaxTime != 0 {
		maxTime = setting.MaxTime
	}

	fmt.Printf("\033[1;32;42m%s\n", "正在合成视频......")

	if setting != nil && setting.MusicRule != "" {
		bgmPath = setting.MusicRule
	}

	if err := video.MultiImageToVideo(path+"/%05d.png", bgmPath, path, fpsRate, (maxTime+1)*1000); err != nil {
		panic(err)
	}

	fmt.Printf("\033[1;32;42m%s\n", "已生成视频!")
}

// GenAdviceVideoWithSetting 配置生成Video
func GenAdviceVideoWithSetting(advice []common.VttContent, outPath string, setting *common.AdviceFoSetting, style *common.AdviceFoStyle, token string) {
	fmt.Printf("\033[1;32;42m%s\n", "开始生成视频!")
	path := fmt.Sprintf("%s/%d", outPath, time.Now().Unix())
	if err := os.MkdirAll(path, 0444); err != nil {
		panic(err)
	}

	var advicePage []common.VttContent
	var preTime float64
	for i, e := range advice {
		bgmPath := fmt.Sprintf("%s/%d.wav", path, i)
		vtt := voice.GenEdgeVoiceOnline([]string{e.Content}, bgmPath, &token)
		for j := 0; j < len(vtt); j++ {
			vtt[j].ContentImage = e.ContentImage
			vtt[j].Content = strings.Replace(vtt[j].Content, " ", "", -1)
			vtt[j].Time[0] = vtt[j].Time[0] + preTime
			vtt[j].Time[1] = vtt[j].Time[1] + preTime
		}

		preTime = vtt[len(vtt)-1].Time[1]
		advicePage = append(advicePage, vtt...)
	}

	counter := 0
	bgmPath := fmt.Sprintf("%s/voice.wav", path)

	for i, e := range advicePage {
		fmt.Printf("\033[1;32;42m%s%d%s\n", "正在生成第 ", i+1, " 幕视频帧......")
		counter = img.GenAdviceImage(path, &e, advicePage[len(advicePage)-1].Time[1]+1, counter, setting, style)
	}

	fmt.Printf("\033[1;32;42m%s\n", "正在生成音频......")
	err := voice.MergeWAV(path+"/*.wav", bgmPath)
	if err != nil {
		panic(err)
	}

	fpsRate := 8.0
	if setting.FpsRate != 0 {
		fpsRate = setting.FpsRate
	}

	fmt.Printf("\033[1;32;42m%s\n", "正在合成视频......")

	if err := video.MultiImageToVideo(path+"/%05d.png", bgmPath, path, fpsRate, (advicePage[len(advicePage)-1].Time[1]+1)*1000); err != nil {
		panic(err)
	}

	fmt.Printf("\033[1;32;42m%s\n", "已生成视频!")
}
