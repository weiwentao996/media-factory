package video

import (
	"fmt"
	"io"
	"os"
	"os/exec"
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

func AddBgm(videoPath, audioPath, outputPath string) error {
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
		cmdArgs = append(cmdArgs, "-stream_loop", "-1", "-i", audioPath, "-filter_complex", "[1:a]volume=0.3[a1];[0:a][a1]amix=inputs=2:duration=first[a]", "-strict", "experimental")
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
