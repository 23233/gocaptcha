package gocaptcha

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
)

const (
	// DefaultDPI 默认的dpi
	DefaultDPI = 72.0
	// DefaultBlurKernelSize 默认模糊卷积核大小
	DefaultBlurKernelSize = 2
	// DefaultBlurSigma 默认模糊sigma值
	DefaultBlurSigma = 0.65
	// DefaultAmplitude 默认图片扭曲的振幅
	DefaultAmplitude = 20
	//DefaultFrequency 默认图片扭曲的波频率
	DefaultFrequency = 0.05
)

var TextCharacters = []rune("ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz0123456789")

const (
	ImageFormatPng ImageFormat = iota
	ImageFormatJpeg
	ImageFormatGif
)

// ImageFormat 图片格式
type ImageFormat int

type CaptchaImage struct {
	nrgba   *image.NRGBA
	width   int
	height  int
	Complex int
	Error   error
}

// New 新建一个图片对象
func New(width int, height int, bgColor color.RGBA) *CaptchaImage {
	m := image.NewNRGBA(image.Rect(0, 0, width, height))

	draw.Draw(m, m.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	return &CaptchaImage{
		nrgba:  m,
		height: height,
		width:  width,
	}
}

// Encode 编码图片
func (captcha *CaptchaImage) Encode(w io.Writer, imageFormat ImageFormat) error {

	if imageFormat == ImageFormatPng {
		return png.Encode(w, captcha.nrgba)
	}
	if imageFormat == ImageFormatJpeg {
		return jpeg.Encode(w, captcha.nrgba, &jpeg.Options{Quality: 100})
	}
	if imageFormat == ImageFormatGif {
		return gif.Encode(w, captcha.nrgba, &gif.Options{NumColors: 256})
	}

	return errors.New("not supported image format")
}

// DrawLine 画直线.
func (captcha *CaptchaImage) DrawLine(drawer LineDrawer, lineColor color.Color) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	y := captcha.nrgba.Bounds().Dy()
	point1 := image.Point{X: captcha.nrgba.Bounds().Min.X + 1, Y: rand.Intn(y)}
	point2 := image.Point{X: captcha.nrgba.Bounds().Max.X - 1, Y: rand.Intn(y)}
	captcha.Error = drawer.DrawLine(captcha.nrgba, point1, point2, lineColor)
	return captcha
}

// DrawBorder 画边框.
func (captcha *CaptchaImage) DrawBorder(borderColor color.RGBA) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	for x := 0; x < captcha.width; x++ {
		captcha.nrgba.Set(x, 0, borderColor)
		captcha.nrgba.Set(x, captcha.height-1, borderColor)
	}
	for y := 0; y < captcha.height; y++ {
		captcha.nrgba.Set(0, y, borderColor)
		captcha.nrgba.Set(captcha.width-1, y, borderColor)
	}
	return captcha
}

// DrawNoise 画噪点.
func (captcha *CaptchaImage) DrawNoise(complex NoiseDensity, noiseDrawer NoiseDrawer) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	captcha.Error = noiseDrawer.DrawNoise(captcha.nrgba, complex)
	return captcha
}

// DrawText 写字.
func (captcha *CaptchaImage) DrawText(textDrawer TextDrawer, text string) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	captcha.Error = textDrawer.DrawString(captcha.nrgba, text)
	return captcha
}

// DrawBlur 对图片进行模糊处理
func (captcha *CaptchaImage) DrawBlur(drawer BlurDrawer, kernelSize int, sigma float64) *CaptchaImage {
	if captcha.Error != nil {
		return captcha
	}
	captcha.Error = drawer.DrawBlur(captcha.nrgba, kernelSize, sigma)
	return captcha
}

// CaptchaDifficulty 验证码难度级别
type CaptchaDifficulty int

const (
	// CaptchaVeryEasy 非常简单难度 - 清晰文字，无噪点，无扭曲
	CaptchaVeryEasy CaptchaDifficulty = iota
	// CaptchaEasy 简单难度 - 轻微扭曲，少量噪点
	CaptchaEasy
	// CaptchaMedium 中等难度 - 原来的简单模式
	CaptchaMedium
	// CaptchaHard 困难难度 - 原来的困难模式
	CaptchaHard
)

