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
func GenPPTVideoWithSetting(essay []common.PageData, outPath string, setting *common.PPTSetting) error {
	fmt.Printf("\033[1;32;42m%s\n", "开始生成视频!")
	path := fmt.Sprintf("%s/%d", outPath, time.Now().Unix())
	if err := os.MkdirAll(path, 0444); err != nil {
		return err
	}

	counter := 0
	allFpsCount := 0
	var maxTime float64
	for i, e := range essay {
		bgmPath := fmt.Sprintf("%s/%d.wav", path, i)
		vtt, err := voice.GenEdgeVoice(e.Content, bgmPath)
		if err != nil {
			return err
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
		return err
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
		return err
	}

	fmt.Printf("\033[1;32;42m%s\n", "已生成视频!")
	return nil
}

// GenAdviceVideoWithSetting 配置生成Video
func GenAdviceVideoWithSetting(advice []common.VttContent, voiceType, outPath string, setting *common.AdviceFoSetting, style *common.AdviceFoStyle, token string) error {
	if setting.FpsRate == 0 {
		setting.FpsRate = 6.0
	}
	fmt.Printf("\033[1;32;42m%s\n", "开始生成视频!")
	path := fmt.Sprintf("%s/%d", outPath, time.Now().Unix())
	if err := os.MkdirAll(path, 0444); err != nil {
		return err
	}
	bgmPath := fmt.Sprintf("%s/voice.wav", path)
	var preTime float64
	for i, content := range advice {
		cllVttList, err := voice.GenEdgeVoiceOnline([]string{content.Content}, voiceType, fmt.Sprintf("%s/%05d.wav", path, i), &token)
		if err != nil {
			return err
		}

		advice[i].Time[0] = preTime
		preTime += cllVttList[len(cllVttList)-1].Time[1] + setting.FpsFix
		advice[i].Time[1] = preTime

	}

	err := voice.MergeWAV(fmt.Sprintf("%s/*.wav", path), bgmPath)
	if err != nil {
		return err
	}

	counter := 0

	advice = append(advice, common.VttContent{
		Content:      "Ending...",
		Avatar:       "",
		Nickname:     "",
		ContentImage: "",
		Time:         [2]float64{advice[len(advice)-1].Time[1], advice[len(advice)-1].Time[1] + 3},
		CommentTime:  time.Time{},
	})
	for i, content := range advice {
		fmt.Printf("\033[1;32;42m%s%d%s\n", "正在生成第 ", i+1, " 幕视频帧......")
		content.Content = strings.Replace(content.Content, " ", "", -1)
		switch setting.Model {
		case "music":
			counter, err = img.GenMusicImage(path, &content, advice[len(advice)-1].Time[1]+1, counter, setting, style)
			if err != nil {
				return err
			}
		default:
			counter, err = img.GenAdviceImage(path, &content, advice[len(advice)-1].Time[1]+1, counter, setting, style)
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("\033[1;32;42m%s\n", "正在合成视频......")
	videoFilePath := path + "/video.mp4"
	if err := video.MultiImageToVideo(path+"/%05d.png", bgmPath, videoFilePath, setting.FpsRate, (advice[len(advice)-1].Time[1]+1)*1000); err != nil {
		return err
	}

	if setting.BgmUrl != "" {
		finishFilePath := path + "/finish.mp4"
		fmt.Printf("\033[1;32;42m%s\n", "正在添加bgm......")
		video.AddBgm(videoFilePath, setting.BgmUrl, finishFilePath, 0.1)
	}

	fmt.Printf("\033[1;32;42m%s\n", "已生成视频!")
	return nil
}

func GenVideoFast(advice []common.VttContent, voiceType, output, bmg string, style *common.AdviceFoStyle, bgmVolume float32, token string) error {
	var audioPath, imgPath, videoPath string
	output = fmt.Sprintf("%s/%d", output, time.Now().Unix())
	for i, content := range advice {
		fmt.Printf("\033[1;32;42m正在生成第%d页......\n", i+1)
		if err := os.MkdirAll(output, 0444); err != nil {
			return err
		}

		audioPath = fmt.Sprintf("%s/%d.wav", output, i)
		imgPath = fmt.Sprintf("%s/%d.png", output, i)
		videoPath = fmt.Sprintf("%s/%d.mp4", output, i)

		// 生成音频
		voice.GenAzureVoiceOnline(content.Content, voiceType, audioPath, token)

		// 生成图片
		img.GenMusicImageFast(imgPath, &content, style)

		// 合成视频
		err := video.ImageAndVoice2Video(imgPath, audioPath, videoPath)
		if err != nil {
			return err
		}
	}

	// 融合视频
	fmt.Printf("\033[1;32;42m%s\n", "融合视频......")
	mergeVideoPath, err := video.Merge(fmt.Sprintf("%s/*.mp4", output), output)
	if err != nil {
		return err
	}

	// 添加BGM
	finishFilePath := output + "/finish.mp4"
	fmt.Printf("\033[1;32;42m%s\n", "正在添加bgm......")
	err = video.AddBgm(mergeVideoPath, bmg, finishFilePath, bgmVolume)
	if err != nil {
		return err
	}

	err = os.RemoveAll(mergeVideoPath)
	if err != nil {
		return err
	}

	return nil
}
