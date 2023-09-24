package voice

import (
	"fmt"
	"testing"
)

func TestGenEdgeVoice(t *testing.T) {
	GenEdgeVoice([]string{"1"}, "./1.wav")
}

func TestMergeMp3(t *testing.T) {
	MergeWAV("./*.wav", "./output.wav")
}
func TestGetWavDuration(t *testing.T) {
	duration, _ := GetWavDuration("./0000.wav")
	fmt.Println(duration)
}

func TestReadVtt(t *testing.T) {
	vtt, err := ReadVtt("./0.wav.vtt")
	if err != nil {
		t.Failed()
	}
	fmt.Println(vtt)
}
