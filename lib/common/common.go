package common

import "math"

type Setting struct {
	FpsRate         float64
	MaxTime         float64
	MusicRule       string
	HighPerformance bool
}

type PageData struct {
	Title   string   `mapstructure:"title"`
	Content []string `mapstructure:"content"`
	Style   Style    `mapstructure:"style"`
}

type Style struct {
	Title      TitleStyle   `mapstructure:"title"`
	Content    ContentStyle `mapstructure:"content"`
	Background string       `mapstructure:"background"`
	LiveTime   int          `mapstructure:"live_time"`
}

type TitleStyle struct {
	Align string  `mapstructure:"align"`
	Size  float64 `mapstructure:"size"`
	Color *Color  `mapstructure:"color"`
}
type ContentStyle struct {
	Align string  `mapstructure:"align"`
	Size  float64 `mapstructure:"size"`
	Color *Color  `mapstructure:"color"`
}

type Color struct {
	R int
	G int
	B int
}

var (
	FpsCount   = 24       // 每幅图帧率
	Black      = 4        // 留白
	Start      = 12       // 开场透明结束帧
	End        = FpsCount // 结束透明开始帧
	JumpHeight = 6        // 进度条跳的高度
	JumpRate   = 6        // 进度条跳的频率，每JumpRate帧完成一次跳跃
	WalkRate   = 1        // 进度条步行的频率，每WalkRate帧完成一次跳跃
)

type VideoConfig struct {
	FpsCount   int `json:"fps_count"`   // 每页帧率
	Black      int `json:"black"`       // 留白
	Start      int `json:"start"`       // 开场透明结束帧
	End        int `json:"end"`         // 结束透明开始帧
	JumpHeight int `json:"jump_height"` // 进度条跳的高度
	JumpRate   int `json:"jump_rate"`   // 进度条跳的频率，每JumpRate帧完成一次跳跃
	WalkRate   int `json:"walk_rate"`   // 进度条步行的频率，每WalkRate帧完成一次跳跃
}

func GetConfig(setting *Setting, data PageData) VideoConfig {
	if setting != nil && setting.FpsRate != 0 {
		FpsCount = int(math.Ceil(setting.FpsRate)) * len(data.Content)
		End = FpsCount
	}

	if data.Style.LiveTime != 0 && setting.FpsRate != 0 {
		FpsCount = int(float64(data.Style.LiveTime) * setting.FpsRate)
		End = FpsCount
	}

	return VideoConfig{
		FpsCount:   FpsCount,
		Black:      Black,
		Start:      Start,
		End:        End,
		JumpHeight: JumpHeight,
		JumpRate:   JumpRate,
		WalkRate:   WalkRate,
	}
}
