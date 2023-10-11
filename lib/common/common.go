package common

import (
	"math"
	"time"
)

type PPTSetting struct {
	FpsRate         float64 `json:"fps_rate"`
	MaxTime         float64 `json:"max_time"`
	MusicRule       string  `json:"music_rule"`
	HighPerformance bool    `json:"high_performance"`
}

type PageData struct {
	Title   string   `json:"title" mapstructure:"title"`
	Content []string `json:"content" mapstructure:"content"`
	Style   PPTStyle `json:"style" mapstructure:"style"`
}

type VttContent struct {
	Content      string     `json:"content" mapstructure:"content"`
	Avatar       string     `json:"avatar" mapstructure:"avatar"`
	Nickname     string     `json:"nickname" mapstructure:"nickname"`
	ContentImage string     `json:"content_image" mapstructure:"content_image"`
	Time         [2]float64 `json:"time" mapstructure:"time"`
	CommentTime  time.Time  `json:"comment_time" mapstructure:"comment_time"`
}

type AdviceFoSetting struct {
	BgmUrl  string  `json:"bgm_url"`
	Model   string  `json:"model"`
	FpsRate float64 `json:"fps_rate"`
	FpsFix  float64 `json:"fps_fix"`
}

type CommentContent struct {
	Content      []string   `json:"content" mapstructure:"content"`
	ContentImage string     `json:"content_image" mapstructure:"content_image"`
	Time         [2]float64 `json:"time" mapstructure:"time"`
}

type AdviceFoStyle struct {
	Align      string  `json:"align" mapstructure:"align"`
	Size       float64 `json:"size" mapstructure:"size"`
	Color      *Color  `json:"color" mapstructure:"color"`
	Background string  `json:"background" mapstructure:"background"`
}

type PPTStyle struct {
	Title      PPTTitleStyle   `json:"title" mapstructure:"title"`
	Content    PPTContentStyle `json:"content" mapstructure:"content"`
	Background string          `json:"background" mapstructure:"background"`
	LiveTime   float64         `json:"live_time" mapstructure:"live_time"`
}

type PPTTitleStyle struct {
	Align string  `json:"align" mapstructure:"align"`
	Size  float64 `json:"size" mapstructure:"size"`
	Color *Color  `json:"color" mapstructure:"color"`
}
type PPTContentStyle struct {
	Align string  `json:"align" mapstructure:"align"`
	Size  float64 `json:"size" mapstructure:"size"`
	Color *Color  `json:"color" mapstructure:"color"`
}

type Color struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

var (
	FpsCount   = 24       // 每幅图帧率
	Black      = 0        // 留白
	Start      = 6        // 开场透明结束帧
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

func GetConfig(setting *PPTSetting, data PageData) VideoConfig {
	if setting != nil && setting.FpsRate != 0 {
		FpsCount = int(math.Ceil(setting.FpsRate)) * len(data.Content)
		End = FpsCount
	}

	if data.Style.LiveTime != 0 && setting.FpsRate != 0 {
		FpsCount = int(data.Style.LiveTime * setting.FpsRate)
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
