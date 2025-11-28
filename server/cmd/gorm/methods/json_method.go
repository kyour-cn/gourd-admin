package methods

import (
	"time"
	"unsafe"

	"github.com/json-iterator/go"
	"gorm.io/gorm"
)

func init() {

	// 注册全局 encoder/empty 判定
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

}

// JsonMethod 实现json序列化、反序列化接口方法
type JsonMethod struct {
}

// MarshalBinary 支持json序列化
func (m *JsonMethod) MarshalBinary() (data []byte, err error) {
	return jsoniter.Marshal(m)
}

// UnmarshalBinary 支持json反序列化
func (m *JsonMethod) UnmarshalBinary(data []byte) error {
	return jsoniter.Unmarshal(data, m)
}
