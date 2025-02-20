package gocaptcha

import (
	"os"
	"testing"
)

func TestCaptchaImage_Encode(t *testing.T) {
	err := SetFontPath("./fonts")
	if err != nil {
		t.Fatal(err)
	}
	captchaImage := New(150, 20, RandLightColor())
	err = captchaImage.
		DrawBorder(RandDeepColor()).
		DrawNoise(NoiseDensityHigh, NewTextNoiseDrawer(72)).
		DrawNoise(NoiseDensityLower, NewPointNoiseDrawer()).
		DrawLine(NewBezier3DLine(), RandDeepColor()).
		DrawText(NewTwistTextDrawer(DefaultDPI, DefaultAmplitude, DefaultFrequency), RandText(4)).
		DrawLine(NewBeeline(), RandDeepColor()).
		DrawLine(NewHollowLine(), RandLightColor()).
		DrawBlur(NewGaussianBlur(), DefaultBlurKernelSize, DefaultBlurSigma).
		Error
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateCaptcha(t *testing.T) {
	gotText, gotImgBytes, err := GenerateCaptcha(180, 60, 4, CaptchaVeryEasy)
	if err != nil {
		t.Fatal(err)
	}
	if len(gotText) != 4 {
		t.Errorf("GenerateCaptcha() gotText = %v, want %v", gotText, 4)
	}

	// 将生成的图片保存到文件
	err = os.WriteFile("test_captcha.png", gotImgBytes, 0644)
	if err != nil {
		t.Fatal("Failed to save captcha image:", err)
	}
}
