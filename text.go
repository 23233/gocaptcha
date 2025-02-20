package gocaptcha

import (
	"errors"
	"image"
	"image/draw"
	"math"
	"math/rand"
	"time"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

var (
	ErrNilCanvas = errors.New("canvas is nil")
	ErrNilText   = errors.New("text is nil")
)

// TextDrawer is a text drawer interface.
type TextDrawer interface {
	DrawString(canvas draw.Image, text string) error
}

type textDrawer struct {
	dpi float64
	r   *rand.Rand
}

// DrawString draws a string on the canvas.
func (t *textDrawer) DrawString(canvas draw.Image, text string) error {
	if len(text) == 0 {
		return ErrNilText
	}
	if canvas == nil {
		return ErrNilCanvas
	}
	c := freetype.NewContext()
	if t.dpi <= 0 {
		t.dpi = 72
	}
	c.SetDPI(t.dpi)
	c.SetClip(canvas.Bounds())
	c.SetDst(canvas)
	c.SetHinting(font.HintingFull)

	fontWidth := canvas.Bounds().Dx() / len(text)

	for i, s := range text {

		fontSize := float64(canvas.Bounds().Dy()) / (1 + float64(t.r.Intn(7))/float64(9))

		c.SetSrc(image.NewUniform(RandDeepColor()))
		c.SetFontSize(fontSize)
		f, err := DefaultFontFamily.Random()

		if err != nil {
			return err
		}
		c.SetFont(f)

		x := (fontWidth)*i + (fontWidth)/int(fontSize)

		y := 5 + t.r.Intn(canvas.Bounds().Dy()/2) + int(fontSize/2)

		pt := freetype.Pt(x, y)

		_, err = c.DrawString(string(s), pt)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewTextDrawer returns a new text drawer.
func NewTextDrawer(dpi float64) TextDrawer {
	return &textDrawer{
		dpi: dpi,
		r:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type twistTextDrawer struct {
	dpi       float64
	r         *rand.Rand
	amplitude float64
	frequency float64
}

// DrawString draws a string on the canvas.
func (t *twistTextDrawer) DrawString(canvas draw.Image, text string) error {
	if len(text) == 0 {
		return ErrNilText
	}
	if canvas == nil {
		return ErrNilCanvas
	}

	bounds := canvas.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建一个新的画布用于存储扭曲后的图像
	textCanvas := image.NewRGBA(bounds)
	draw.Draw(textCanvas, textCanvas.Bounds(), image.Transparent, image.Point{}, draw.Src)

	c := freetype.NewContext()
	if t.dpi <= 0 {
		t.dpi = 72
	}
	c.SetDPI(t.dpi)
	c.SetClip(bounds)
	c.SetDst(textCanvas)
	c.SetHinting(font.HintingFull)

	// 计算每个字符的最大宽度，预留边距
	fontWidth := (width - 20) / len(text) // 左右各预留10像素边距

	// 计算字体大小范围
	maxFontSize := float64(height) * 0.8 // 使用80%的高度作为最大字体大小
	// 提高最小字体大小的比例，确保文字足够大
	minFontSize := float64(height) * 0.65 // 使用65%的高度作为最小字体大小

	// 确保最小字体大小不会因为宽度限制而变得太小
	if float64(fontWidth)*0.9 > minFontSize {
		minFontSize = float64(fontWidth) * 0.9 // 使用90%的字符宽度作为最小值
	}

	for i, s := range text {
		// 基准字体大小设置为最小字体大小
		baseFontSize := minFontSize
		// 只允许向上浮动，不允许比最小值更小
		fontSize := baseFontSize * (1.0 + float64(t.r.Intn(15))/100.0) // 在100%-115%之间随机
		if fontSize > maxFontSize {
			fontSize = maxFontSize
		}

		c.SetSrc(image.NewUniform(RandDeepColor()))
		c.SetFontSize(fontSize)
		f, err := DefaultFontFamily.Random()
		if err != nil {
			return err
		}
		c.SetFont(f)

		// 计算文字位置
		x := 10 + fontWidth*i + (fontWidth-int(fontSize))/2 // 居中对齐
		// 垂直方向上的位置，确保不会超出边界
		baseY := height/2 + int(fontSize/2)         // 基准位置在中间
		maxOffset := height/2 - int(fontSize/2) - 5 // 最大偏移量
		if maxOffset < 0 {
			maxOffset = 0
		}
		y := baseY + t.r.Intn(2*maxOffset+1) - maxOffset // 在允许范围内随机偏移

		pt := freetype.Pt(x, y)
		_, err = c.DrawString(string(s), pt)
		if err != nil {
			return err
		}
	}

	return t.twistEffect(textCanvas, canvas)
}

func (t *twistTextDrawer) twistEffect(src image.Image, dst draw.Image) error {
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()

	// 遍历源图像像素
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 计算扭曲后的坐标
			dx := int(t.amplitude * math.Sin(t.frequency*float64(y)))
			newX := x + dx
			newY := y

			// 如果新坐标在目标图像范围内，设置像素
			if newX >= 0 && newX < width && newY >= 0 && newY < height {
				_, _, _, a := src.At(x, y).RGBA()
				if a != 0 {
					dst.Set(newX, newY, src.At(x, y))
				}
			}
		}
	}
	return nil
}

// NewTwistTextDrawer returns a new text drawer with twist effect.
func NewTwistTextDrawer(dpi float64, amplitude float64, frequency float64) TextDrawer {
	return &twistTextDrawer{
		dpi:       dpi,
		r:         rand.New(rand.NewSource(time.Now().UnixNano())),
		amplitude: amplitude,
		frequency: frequency,
	}
}
