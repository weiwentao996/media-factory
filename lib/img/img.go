package img

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/weiwentao996/media-factory/lib/common"
	"github.com/weiwentao996/media-factory/sources"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"math"
	"net/http"
	"strings"
	"sync"
)

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

var ContentColors = []common.Color{
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

func LoadFontFace(fontBytes []byte, points float64) (font.Face, error) {
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face, nil
}

func GenPPTImage(outPath string, data common.PageData, counter, allFpsCount int, setting *common.PPTSetting) {
	conf := common.GetConfig(setting, data)
	dc := gg.NewContext(Width, Height)
	offsetY := 10.0
	// 标题
	if data.Style.Title.Size <= 0 {
		data.Style.Title.Size = 80
	}

	titleFrontBytes, err := sources.Sources.ReadFile("fronts/Aa厚底黑.ttf")
	if err != nil {
		panic(err)
	}

	face, err := LoadFontFace(titleFrontBytes, data.Style.Title.Size)
	if err != nil {
		panic(err)
	}

	dc.SetFontFace(face)

	if data.Style.Title.Color != nil {
		dc.SetRGB255(data.Style.Title.Color.R, data.Style.Title.Color.G, data.Style.Title.Color.B)
	}

	sWidth, sHeight := dc.MeasureString(data.Title)
	var x, y float64
	switch strings.ToLower(data.Style.Title.Align) {
	case "left":
		x, y = 0, sHeight+offsetY
	case "right":
		x, y = Width-sWidth, sHeight+offsetY
	case "center":
		x, y = (Width-sWidth)/2, sHeight+offsetY
	default:
		x, y = (Width-sWidth)/2, sHeight+offsetY
	}

	rectColor := color.RGBA{0, 0, 0, 120} // 背景色
	dc.SetColor(rectColor)
	dc.DrawRectangle(x, y-0.9*sHeight, sWidth, sHeight*1.2)
	dc.Fill()

	if data.Style.Title.Color == nil {
		dc.SetRGB255(237, 90, 101)
	}
	dc.DrawString(data.Title, x, y)

	offsetY += sHeight * 3
	if data.Content != nil {
		if data.Style.Content.Size <= 0 {
			data.Style.Content.Size = 60
		}

		if data.Style.Content.Color != nil {
			dc.SetRGB255(data.Style.Content.Color.R, data.Style.Content.Color.G, data.Style.Content.Color.B)
		}

		contentFrontBytes, err := sources.Sources.ReadFile("fronts/Leefont蒙黑体.ttf")
		if err != nil {
			panic(err)
		}

		contentFace, err := LoadFontFace(contentFrontBytes, data.Style.Content.Size)
		if err != nil {
			panic(err)
		}

		dc.SetFontFace(contentFace)

		for i, c := range data.Content {
			cWidth, cHeight := dc.MeasureString(c)
			var x, y float64
			switch strings.ToLower(data.Style.Content.Align) {
			case "left":
				x, y = 0, cWidth+cHeight
			case "right":
				x, y = Width-cWidth, cHeight+offsetY
			case "center":
				x, y = (Width-cWidth)/2, cHeight+offsetY
			default:
				x, y = (Width-cWidth)/2, cHeight+offsetY
			}

			dc.SetColor(rectColor)
			dc.DrawRectangle(x, y-0.9*cHeight, cWidth, cHeight*1.2)
			dc.Fill()

			if data.Style.Content.Color == nil {
				colorIndex := i % len(ContentColors)
				dc.SetRGB255(ContentColors[colorIndex].R, ContentColors[colorIndex].G, ContentColors[colorIndex].B)
			}
			dc.DrawString(c, x, y)

			offsetY += cHeight * 2
		}
	}

	imageBytes, err := sources.Sources.ReadFile("img/BG.png")
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(imageBytes)
	bg, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}

	var bgImg image.Image
	if data.Style.Background != "" {
		bgImg, err = GetImage(data.Style.Background)
		if err != nil {
			panic(err)
		}
		//bgImg = adjustOpacity(bgImg, 0.3)
	}

	if err != nil {
		panic(err)
	}

	processList := NewCircularLinkedList()

	for i := 0; i < 12; i++ {
		imageBytes, err = sources.Sources.ReadFile(fmt.Sprintf("img/bugs/process%d.png", i))
		if err != nil {
			panic(err)
		}
		reader = bytes.NewReader(imageBytes)
		p, _, err := image.Decode(reader)
		if err != nil {
			panic(err)
		}

		processList.Insert(p)
	}

	if setting.HighPerformance {
		// 多线程
		wg := sync.WaitGroup{}
		wg.Add(conf.FpsCount)
		for i := 0; i < conf.FpsCount; i++ {
			var proImage image.Image
			proImage = processList.GetProcess()
			go func(fpsIndex int, counter int) {
				bgc := gg.NewContextForImage(bg)
				if bgImg != nil {
					putBackGroundImage(bgc, bgImg)
				}
				bgc.DrawImage(proImage, int(float64(Width)*(float64(counter)/float64(allFpsCount))), Height-(128+(conf.JumpRate-fpsIndex%conf.JumpRate)*(conf.JumpHeight/conf.JumpRate)))
				switch {
				case fpsIndex < conf.Black:
				case fpsIndex < conf.Start:
					img := adjustOpacity(dc.Image(), float64(fpsIndex+1)/float64(conf.Start))
					bgc.DrawImage(img, 0, int(sHeight))
				case fpsIndex > conf.End:
					img := adjustOpacity(dc.Image(), float64(1)-float64(fpsIndex%conf.End)/float64(conf.FpsCount-conf.End))
					bgc.DrawImage(img, 0, int(sHeight))
				default:
					bgc.DrawImage(dc.Image(), 0, int(sHeight))
				}
				fileName := fmt.Sprintf("%s/%05d.png", outPath, counter)
				if err := bgc.SavePNG(fileName); err != nil {
					panic(err)
				}
				wg.Done()
			}(i, counter)
			counter++
		}

		wg.Wait()
	} else {
		for i := 0; i < conf.FpsCount; i++ {
			proImage := processList.GetProcess()
			bgc := gg.NewContextForImage(bg)
			if bgImg != nil {
				putBackGroundImage(bgc, bgImg)
			}
			bgc.DrawImage(proImage, int(float64(Width)*(float64(counter+i)/float64(allFpsCount))), Height-(128+(conf.JumpRate-i%conf.JumpRate)*(conf.JumpHeight/conf.JumpRate)))
			switch {
			case i < conf.Black:
			case i < conf.Start:
				img := adjustOpacity(dc.Image(), float64(i+1)/float64(conf.Start))
				bgc.DrawImage(img, 0, int(sHeight))
			case i > conf.End:
				img := adjustOpacity(dc.Image(), float64(1)-float64(i%conf.End)/float64(conf.FpsCount-conf.End))
				bgc.DrawImage(img, 0, int(sHeight))
			default:
				bgc.DrawImage(dc.Image(), 0, int(sHeight))
			}
			fileName := fmt.Sprintf("%s/%05d.png", outPath, counter)
			if err := bgc.SavePNG(fileName); err != nil {
				panic(err)
			}
			counter++
		}
	}
}

