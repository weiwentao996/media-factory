package voice

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"time"
)

func GenVoice(contents []string, path string) string {
	ole.CoInitialize(0)
	unknown, _ := oleutil.CreateObject("SAPI.SpVoice")
	voice, _ := unknown.QueryInterface(ole.IID_IDispatch)
	saveFile, _ := oleutil.CreateObject("SAPI.SpFileStream")
	ff, _ := saveFile.QueryInterface(ole.IID_IDispatch)
	fileName := fmt.Sprintf("%s/voice.wav", path)
	// 打开wav文件
	oleutil.CallMethod(ff, "Open", fileName, 3, true)
	// 设置voice的AudioOutputStream属性，必须是PutPropertyRef，如果是PutProperty就无法生效
	oleutil.PutPropertyRef(voice, "AudioOutputStream", ff)
	// 设置语速
	oleutil.PutProperty(voice, "Rate", 4)
	// 设置音量
	oleutil.PutProperty(voice, "Volume", 200)
	// 说话
	for _, content := range contents {
		oleutil.CallMethod(voice, "Speak", content)
		// 等待5秒钟
		time.Sleep(5 * time.Second)

		// 暂停语音引擎
		oleutil.CallMethod(voice, "Pause")

		// 等待一段时间
		time.Sleep(5 * time.Second)

		// 恢复语音引擎
		oleutil.CallMethod(voice, "Resume")

	}
	// 停止说话
	oleutil.CallMethod(voice, "Pause")
	// 恢复说话
	oleutil.CallMethod(voice, "Resume")
	// 等待结束
	oleutil.CallMethod(voice, "WaitUntilDone", 1000000)
	// 关闭文件
	oleutil.CallMethod(ff, "Close")
	ff.Release()
	voice.Release()
	ole.CoUninitialize()
	return fileName
}
