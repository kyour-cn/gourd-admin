package methods

import "github.com/json-iterator/go"

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
