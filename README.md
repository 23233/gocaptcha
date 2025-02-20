# gocaptcha
一个简单的Go语言实现的验证码

### 图片实例

![image](https://raw.githubusercontent.com/lifei6671/gocaptcha/master/example/image_1.jpg)
![image](https://raw.githubusercontent.com/lifei6671/gocaptcha/master/example/image_2.jpg)
![image](https://raw.githubusercontent.com/lifei6671/gocaptcha/master/example/image_3.jpg)
![image](https://raw.githubusercontent.com/lifei6671/gocaptcha/master/example/image_4.jpg)

## 简介

基于Golang实现的图片验证码生成库，可以实现随机字母个数，随机直线，随机噪点等。可以设置任意多字体，每个验证码随机选一种字体展示。

## 实例

#### 使用：

```
	go get github.com/23233/gocaptcha
```

#### 使用的类库

内置了字体与难度调节

```
	go get github.com/golang/freetype
	go get github.com/golang/freetype/truetype
	go get golang.org/x/image
```

#### 代码
具体实例可以查看example目录，有生成的验证码图片。






