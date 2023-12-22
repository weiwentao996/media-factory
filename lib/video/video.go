package video

import (
	"bytes"
	"fmt"
	"github.com/weiwentao996/media-factory/lib/img"
	"github.com/weiwentao996/media-factory/sources"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

//多张图片组转视频 video -f image2 -i ~/Desktop/images/image%d.jpg  -vcodec libx264  ~/Desktop/test.mp4 -y
//图片按固定的帧率转视频 video -f image2 -i ~/Desktop/images/image%d.jpg  -vcodec libx264  -r 10 ~/Desktop/test.mp4 -y
//设置图片切换帧率为每秒一张图片 video -r 1 -i ~/Desktop/%d.jpg -vf fps=1 -vcodec libx264 ~/Desktop/test.mp4 -y

// video -threads2 -y -r 10 -i /tmpdir/image%04d.jpg -i audio.mp3 -absf aac_adtstoasc output.mp4
//参数的解释含义：
//-threads 2 以两个线程进行运行， 加快处理的速度。
//-y 对输出文件进行覆盖
//-r 10 fps设置为10帧/秒（不同位置有不同含义，后面再解释）
//-i /tmpdir/image%04d.jpg 输入图片文件，图片文件保存为 image0001.jpg image0002.jpg ….
//-i audio.mp3 输入的音频文件
//-absf aac_adtstoasc 将结果的音频格式转为faac格式时需要这个选项。将音频格式转为faac是因为在iphone上某些音频格式的视频无法播放，例如mp3. 但faac格式的音频的视频在iphone上可以播放。-absf 的意思是设置一个bitstream filter进行某些转换。可以用ffmpeg -bsfs 查看所有支持的bitstream filter。 bitstream filter和 aac_adtstoasc的具体含义我也说不上。但是如果不用这个选项又会导致转换失败。

//不带音频
//video -loop 1 -f image2 -i /tmpdir/image%04d.jpg -vcodec libx264 -r 10 -t 10 test.mp4
//-loop 1循环读输入 0读完就不读了
//-vcode 编码格式libx264
//-b 指定200k码率
//-t 输出视频总时长：

func SingleImageToVideo(imgPath, outPath string) error {
	cmd := exec.Command("video", "-ss", "0", "-t", "10", "-f", "lavfi", "-i", "color=c=0x000000:s=1326x900:r=30", "-i", imgPath, "-filter_complex", "[1:v]scale=1920:1080[v1];[0:v][v1]overlay=0:0[outv]", "-map", "[outv]", "-c:v", "libx264", outPath, "-y")
	// 获取输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// 开始执行命令
	if err := cmd.Start(); err != nil {
		return err
	}

	// 读取输出
	go func() {
		if _, err := io.Copy(os.Stdout, stdout); err != nil {
			return
		}
	}()

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func MultiImageToVideo(imgPathRule, bgmPath, outPath string, fps float64, maxTime float64) error {
	// video -r 0.1  -f image2 -i ./sources/img/%d.jpg  -s 640x480 ./sources/video/output.mp4
	cmd := exec.Command("ffmpeg", "-r", fmt.Sprintf("%f", fps), "-f", "image2", "-i", imgPathRule, "-i", bgmPath, "-t", fmt.Sprintf("%f", maxTime), "-pix_fmt", "yuv420p", outPath, "-y")
	fmt.Println(cmd.String())
	// 获取输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return err
	}

	go copyOutput(stdout)
	go copyOutput(stderr)

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
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

func AddBgm(videoPath, audioPath, outputPath string, volume float32) error {
	// 获取视频时长
	videoDurationCmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoPath)
	videoDuration, err := videoDurationCmd.Output()
	if err != nil {
		return err
	}

	// 获取音频时长
	audioDurationCmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", audioPath)
	audioDuration, err := audioDurationCmd.Output()
	if err != nil {
		return err
	}
	// ffmpeg命令及参数
	cmdArgs := []string{
		"-i", videoPath,
	}
	if string(videoDuration) < string(audioDuration) {
		cmdArgs = append(cmdArgs, "-i", audioPath, "-filter_complex", "[1:a]volume=0.1[a1];[0:a][a1]amix=inputs=2:duration=first[a]", "-strict", "experimental", "-shortest")
	} else {
		cmdArgs = append(cmdArgs, "-stream_loop", "-1", "-i", audioPath, "-filter_complex", "[1:a]volume="+fmt.Sprintf("%f", volume)+"[a1];[0:a][a1]amix=inputs=2:duration=first[a]", "-strict", "experimental")
	}

	cmdArgs = append(cmdArgs,
		"-map", "0:v",
		"-map", "[a]",
		"-c:v", "copy",
		"-c:a", "aac", "-y", outputPath)

	// 执行ffmpeg命令
	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return err
	}
	return nil
}

