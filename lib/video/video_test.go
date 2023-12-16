package video

import "testing"

func TestMerge(t *testing.T) {
	Merge("C:\\work\\media-factory\\output\\*.mp4", "C:\\work\\media-factory\\output\\output.mp4")
}

func TestAddProcess(t *testing.T) {
	AddProcess("", "")
}
