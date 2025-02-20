package gocaptcha

import (
	"testing"
)

func TestFontFamily_Random(t *testing.T) {
	type args struct {
	}
	fontFamily := NewFontFamily()
	tests := []struct {
		name     string
		t        *FontFamily
		args     args
		fn       func(*FontFamily) error
		wantErr  bool
		wantFont bool
	}{
		{
			name: "test1",
			t:    fontFamily,
			args: args{},
			fn: func(family *FontFamily) error {
				family.fonts = nil // 清空字体列表测试
				return nil
			},
			wantErr:  true,
			wantFont: false,
		},
		{
			name: "test2",
			t:    fontFamily,
			args: args{},
			fn: func(family *FontFamily) error {
				// 使用内嵌字体路径
				family.fonts = []string{"fonts/3Dumb.ttf"}
				return nil
			},
			wantErr:  false,
			wantFont: true,
		},
		{
			name: "test3",
			t:    fontFamily,
			args: args{},
			fn: func(family *FontFamily) error {
				// 使用内嵌字体路径
				family.fonts = []string{"fonts/3Dumb.ttf"}
				return nil
			},
			wantErr:  false,
			wantFont: true,
		},
		{
			name: "test4",
			t:    fontFamily,
			args: args{},
			fn: func(family *FontFamily) error {
				// 使用内嵌字体路径
				return family.AddFont("fonts/3Dumb.ttf")
			},
			wantErr:  false,
			wantFont: true,
		},
		{
			name: "test5",
			t:    fontFamily,
			args: args{},
			fn: func(family *FontFamily) error {
				// 使用内嵌字体路径
				return family.AddFont("fonts/Comismsh.ttf")
			},
			wantErr:  false,
			wantFont: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.fn(tt.t)

			if font, err := tt.t.Random(); (font == nil) == tt.wantFont || (err != nil) != tt.wantErr {
				t.Errorf("FontFamily.Random() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetFonts(t *testing.T) {
	type args struct {
		fonts []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				// 使用内嵌字体路径
				fonts: []string{"fonts/3Dumb.ttf"},
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				// 使用不存在的字体测试错误情况
				fonts: []string{"fonts/nonexistent.ttf"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetFonts(tt.args.fonts...); (err != nil) != tt.wantErr {
				t.Errorf("SetFonts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFontFamily_AddFontPath(t *testing.T) {
	type args struct {
		dirPath string
	}
	fontFamily := NewFontFamily()
	tests := []struct {
		name    string
		t       *FontFamily
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			t:       fontFamily,
			args:    args{dirPath: "fonts"},
			wantErr: false,
		},
		{
			name:    "test2",
			t:       fontFamily,
			args:    args{dirPath: "nonexistent"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.AddFontPath(tt.args.dirPath); (err != nil) != tt.wantErr {
				t.Errorf("FontFamily.AddFontPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
