package captcha

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wenlng/go-captcha-assets/helper"
	"github.com/wenlng/go-captcha-assets/resources/images"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/slide"
	"gourd/internal/util/redisutil"
	"log"
	"time"
)

var slideCapt slide.Captcha

func init() {
	builder := slide.NewBuilder(
		//slide.WithGenGraphNumber(2),
		slide.WithEnableGraphVerticalRandom(true),
	)

	// background images
	imgs, err := images.GetImages()
	if err != nil {
		log.Fatalln(err)
	}

	graphs, err := tiles.GetTiles()
	if err != nil {
		log.Fatalln(err)
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
		slide.WithBackgrounds(imgs),
	)

	slideCapt = builder.Make()
}

// GenerateSlide 生成滑动验证码
func GenerateSlide() (any, error) {
	captData, err := slideCapt.Generate()
	if err != nil {
		return nil, err
	}

	blockData := captData.GetData()
	if blockData == nil {
		return nil, errors.New("captcha data is nil")
	}

	var masterImageBase64, tileImageBase64 string
	masterImageBase64 = captData.GetMasterImage().ToBase64()
	tileImageBase64 = captData.GetTileImage().ToBase64()

	dotsByte, _ := json.Marshal(blockData)
	key := "captcha:" + helper.StringToMD5(string(dotsByte))

	// 缓存
	redis, err := redisutil.GetRedis(context.Background())
	if err != nil {
		return nil, err
	}

	redis.Set(context.Background(), key, dotsByte, time.Second*300)

	bt := map[string]any{
		"code":         0,
		"captcha_key":  key,
		"image_base64": masterImageBase64,
		"tile_base64":  tileImageBase64,
		"tile_width":   blockData.Width,
		"tile_height":  blockData.Height,
		"tile_x":       blockData.TileX,
		"tile_y":       blockData.TileY,
	}

	return bt, nil
}

// VerifySlide 验证滑动验证码
func VerifySlide(captchaKey string, x int64, y int64) bool {
	redis, err := redisutil.GetRedis(context.Background())
	if err != nil {
		return false
	}

	ctx := context.Background()
	captchaByte, err := redis.Get(ctx, captchaKey).Bytes()
	if err != nil {
		return false
	}

	// 删除缓存
	redis.Del(ctx, captchaKey)

	var dct *slide.Block
	if err := json.Unmarshal(captchaByte, &dct); err != nil {
		return false
	}

	return slide.CheckPoint(x, y, int64(dct.X), int64(dct.Y), 4)
}
