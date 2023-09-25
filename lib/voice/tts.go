package voice

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type ttsRsp struct {
	Success bool `json:"success"`
	Code    int  `json:"code"`
	Message struct {
		Description string `json:"description"`
	} `json:"message"`
	Result struct {
		Lang       string  `json:"lang"`
		SpkId      int     `json:"spk_id"`
		Speed      float64 `json:"speed"`
		Volume     float64 `json:"volume"`
		SampleRate int     `json:"sample_rate"`
		Duration   float64 `json:"duration"`
		SavePath   string  `json:"save_path"`
		Audio      string  `json:"audio"`
	} `json:"result"`
}

func CalVoiceTime(content []string, output string) float64 {
	// 要发送的数据
	requestData := []byte(fmt.Sprintf(`{
	   "text": "%s",
	   "spk_id": 0,
	   "speed": 1.0,
	   "volume": 1.0,
	   "sample_rate": 0,
	   "save_path": "/mnt/TTS.wav"
	}`, strings.Join(content, "。")))

	// 目标 URL
	url := "http://1.15.92.254:8887/paddlespeech/tts"

	// 执行 POST 请求
	response, err := httpPost(url, requestData, nil)
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
	rsp := ttsRsp{}

	json.Unmarshal(responseBody, &rsp)
	return rsp.Result.Duration
}

func GetVoiceTTS(content []string, output string) float64 {
	// 要发送的数据
	requestData := []byte(fmt.Sprintf(`{
	   "text": "%s",
	   "spk_id": 0,
	   "speed": 1.0,
	   "volume": 1.0,
	   "sample_rate": 0,
	   "save_path": "/mnt/TTS.wav"
	}`, strings.Join(content, ",")))

	// 目标 URL
	url := "http://1.15.92.254:8887/paddlespeech/tts"

	// 执行 POST 请求
	response, err := httpPost(url, requestData, nil)
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
	rsp := ttsRsp{}

	json.Unmarshal(responseBody, &rsp)

	base64ToWAV(rsp.Result.Audio, output)

	return rsp.Result.Duration
}

// httpPost 执行 HTTP POST 请求
func httpPost(url string, data []byte, token *string) (*http.Response, error) {
	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+*token)

	// 执行请求
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func base64ToWAV(base64Data string, output string) {
	// 解码Base64数据
	audioData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		panic(err)
	}

	// 创建并写入WAV文件
	err = writeWAVFile(output, audioData)
	if err != nil {
		panic(err)
	}

	fmt.Println("成功转换为WAV文件：", output)
}

func writeWAVFile(fileName string, audioData []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入音频数据
	file.Write(audioData)
	return nil
}