func GetImage(path string) (image.Image, error) {
	split := strings.Split(path, ":")
	if len(split) > 0 && (split[0] == "http" || split[0] == "https") {
		return getImageFromNet(path)
	}
	imageBytes, err := sources.Sources.ReadFile("img/BG.png")
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(imageBytes)
	bg, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}
	return bg, nil
}

func GenAdviceImage(outPath string, data *common.VttContent, videoEndTime float64, counter int, setting *common.AdviceFoSetting, style *common.AdviceFoStyle) int {
	dc := gg.NewContext(Width, Height)

	titleFrontBytes, err := sources.Sources.ReadFile("fronts/Aa厚底黑.ttf")
	if err != nil {
		panic(err)
	}

	face, err := LoadFontFace(titleFrontBytes, style.Size)
	if err != nil {
		panic(err)
	}

	dc.SetFontFace(face)

	if style.Color != nil {
		dc.SetRGB255(style.Color.R, style.Color.G, style.Color.B)
	}

	sWidth, sHeight := dc.MeasureString(data.Content)

	offsetY := Height - sHeight - 200
	var x, y float64
	switch strings.ToLower(style.Align) {
	case "left":
		x, y = 0, sHeight+offsetY
	case "right":
		x, y = Width-sWidth, sHeight+offsetY
	case "center":
		x, y = (Width-sWidth)/2, sHeight+offsetY
	default:
		x, y = (Width-sWidth)/2, sHeight+offsetY
	}

	rectColor := color.RGBA{249, 205, 173, 100} // 背景色
	dc.SetColor(rectColor)
	dc.DrawRectangle(x, y-0.9*sHeight, sWidth, sHeight*1.2)
	dc.Fill()

	if style.Color == nil {
		dc.SetRGB255(237, 90, 101)
	}

	if style.Color != nil {
		dc.SetRGB255(style.Color.R, style.Color.G, style.Color.B)
	}

	dc.DrawString(data.Content, x, y)

	imageBytes, err := sources.Sources.ReadFile("img/BG.png")
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(imageBytes)
	bg, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}

	var bgImg image.Image
	if style.Background != "" {
		bgImg, err = GetImage(style.Background)
		if err != nil {
			panic(err)
		}
		//bgImg = adjustOpacity(bgImg, 0.3)
	}

	var conetentImage image.Image
	if style.Background != "" {
		conetentImage, err = GetImage(data.ContentImage)
		if err != nil {
			panic(err)
		}
		//bgImg = adjustOpacity(bgImg, 0.3)
	}

	processList := NewCircularLinkedList()

	for i := 0; i < 12; i++ {
		imageBytes, err = sources.Sources.ReadFile(fmt.Sprintf("img/bugs/process%d.png", i))
		if err != nil {
			panic(err)
		}
		reader = bytes.NewReader(imageBytes)
		p, _, err := image.Decode(reader)
		if err != nil {
			panic(err)
		}

		processList.Insert(p)
	}

	fpsCount := int(math.Round((data.Time[1] - data.Time[0] + setting.FpsFix) * setting.FpsRate))
	// 多线程
	wg := sync.WaitGroup{}
	wg.Add(fpsCount)
	for i := 0; i < fpsCount; i++ {
		var proImage image.Image
		proImage = processList.GetProcess()
		allProcessPercentage := (data.Time[0] + (float64(i) / setting.FpsRate)) / videoEndTime
		go func(pageIndex int, process float64) {
			bgc := gg.NewContextForImage(bg)
			if bgImg != nil {
				putBackGroundImage(bgc, bgImg)
			}
			if conetentImage != nil {
				putContentImage(bgc, conetentImage)
			}
			bgc.DrawImage(proImage, int(float64(Width)*process), Height-128)
			bgc.DrawImage(dc.Image(), 0, int(sHeight))
			fileName := fmt.Sprintf("%s/%05d.png", outPath, pageIndex)
			if err := bgc.SavePNG(fileName); err != nil {
				panic(err)
			}
			wg.Done()
		}(counter, allProcessPercentage)
		counter++
	}
	wg.Wait()
	return counter
}

