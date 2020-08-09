package captcha

import (
	"github.com/mojocn/base64Captcha"
	"github.com/pkg/errors"
	"image/color"
)

type Captcha struct {
	store  base64Captcha.Store
	engine *base64Captcha.Captcha
}

type AudioCaptchaOptions struct {
	Length   int
	Language string
}

type StringCaptchaOptions struct {
	ImageHeight     int
	ImageWidth      int
	Length          int
	NoiseCount      int
	ShowLineOptions int
	Source          string
	BgColor         *color.RGBA
	fonts           []string
}

type MathCaptchaOptions struct {
	ImageHeight     int
	ImageWidth      int
	NoiseCount      int
	ShowLineOptions int
	BgColor         *color.RGBA
	fonts           []string
}

type ChineseCaptchaOptions struct {
	ImageHeight     int
	ImageWidth      int
	Length          int
	NoiseCount      int
	ShowLineOptions int
	Source          string
	BgColor         *color.RGBA
	fonts           []string
}

type DigitCaptchaOptions struct {
	ImageHeight int
	ImageWidth  int
	Length      int
	MaxSkew     float64
	DotCount    int
}

type CaptchaData struct {
	Id        string
	Base64Str string
}

// New 创建新的验证码实例
// base64Captcha.DefaultMemStore
//
func New(options interface{}) (*Captcha, error) {
	store, err := getDefaultCaptchaStore()
	if err != nil {
		return nil, err
	}

	var driver base64Captcha.Driver
	if op, ok := options.(AudioCaptchaOptions); ok {
		driver = base64Captcha.NewDriverAudio(op.Length, op.Language)
	} else if op, ok := options.(StringCaptchaOptions); ok {
		driver = base64Captcha.NewDriverString(
			op.ImageHeight,
			op.ImageWidth,
			op.NoiseCount,
			op.ShowLineOptions,
			op.Length,
			op.Source,
			op.BgColor,
			op.fonts,
		).ConvertFonts()
	} else if op, ok := options.(MathCaptchaOptions); ok {
		driver = base64Captcha.NewDriverMath(
			op.ImageHeight,
			op.ImageWidth,
			op.NoiseCount,
			op.ShowLineOptions,
			op.BgColor,
			op.fonts,
		).ConvertFonts()
	} else if op, ok := options.(ChineseCaptchaOptions); ok {
		driver = base64Captcha.NewDriverChinese(
			op.ImageHeight,
			op.ImageWidth,
			op.NoiseCount,
			op.ShowLineOptions,
			op.Length,
			op.Source,
			op.BgColor,
			op.fonts,
		).ConvertFonts()
	} else if op, ok := options.(DigitCaptchaOptions); ok {
		driver = base64Captcha.NewDriverDigit(
			op.ImageHeight,
			op.ImageWidth,
			op.Length,
			op.MaxSkew,
			op.DotCount,
		)
	} else {
		return nil, errors.New("options参数非法")
	}

	// 生成driver
	engine := base64Captcha.NewCaptcha(driver, store)
	return &Captcha{
		store:  store,
		engine: engine,
	}, nil
}

// Generate 生成验证码
func (c *Captcha) Generate() (*CaptchaData, error) {
	id, b64s, err := c.engine.Generate()
	if err != nil {
		return nil, errors.WithMessage(err, "生成验证码失败")
	}
	return &CaptchaData{
		Id:        id,
		Base64Str: b64s,
	}, nil
}

// Verify 校验验证码
func (c *Captcha) Verify(id string, captcha string, isClear bool) bool {
	return c.store.Verify(id, captcha, isClear)
}