// ImageAndVoice2Video 图片声音生成视频
func ImageAndVoice2Video(imgPath, voicePath, outPath string) {
	imagePath := imgPath
	audioPath := voicePath
	outputPath := outPath

	cmd := exec.Command("ffmpeg", "-loop", "1", "-i", imagePath, "-i", audioPath, "-c:v", "libx264", "-c:a", "aac", "-strict", "experimental", "-b:a", "192k", "-shortest", outputPath)

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error executing FFmpeg command: %s", err)
	}

	os.RemoveAll(imgPath)
	os.RemoveAll(voicePath)
}

func Merge(inputPattern string, output string) (string, error) {
	tmpVideoPath := fmt.Sprintf("%s/tmp.mp4", output)
	// 列出文件列表
	// 获取匹配通配符的文件列表
	files, err := filepath.Glob(inputPattern)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", err
	}

	// 创建文件列表文本文件
	filelist, err := os.Create("filelist")
	if err != nil {
		return "", err
	}

	defer filelist.Close()

	// 写入文件列表
	for _, file := range files {
		filelist.WriteString(fmt.Sprintf("file '%s'\n", file))
	}

	// 使用 FFmpeg 合并文件
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", "filelist", "-c", "copy", tmpVideoPath, "-y")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	filelist.Close()
	os.RemoveAll("filelist")

	// 写入文件列表
	for _, file := range files {
		os.RemoveAll(file)
	}

	AddProcess(tmpVideoPath, fmt.Sprintf("%s/output.mp4", output))

	os.RemoveAll(tmpVideoPath)

	return fmt.Sprintf("%s/output.mp4", output), nil
}

func AddProcess(videoInput, outputFile string) {
	//以下是一个示例命令，假设进度条 GIF 的宽度是 200 像素，高度是 50 像素，视频的宽度是 1280 像素，高度是 720 像素，持续时间为 10 秒：
	// 这个命令中的 overlay 过滤器包含了 x 和 y 参数。其中 x 参数控制 GIF 的水平位置，使用了一个表达式来实现从左向右移动的效果。y 参数保持在视频底部。
	// lt(-w+(t)*200,0) 将在视频底部左侧将 GIF 从左到右移动，t 代表时间。
	// H-h 用于将 GIF 放置在视频底部。
	file, _ := sources.Sources.ReadFile("img/bugs/index.gif")

	create, _ := os.Create(outputFile + ".gif")
	defer func() {
		create.Close()
		os.RemoveAll(outputFile + ".gif")
	}()

	io.Copy(create, bytes.NewBuffer(file))

	totalTIme := int(GetVideoTime(videoInput) + 3)
	cellWidth := img.Width / totalTIme
	cmd := exec.Command("ffmpeg",
		"-i", videoInput,
		"-ignore_loop", "0", // 循环播放
		"-i", outputFile+".gif",
		"-filter_complex", fmt.Sprintf("[0:v][1:v] overlay=x='if(lt(-w+(t)*%d,0),0,-w+(t)*%d)':y=H-h", cellWidth, cellWidth),
		"-t", fmt.Sprintf("%d", int(totalTIme)),
		"-y",
		outputFile,
	)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func GetVideoTime(videoInput string) float64 {
	videoDurationCmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoInput)
	audioDuration, _ := videoDurationCmd.Output()

	// 使用 strings.Map 函数去除不可见字符
	result := strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1 // 删除不可见字符
	}, string(audioDuration))

	videoTimes, _ := strconv.ParseFloat(strings.TrimSpace(result), 10)
	return videoTimes
}

func timeToSeconds(timeStr string) (int, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("时间格式无效")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("小时格式错误")
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("分钟格式错误")
	}

	secWithMs := strings.Split(parts[2], ".")
	seconds, err := strconv.Atoi(secWithMs[0])
	if err != nil {
		return 0, fmt.Errorf("秒格式错误")
	}

	totalSeconds := hours*3600 + minutes*60 + seconds
	return totalSeconds, nil
}
