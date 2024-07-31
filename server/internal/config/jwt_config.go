package config

// JwtConfig 日志配置
type JwtConfig struct {
	Secret string `toml:"secret" json:"secret"` // 加密密钥
	Expire int64  `toml:"expire" json:"expire"` // 过期时间（单位：秒）
}

var jwtConf *JwtConfig

// GetJwtConfig 获取Log服务器配置
func GetJwtConfig() (*JwtConfig, error) {

	// 初始化配置
	if jwtConf == nil {
		_conf := &JwtConfig{
			Secret: "gourd-admin",
		}
		err := Unmarshal("jwt", _conf)
		if err != nil {
			return nil, err
		}
		jwtConf = _conf
	}

	return jwtConf, nil
}

// SetJwtConfig 设置Log服务器配置
func SetJwtConfig(conf *JwtConfig) {
	jwtConf = conf
}
