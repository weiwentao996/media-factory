package img

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/weiwentao996/media-factory/lib/common"
	"github.com/weiwentao996/media-factory/sources"
	"image"
	"image/color"
	"math"
	"net/http"
	"strings"
	"sync"
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
	Title      TitleStyle   `mapstructure:"title"`
	Content    ContentStyle `mapstructure:"content"`
	Background string       `mapstructure:"background"`
	LiveTime   int          `mapstructure:"live_time"`
}

type Color struct {
	R int
	G int
	B int
}

// 定义环形链表
type CircularLinkedList struct {
	Head    *ProcessCircle
	Current *ProcessCircle
	Tail    *ProcessCircle
}

type ProcessCircle struct {
	Value image.Image
	Next  *ProcessCircle
}

// 初始化一个空的环形链表
func NewCircularLinkedList() *CircularLinkedList {
	return &CircularLinkedList{}
}

// 插入节点到链表尾部
func (cll *CircularLinkedList) Insert(data image.Image) {
	newNode := &ProcessCircle{Value: data}
	if cll.Head == nil {
		cll.Head = newNode
		cll.Tail = newNode
		cll.Current = newNode
		newNode.Next = newNode
	} else {
		newNode.Next = cll.Head
		cll.Tail.Next = newNode
		cll.Tail = newNode
	}
}

func (cll *CircularLinkedList) GetProcess() image.Image {
	rs := cll.Current.Value
	cll.Current = cll.Current.Next
	return rs
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
	FpsCount   = 12       // 每幅图帧率
	Black      = 4        // 留白
	Start      = 12       // 开场透明结束帧
	End        = FpsCount // 结束透明开始帧
	JumpHeight = 6        // 进度条跳的高度
	JumpRate   = 6        // 进度条跳的频率，每JumpRate帧完成一次跳跃
	WalkRate   = 1        // 进度条步行的频率，每WalkRate帧完成一次跳跃
)

func GenImage(outPath string, data ImageData, currentLoop, imageCount int, setting *common.Setting) {
	if setting != nil && setting.FpsRate != 0 {
		FpsCount = int(math.Ceil(setting.FpsRate)) * len(data.Content)
		End = FpsCount
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
	offsetY += sHeight * 4

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

	var bgImg image.Image
	if data.Style.Background != "" {
		bgImg, err = GetImage(data.Style.Background)
		if err != nil {
			panic(err)
		}
		bgImg = adjustOpacity(bgImg, 0.6)
	}

	if err != nil {
		panic(err)
	}

	processList := NewCircularLinkedList()

	for i := 0; i < 12; i++ {
		p, err := gg.LoadImage(fmt.Sprintf("%s/img/bugs/process%d.png", sources.Path, i))
		if err != nil {
			panic(err)
		}

		processList.Insert(p)
	}

	if setting.HighPerformance {
		// 多线程
		wg := sync.WaitGroup{}
		wg.Add(FpsCount)
		for i := 0; i < FpsCount; i++ {
			var proImage image.Image
			proImage = processList.GetProcess()
			go func(fpsIndex int) {
				bgc := gg.NewContextForImage(bg)
				if bgImg != nil {
					putImage(bgc, bgImg)
				}
				bgc.DrawImage(proImage, int(float64(Width)*(float64(FpsCount*currentLoop+fpsIndex)/float64(FpsCount*imageCount))), Height-(128+(JumpRate-fpsIndex%JumpRate)*(JumpHeight/JumpRate)))
				//bgc.DrawImage(proImage, int(float64(Width)*(float64(FpsCount*currentLoop+fpsIndex)/float64(FpsCount*imageCount))), Height-128)
				switch {
				case fpsIndex < Black:
				case fpsIndex < Start:
					img := adjustOpacity(dc.Image(), float64(fpsIndex+1)/float64(Start))
					bgc.DrawImage(img, 0, int(sHeight))
				case fpsIndex > End:
					img := adjustOpacity(dc.Image(), float64(1)-float64(fpsIndex%End)/float64(FpsCount-End))
					bgc.DrawImage(img, 0, int(sHeight))
				default:
					bgc.DrawImage(dc.Image(), 0, int(sHeight))
				}
				fileName := fmt.Sprintf("%s/%05d.png", outPath, fpsIndex+(currentLoop*FpsCount))
				if err := bgc.SavePNG(fileName); err != nil {
					panic(err)
				}
				wg.Done()
			}(i)
		}

		wg.Wait()
	} else {
		for i := 0; i < FpsCount; i++ {
			proImage := processList.GetProcess()
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
}

func GetImage(path string) (image.Image, error) {
	split := strings.Split(path, ":")
	if len(split) > 0 && (split[0] == "http" || split[0] == "https") {
		return getImageFromNet(path)
	}
	return gg.LoadImage(sources.Path + "/img/BG.png")
}

// GetImageFromNet 从远程读取图片
func getImageFromNet(url string) (image.Image, error) {
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		return nil, err
	}
	defer res.Body.Close()
	m, _, err := image.Decode(res.Body)
	return m, err
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

func putImage(dc *gg.Context, newBg image.Image) error {

	// 设置背景颜色
	dc.SetColor(color.White)
	dc.Clear()

	// 计算图片缩放比例
	scaleX := float64(newBg.Bounds().Dy()) / float64(Height)
	scaleY := float64(newBg.Bounds().Dx()) / float64(Width)

	// 选择较小的缩放比例，确保图片完全适应背景
	scale := math.Min(scaleX, scaleY)

	// 计算缩放后图片的宽度和高度
	scaledWidth := float64(newBg.Bounds().Dx()) * scale
	scaledHeight := float64(newBg.Bounds().Dy()) * scale

	// 计算图片在背景中居中的位置
	xOffset := (float64(Width) - scaledWidth) / 2
	yOffset := (float64(Height) - scaledHeight) / 2

	// 创建一个新的绘图上下文
	newBgC := gg.NewContext(int(scaledWidth), int(scaledHeight))

	// 缩放并绘制图像
	newBgC.Scale(scaleX, scaleY)
	newBgC.DrawImage(newBg, 0, 0)

	// 将图片绘制到背景中央
	dc.DrawImage(newBgC.Image(), int(xOffset), int(yOffset))

	return nil
}