func splitString(input string, segmentCount int) []string {
	// 计算每个段的长度
	strLen := []rune(input)
	segmentSize := int(math.Ceil(float64(len(strLen)) / float64(segmentCount)))

	var rs []string
	current := 0
	end := 0
	for {
		end += segmentSize
		if end > len(strLen) {
			rs = append(rs, string(strLen[current:]))
			break
		} else {
			rs = append(rs, string(strLen[current:end]))
			current = end
		}

	}

	return rs
}

var (
	dcOffsetY       = 240.0
	dcPaddingX      = 400.0
	dcPaddingBottom = 20.0
	avatarSize      = 100.0
)

func GenMusicImageFast(imgPath string, data *common.VttContent, style *common.AdviceFoStyle) {
	dc := gg.NewContext(Width, Height)
	titleFrontBytes, err := sources.Sources.ReadFile("fronts/Aa厚底黑.ttf")
	if err != nil {
		panic(err)
	}

	face, err := LoadFontFace(titleFrontBytes, style.Size)
	if err != nil {
		panic(err)
	}

	dc.SetFontFace(face)

	if style.Color != nil {
		dc.SetRGB255(style.Color.R, style.Color.G, style.Color.B)
	}
	rectColor := color.RGBA{255, 255, 255, 100} // 背景色

	data.Content = strings.TrimSpace(data.Content)

	sWidth, sHeight := dc.MeasureString(data.Content)
	rowCount := sWidth / (Width - dcPaddingX)

	contentArr := splitString(data.Content, int(math.Ceil(rowCount)))

	offsetY := dcOffsetY
	// 昵称
	_, cHeight := dc.MeasureString(data.Nickname)
	var x, y = dcPaddingX / 2, cHeight + offsetY
	dc.SetColor(rectColor)
	//dc.DrawRectangle(x, y-0.9*sHeight, cWidth, sHeight*1.2)
	dc.DrawRectangle(dcPaddingX/2, y-avatarSize, Width-dcPaddingX, dcPaddingBottom+avatarSize)
	dc.Fill()

	if style.Color == nil {
		dc.SetRGB255(0, 0, 0)
	}

	if style.Color != nil {
		dc.SetRGB255(style.Color.R, style.Color.G, style.Color.B)
	}

	if data.Nickname != "" {
		dc.DrawString(fmt.Sprintf("%s  [%s]: ", data.Nickname, data.CommentTime.Format("2006-01-02 15:04:05")), x+avatarSize, y)
	}

	offsetY += cHeight + dcPaddingBottom
	for _, s := range contentArr {
		cWidth, cHeight := dc.MeasureString(s)
		var x, y = (Width - cWidth) / 2, cHeight + offsetY
		dc.SetColor(rectColor)
		//dc.DrawRectangle(x, y-0.9*sHeight, cWidth, sHeight*1.2)
		dc.DrawRectangle(dcPaddingX/2, y-sHeight, Width-dcPaddingX, sHeight+dcPaddingBottom)
		dc.Fill()

		if style.Color == nil {
			dc.SetRGB255(0, 0, 0)
		}

		if style.Color != nil {
			dc.SetRGB255(style.Color.R, style.Color.G, style.Color.B)
		}

		dc.DrawString(s, x, y)

		offsetY += cHeight + dcPaddingBottom
	}

	imageBytes, err := sources.Sources.ReadFile("img/BG.png")
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(imageBytes)
	bg, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}

	var bgImg image.Image
	if style.Background != "" {
		bgImg, err = GetImage(style.Background)
		if err != nil {
			panic(err)
		}
	}

	var avatar image.Image
	if data.Avatar != "" {
		avatar, err = GetImage(data.Avatar)
		if err != nil {
			panic(err)
		}
	}

	bgc := gg.NewContextForImage(bg)

	if bgImg != nil {
		putBackGroundImage(bgc, bgImg)
	}

	if avatar != nil {
		putAvatarImage(dc, avatar)
	}

	bgc.DrawImage(dc.Image(), 0, int(sHeight))

	if err := bgc.SavePNG(imgPath); err != nil {
		panic(err)
	}
}

