package initialize

import (
	"context"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"

	"app/internal/util/cache"
)

func InitCommon(ctx context.Context) {

	// 初始化日志
	err := InitLog()
	if err != nil {
		panic(err)
	}

	// 注册全局 encoder/empty 判定（在程序初始化时执行一次）
	// 主要用来处理 time.Time 类型的零值，后期建议使用json/v2替换
	jsoniter.RegisterTypeEncoderFunc("time.Time",
		func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			t := *((*time.Time)(ptr))
			if t.IsZero() {
				stream.WriteString("")
				return
			}
			stream.WriteString(t.Format("2006-01-02 15:04:05"))
		},
		func(ptr unsafe.Pointer) bool {
			t := *((*time.Time)(ptr))
			return t.IsZero()
		},
	)
	jsoniter.RegisterTypeEncoderFunc("gorm.DeletedAt",
		func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			t := *((*gorm.DeletedAt)(ptr))
			if !t.Valid {
				stream.WriteString("")
				return
			}
			stream.WriteString(t.Time.Format("2006-01-02 15:04:05"))
		},
		func(ptr unsafe.Pointer) bool {
			t := *((*time.Time)(ptr))
			return t.IsZero()
		},
	)

	// 初始化缓存
	cache.InitDefaultCache(ctx)
}
