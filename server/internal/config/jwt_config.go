package config

// JwtConfig 日志配置
type JwtConfig struct {
	Secret string `toml:"secret" json:"secret" comment:"加密密钥"`
	Expire int64  `toml:"expire" json:"expire" comment:"过期时间（单位：秒）"`
}

var jwtConf *JwtConfig

// GetJwtConfig 获取Log服务器配置
func GetJwtConfig() (*JwtConfig, error) {
	key := "jwt"

	// 初始化配置
	if jwtConf == nil {
		_conf := &JwtConfig{}

		// 如果配置不存在，则创建默认配置
		if !Exists(key) {
			err := SetJwtConfig(_conf)
			if err != nil {
				return nil, err
			}
		}

		err := Unmarshal(key, _conf)
		if err != nil {
			return nil, err
		}
		jwtConf = _conf
	}

	return jwtConf, nil
}

// SetJwtConfig 设置Log服务器配置
func SetJwtConfig(conf *JwtConfig) error {
	key := "jwt"
	jwtConf = conf
	return Marshal(key, conf)
}