func GenMusicImage(outPath string, data *common.VttContent, videoEndTime float64, counter int, setting *common.AdviceFoSetting, style *common.AdviceFoStyle) int {
	dc := gg.NewContext(Width, Height)

	titleFrontBytes, err := sources.Sources.ReadFile("fronts/Aa厚底黑.ttf")
	if err != nil {
		panic(err)
	}

	face, err := LoadFontFace(titleFrontBytes, style.Size)
	if err != nil {
		panic(err)
	}

	dc.SetFontFace(face)

	if style.Color != nil {
		dc.SetRGB255(style.Color.R, style.Color.G, style.Color.B)
	}
	rectColor := color.RGBA{255, 255, 255, 100} // 背景色

	sWidth, sHeight := dc.MeasureString(data.Content)
	rowCount := sWidth / (Width - dcPaddingX)

	contentArr := splitString(data.Content, int(math.Ceil(rowCount)))

	offsetY := dcOffsetY
	// 昵称
	_, cHeight := dc.MeasureString(data.Nickname)
	var x, y = dcPaddingX / 2, cHeight + offsetY
	dc.SetColor(rectColor)
	//dc.DrawRectangle(x, y-0.9*sHeight, cWidth, sHeight*1.2)
	dc.DrawRectangle(dcPaddingX/2, y-avatarSize, Width-dcPaddingX, dcPaddingBottom+avatarSize)
	dc.Fill()

	if style.Color == nil {
		dc.SetRGB255(0, 0, 0)
	}

	if style.Color != nil {
		dc.SetRGB255(style.Color.R, style.Color.G, style.Color.B)
	}

	if data.Nickname != "" {
		dc.DrawString(fmt.Sprintf("%s  [%s]: ", data.Nickname, data.CommentTime.Format("2006-01-02 15:04:05")), x+avatarSize, y)
	}

	offsetY += cHeight + dcPaddingBottom

	for _, s := range contentArr {
		cWidth, cHeight := dc.MeasureString(s)
		var x, y = (Width - cWidth) / 2, cHeight + offsetY
		dc.SetColor(rectColor)
		//dc.DrawRectangle(x, y-0.9*sHeight, cWidth, sHeight*1.2)
		dc.DrawRectangle(dcPaddingX/2, y-sHeight, Width-dcPaddingX, sHeight+dcPaddingBottom)
		dc.Fill()

		if style.Color == nil {
			dc.SetRGB255(0, 0, 0)
		}

		if style.Color != nil {
			dc.SetRGB255(style.Color.R, style.Color.G, style.Color.B)
		}

		dc.DrawString(s, x, y)

		offsetY += cHeight + dcPaddingBottom
	}

	imageBytes, err := sources.Sources.ReadFile("img/BG.png")
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(imageBytes)
	bg, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}

	var bgImg image.Image
	if style.Background != "" {
		bgImg, err = GetImage(style.Background)
		if err != nil {
			panic(err)
		}
	}

	var avatar image.Image
	if data.Avatar != "" {
		avatar, err = GetImage(data.Avatar)
		if err != nil {
			panic(err)
		}
	}

	processList := NewCircularLinkedList()

	for i := 0; i < 12; i++ {
		imageBytes, err = sources.Sources.ReadFile(fmt.Sprintf("img/bugs/process%d.png", i))
		if err != nil {
			panic(err)
		}
		reader = bytes.NewReader(imageBytes)
		p, _, err := image.Decode(reader)
		if err != nil {
			panic(err)
		}

		processList.Insert(p)
	}

	fpsCount := int(math.Round((data.Time[1] - data.Time[0] + setting.FpsFix) * setting.FpsRate))
	// 多线程
	wg := sync.WaitGroup{}
	wg.Add(fpsCount)
	for i := 0; i < fpsCount; i++ {
		var proImage image.Image
		proImage = processList.GetProcess()
		allProcessPercentage := (data.Time[0] + (float64(i) / setting.FpsRate)) / videoEndTime
		go func(pageIndex int, process float64) {
			bgc := gg.NewContextForImage(bg)
			if bgImg != nil {
				putBackGroundImage(bgc, bgImg)
			}

			if avatar != nil {
				putAvatarImage(dc, avatar)
			}

			bgc.DrawImage(proImage, int(float64(Width)*process), Height-128)
			bgc.DrawImage(dc.Image(), 0, int(sHeight))
			fileName := fmt.Sprintf("%s/%05d.png", outPath, pageIndex)
			if err := bgc.SavePNG(fileName); err != nil {
				panic(err)
			}
			wg.Done()
		}(counter, allProcessPercentage)
		counter++
	}
	wg.Wait()
	return counter
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

