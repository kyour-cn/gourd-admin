package captcha

import (
	"app/internal/util/cache"
	"encoding/json"
	"errors"
	"github.com/wenlng/go-captcha-assets/helper"
	"github.com/wenlng/go-captcha-assets/resources/images"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/slide"
	"log/slog"
	"time"
)

var (
	slideCapt slide.Captcha
	initEd    = false
)

func Init() {
	if initEd {
		return
	}
	initEd = true

	builder := slide.NewBuilder(
		//slide.WithGenGraphNumber(2),
		slide.WithEnableGraphVerticalRandom(true),
	)

	// background images
	bgImages, err := images.GetImages()
	if err != nil {
		slog.Error(err.Error())
	}

	graphs, err := tiles.GetTiles()
	if err != nil {
		slog.Error(err.Error())
	}

	var newGraphs = make([]*slide.GraphImage, 0, len(graphs))
	for i := 0; i < len(graphs); i++ {
		graph := graphs[i]
		newGraphs = append(newGraphs, &slide.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}

	// set resources
	builder.SetResources(
		slide.WithGraphImages(newGraphs),
		slide.WithBackgrounds(bgImages),
	)

	slideCapt = builder.Make()
}

// GenerateSlide 生成滑动验证码
func GenerateSlide() (any, error) {
	Init()

	captData, err := slideCapt.Generate()
	if err != nil {
		return nil, err
	}

	blockData := captData.GetData()
	if blockData == nil {
		return nil, errors.New("captcha data is nil")
	}

	var masterImageBase64, tileImageBase64 string
	masterImageBase64, _ = captData.GetMasterImage().ToBase64()
	tileImageBase64, _ = captData.GetTileImage().ToBase64()

	dotsByte, _ := json.Marshal(blockData)
	key := "captcha:" + helper.StringToMD5(string(dotsByte))

	// 缓存
	cache.Set(key, dotsByte, time.Second*300)

	bt := map[string]any{
		"code":         0,
		"captcha_key":  key,
		"image_base64": masterImageBase64,
		"tile_base64":  tileImageBase64,
		"tile_width":   blockData.Width,
		"tile_height":  blockData.Height,
		"tile_x":       blockData.DX,
		"tile_y":       blockData.DY,
	}

	return bt, nil
}

// VerifySlide 验证滑动验证码
func VerifySlide(captchaKey string, x int, y int) bool {
	captcha, ok := cache.Get(captchaKey)
	if !ok {
		return false
	}

	// 删除缓存
	cache.Delete(captchaKey)

	var dct *slide.Block
	if err := json.Unmarshal(captcha.([]byte), &dct); err != nil {
		return false
	}

	return slide.Validate(x, y, dct.X, dct.Y, 4)
}