// GenerateCaptcha 生成验证码图片和对应的文本
func GenerateCaptcha(width, height int, textLength int, difficulty CaptchaDifficulty) (text string, imgBytes []byte, err error) {
	// 生成随机文本
	text = RandText(textLength)

	var bgColor color.RGBA
	var textColor color.RGBA
	if difficulty == CaptchaVeryEasy {
		// 随机选择一组高对比度的颜色组合
		switch rand.Intn(10) {
		case 0:
			// 白底深蓝
			bgColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
			textColor = color.RGBA{R: 0, G: 0, B: 128, A: 255}
		case 1:
			// 白底深绿
			bgColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
			textColor = color.RGBA{R: 0, G: 100, B: 0, A: 255}
		case 2:
			// 浅黄底黑
			bgColor = color.RGBA{R: 255, G: 255, B: 200, A: 255}
			textColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}
		case 3:
			// 浅蓝底深红
			bgColor = color.RGBA{R: 220, G: 240, B: 255, A: 255}
			textColor = color.RGBA{R: 180, G: 0, B: 0, A: 255}
		case 4:
			// 白底深紫
			bgColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
			textColor = color.RGBA{R: 76, G: 0, B: 153, A: 255}
		case 5:
			// 浅绿底深蓝
			bgColor = color.RGBA{R: 220, G: 255, B: 220, A: 255}
			textColor = color.RGBA{R: 0, G: 0, B: 153, A: 255}
		case 6:
			// 浅灰底深红
			bgColor = color.RGBA{R: 245, G: 245, B: 245, A: 255}
			textColor = color.RGBA{R: 153, G: 0, B: 0, A: 255}
		case 7:
			// 米色底深棕
			bgColor = color.RGBA{R: 255, G: 248, B: 220, A: 255}
			textColor = color.RGBA{R: 139, G: 69, B: 19, A: 255}
		case 8:
			// 淡青底深紫红
			bgColor = color.RGBA{R: 225, G: 255, B: 255, A: 255}
			textColor = color.RGBA{R: 139, G: 0, B: 139, A: 255}
		case 9:
			// 浅粉底深蓝绿
			bgColor = color.RGBA{R: 255, G: 240, B: 245, A: 255}
			textColor = color.RGBA{R: 0, G: 102, B: 102, A: 255}
		}
	} else {
		bgColor = RandLightColor()
		textColor = RandDeepColor()
	}

	// 创建验证码图片
	captchaImage := New(width, height, bgColor)

	// 根据难度选择不同的绘制参数
	switch difficulty {
	case CaptchaVeryEasy:
		err = captchaImage.
			DrawBorder(textColor).
			// 无扭曲的文字
			DrawText(NewTwistTextDrawer(DefaultDPI, 0, 0), text).
			Error

	case CaptchaEasy:
		err = captchaImage.
			DrawBorder(RandDeepColor()).
			// 极轻微的扭曲
			DrawText(NewTwistTextDrawer(DefaultDPI, DefaultAmplitude/4, DefaultFrequency/4), text).
			// 极少量噪点
			DrawNoise(NoiseDensityLower/2, NewPointNoiseDrawer()).
			Error

	case CaptchaMedium:
		err = captchaImage.
			DrawBorder(RandDeepColor()).
			// 只使用较低密度的点状噪点
			DrawNoise(NoiseDensityLower, NewPointNoiseDrawer()).
			// 使用更温和的文字扭曲参数
			DrawText(NewTwistTextDrawer(DefaultDPI, DefaultAmplitude/2, DefaultFrequency/2), text).
			// 只保留一条干扰线
			DrawLine(NewBeeline(), RandDeepColor()).
			// 减轻模糊效果
			DrawBlur(NewGaussianBlur(), 1, 0.3).
			Error

	default: // CaptchaHard
		err = captchaImage.
			DrawBorder(RandDeepColor()).
			DrawNoise(NoiseDensityHigh, NewTextNoiseDrawer(72)).
			DrawNoise(NoiseDensityLower, NewPointNoiseDrawer()).
			DrawLine(NewBezier3DLine(), RandDeepColor()).
			DrawText(NewTwistTextDrawer(DefaultDPI, DefaultAmplitude, DefaultFrequency), text).
			DrawLine(NewBeeline(), RandDeepColor()).
			DrawBlur(NewGaussianBlur(), DefaultBlurKernelSize, DefaultBlurSigma).
			Error
	}

	if err != nil {
		return "", nil, err
	}

	// 将图片编码为字节数组
	buf := new(bytes.Buffer)
	err = captchaImage.Encode(buf, ImageFormatJpeg)
	if err != nil {
		return "", nil, err
	}

	return text, buf.Bytes(), nil
}