func putBackGroundImage(dc *gg.Context, newBg image.Image) error {

	// 设置背景颜色
	//dc.SetColor(color.White)
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

func putContentImage(dc *gg.Context, img image.Image) error {

	// 计算图片缩放比例
	scale := float64(Height-200) / float64(img.Bounds().Dy())

	// 计算缩放后图片的宽度和高度
	scaledWidth := float64(img.Bounds().Dx()) * scale
	scaledHeight := float64(img.Bounds().Dy()) * scale

	// 计算图片在背景中居中的位置
	xOffset := (float64(Width) - scaledWidth) / 2

	// 创建一个新的绘图上下文
	newBgC := gg.NewContext(int(scaledWidth), int(scaledHeight))

	// 缩放并绘制图像
	newBgC.Scale(scale, scale)
	newBgC.DrawImage(img, 0, 0)

	// 将图片绘制到背景中央
	dc.DrawImage(newBgC.Image(), int(xOffset), 0)

	return nil
}

func putAvatarImage(dc *gg.Context, img image.Image) error {

	// 计算图片缩放比例
	scale := avatarSize / float64(img.Bounds().Dy())

	// 计算缩放后图片的宽度和高度
	scaledWidth := float64(img.Bounds().Dx()) * scale
	scaledHeight := float64(img.Bounds().Dy()) * scale

	// 创建一个新的绘图上下文
	avatar := gg.NewContext(int(scaledWidth), int(scaledHeight))

	// 缩放并绘制图像
	avatar.Scale(scale, scale)
	avatar.DrawImage(img, 0, 0)

	// 将图片绘制到背景
	dc.DrawImage(avatar.Image(), int(dcPaddingX/2)+4, int(dcOffsetY-avatarSize/2+2))

	return nil
}
