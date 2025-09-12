package initialize

import (
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

func InitCommon() {

	// 注册全局 encoder/empty 判定（在程序初始化时执行一次）
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
}
