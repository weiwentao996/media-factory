package img

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/weiwentao996/media-factory/lib/common"
	"github.com/weiwentao996/media-factory/sources"
	"image"
	"image/color"
	"math"
	"strings"
)

// ------------ 生成图片 ----------

type ImageData struct {
	Title   string   `mapstructure:"title"`
	Content []string `mapstructure:"content"`
	Style   Style    `mapstructure:"style"`
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

type Style struct {
	Title   TitleStyle   `mapstructure:"title"`
	Content ContentStyle `mapstructure:"content"`
}

type Color struct {
	R int
	G int
	B int
}

const Width = 1920
const Height = 1080

var ContentColors = []Color{
	{
		R: 254,
		G: 186,
		B: 7,
	},
	{
		R: 248,
		G: 179,
		B: 127,
	},
	{
		R: 199,
		G: 237,
		B: 204,
	},
}

var (
	FpsCount   = 24       // 每幅图帧率
	Black      = 4        // 留白
	Start      = 12       // 开场透明结束帧
	End        = FpsCount // 结束透明开始帧
	JumpHeight = 4        // 进度条跳的高度
	JumpRate   = 6        // 进度条跳的频率，每JumpRate帧完成一次跳跃
	WalkRate   = 1        // 进度条步行的频率，每WalkRate帧完成一次跳跃
)

func GenImage(outPath string, data ImageData, currentLoop, imageCount int, setting *common.Setting) {
	if setting != nil && setting.FpsCount > FpsCount {
		FpsCount = setting.FpsCount
		End = setting.FpsCount
	}

	dc := gg.NewContext(Width, Height)
	offsetY := 10.0

	// 标题
	if data.Style.Title.Size <= 0 {
		data.Style.Title.Size = 80
	}

	if err := dc.LoadFontFace(sources.Path+"/front/Aa厚底黑.ttf", data.Style.Title.Size); err != nil {
		panic(err)
	}

	if data.Style.Title.Color != nil {
		dc.SetRGB255(data.Style.Title.Color.R, data.Style.Title.Color.G, data.Style.Title.Color.B)
	}

	if data.Style.Title.Color == nil {
		dc.SetRGB255(237, 90, 101)
	}

	sWidth, sHeight := dc.MeasureString(data.Title)
	switch strings.ToLower(data.Style.Title.Align) {
	case "left":
		dc.DrawString(data.Title, 0, sHeight+offsetY)
	case "right":
		dc.DrawString(data.Title, Width-sWidth, sHeight+offsetY)
	case "center":
		dc.DrawString(data.Title, (Width-sWidth)/2, sHeight+offsetY)
	default:
		dc.DrawString(data.Title, (Width-sWidth)/2, sHeight+offsetY)
	}
	offsetY += sHeight * 2

	if data.Content != nil {
		if data.Style.Content.Size <= 0 {
			data.Style.Content.Size = 60
		}

		if data.Style.Content.Color != nil {
			dc.SetRGB255(data.Style.Content.Color.R, data.Style.Content.Color.G, data.Style.Content.Color.B)
		}

		if err := dc.LoadFontFace(sources.Path+"/front/Leefont蒙黑体.ttf", data.Style.Content.Size); err != nil {
			panic(err)
		}

		for i, c := range data.Content {
			cWidth, cHeight := dc.MeasureString(c)
			if data.Style.Content.Color == nil {
				colorIndex := i % len(ContentColors)
				dc.SetRGB255(ContentColors[colorIndex].R, ContentColors[colorIndex].G, ContentColors[colorIndex].B)
			}

			switch strings.ToLower(data.Style.Content.Align) {
			case "left":
				dc.DrawString(c, 0, cHeight+offsetY)
			case "right":
				dc.DrawString(c, Width-cWidth, cHeight+offsetY)
			case "center":
				dc.DrawString(c, (Width-cWidth)/2, cHeight+offsetY)
			default:
				dc.DrawString(c, (Width-cWidth)/2, cHeight+offsetY)
			}
			offsetY += cHeight * 2
		}
	}

	bg, err := gg.LoadImage(sources.Path + "/img/BG.png")
	if err != nil {
		panic(err)
	}

	var proArr []image.Image

	for i := 0; i < 12; i++ {
		p, err := gg.LoadImage(fmt.Sprintf("%s/img/process%d.png", sources.Path, i))
		if err != nil {
			panic(err)
		}
		proArr = append(proArr, p)
	}

	for i := 0; i < FpsCount; i++ {
		proImage := proArr[int(math.Floor(float64(i+1)/float64(WalkRate)))%len(proArr)]
		bgc := gg.NewContextForImage(bg)
		bgc.DrawImage(proImage, int(float64(Width)*(float64(FpsCount*currentLoop+i)/float64(FpsCount*imageCount))), Height-(128+(JumpRate-i%JumpRate)*(JumpHeight/JumpRate)))
		//bgc.DrawImage(proImage, int(float64(Width)*(float64(FpsCount*currentLoop+i)/float64(FpsCount*imageCount))), Height-128)
		switch {
		case i < Black:
		case i < Start:
			img := adjustOpacity(dc.Image(), float64(i+1)/float64(Start))
			bgc.DrawImage(img, 0, int(sHeight))
		case i > End:
			img := adjustOpacity(dc.Image(), float64(1)-float64(i%End)/float64(FpsCount-End))
			bgc.DrawImage(img, 0, int(sHeight))
		default:
			bgc.DrawImage(dc.Image(), 0, int(sHeight))
		}
		fileName := fmt.Sprintf("%s/%05d.png", outPath, i+(currentLoop*FpsCount))
		if err := bgc.SavePNG(fileName); err != nil {
			panic(err)
		}
	}
}

// adjustOpacity 将输入图像m的透明度变为原来的倍数。若原来为完成全不透明，则percentage = 0.5将变为半透明
func adjustOpacity(m image.Image, percentage float64) image.Image {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA64(bounds)
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			r, g, b, a := colorRgb.RGBA()
			opacity := uint16(float64(a) * percentage)
			v := newRgba.ColorModel().Convert(color.NRGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: opacity})
			_r, _g, _b, _a := v.RGBA()
			newRgba.SetRGBA64(i, j, color.RGBA64{R: uint16(_r), G: uint16(_g), B: uint16(_b), A: uint16(_a)})
		}
	}
	return newRgba
}
