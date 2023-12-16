package voice

import (
	"encoding/json"
	"fmt"
	"github.com/go-audio/wav"
	"github.com/weiwentao996/media-factory/lib/common"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func GetWavDuration(fileName string) (time.Duration, error) {
	// 打开.wav文件
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 创建一个.wav解码器
	decoder := wav.NewDecoder(file)

	// 读取.wav文件的元数据信息
	_, err = decoder.Seek(0, 0)
	if err != nil {
		return 0, err
	}

	return decoder.Duration()
}

func MergeWAV(inputPattern string, outputFileName string) error {
	// 获取匹配通配符的文件列表
	files, err := filepath.Glob(inputPattern)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return err
	}

	output, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer output.Close()

	for _, file := range files {
		input, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer input.Close()

		_, err = io.Copy(output, input)
		if err != nil {
			return err
		}
	}
	return nil
}

// 计算匹配通配符的文件数量
func countFiles(pattern string) int {
	files, _ := filepath.Glob(pattern)
	return len(files)
}

func GenEdgeVoice(content []string, outPath string) ([]common.VttContent, error) {
	// video -r 0.1  -f image2 -i ./sources/img/%d.jpg  -s 640x480 ./sources/video/output.mp4
	cmd := exec.Command("edge-tts", "--voice", "zh-CN-XiaoyiNeural", "--text", strings.Join(content, "。"), "--write-media", outPath, "--write-subtitles", outPath+".vtt")
	fmt.Println(cmd.String())
	// 获取输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	go copyOutput(stdout)
	go copyOutput(stderr)

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	vtt, err := ReadVtt(outPath + ".vtt")
	if err != nil {
		return nil, err
	}

	return vtt, nil
}

func GenEdgeVoiceOnly(content []string, outPath string) ([]common.VttContent, error) {
	// video -r 0.1  -f image2 -i ./sources/img/%d.jpg  -s 640x480 ./sources/video/output.mp4
	cmd := exec.Command("edge-tts", "--voice", "zh-CN-XiaoyiNeural", "--text", strings.Join(content, "。"), "--write-media", outPath, "--write-subtitles", outPath+".vtt")
	fmt.Println(cmd.String())
	// 获取输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	go copyOutput(stdout)
	go copyOutput(stderr)

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	vtt, err := ReadVtt(outPath + ".vtt")
	if err != nil {
		return nil, err
	}

	return vtt, nil
}

func ReadVtt(filePath string) ([]common.VttContent, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	content := string(fileContent)
	// 解析字幕文件
	subtitles := parseWebVTT(content)
	return subtitles, nil
}

func parseWebVTT(vtt string) []common.VttContent {
	var subtitles []common.VttContent
	lines := strings.Split(vtt, "\n")

	// 使用正则表达式匹配时间戳和字幕内容
	re := regexp.MustCompile(`(\d+:\d+:\d+\.\d+) --> (\d+:\d+:\d+\.\d+)`)
	var currentSubtitle common.VttContent

	for _, line := range lines {
		// 忽略或跳过 "WEBVTT" 行
		if strings.TrimSpace(line) == "WEBVTT" {
			continue
		}
		if re.MatchString(line) {
			if currentSubtitle.Content != "" {
				// 如果当前字幕块不为空，则添加到字幕切片中
				subtitles = append(subtitles, currentSubtitle)
			}
			matches := re.FindAllStringSubmatch(strings.TrimSpace(line), -1)
			if len(matches) > 0 {
				// 提取开始时间和结束时间
				startTime := parseTime(matches[0][1])
				endTime := parseTime(matches[0][2])
				currentSubtitle = common.VttContent{
					Time: [2]float64{startTime, endTime},
				}
			}
		} else if len(strings.TrimSpace(line)) > 0 {
			// 非空行为字幕内容
			currentSubtitle.Content += strings.TrimSpace(line) + "\n"
		}
	}

	// 处理最后一个字幕块
	if currentSubtitle.Content != "" {
		subtitles = append(subtitles, currentSubtitle)
	}

	return subtitles
}

func parseTime(timestamp string) float64 {
	parts := strings.Split(timestamp, ":")
	hours := parseFloat(parts[0])
	minutes := parseFloat(parts[1])
	seconds := parseFloat(parts[2])
	return hours*3600 + minutes*60 + seconds
}

func parseFloat(s string) float64 {
	f := 0.0
	fmt.Sscanf(s, "%f", &f)
	return f
}

func copyOutput(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading from pipe:", err)
			break
		}
		if n == 0 {
			break
		}
		fmt.Print(string(buf[:n]))
	}
}

type EdgeTtsRsp struct {
	Voice string              `json:"voice"`
	Vtt   []common.VttContent `json:"vtt"`
}

func GenEdgeVoiceOnline(content []string, voiceType, outPath string, token *string) []common.VttContent {
	// 要发送的数据
	requestData := map[string]interface{}{
		"voice":    voiceType,
		"content":  content,
		"out_path": "./output.wav",
	}

	// 目标 URL
	url := "https://tbg.iuhub.cn/voice/edge"

	ms, err := json.Marshal(requestData)
	if err != nil {
		panic(err)
	}
	// 执行 POST 请求
	response, err := httpPost(url, ms, token)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// 处理响应
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("HTTP 响应状态码:", response.Status)
	rsp := EdgeTtsRsp{}

	json.Unmarshal(responseBody, &rsp)
	base64ToWAV(rsp.Voice, outPath)
	return rsp.Vtt
}

func GenAzureVoiceOnline(content string, voiceType, audioPath string, token string) {
	// 要发送的数据
	requestData := map[string]interface{}{
		"voice":   voiceType,
		"content": content,
	}

	// 目标 URL
	url := "https://tbg.iuhub.cn/voice/azure"

	ms, err := json.Marshal(requestData)
	if err != nil {
		panic(err)
	}
	// 执行 POST 请求
	response, err := httpPost(url, ms, &token)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// 处理响应
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("HTTP 响应状态码:", response.Status)

	base64ToWAV(string(responseBody), audioPath)
}
