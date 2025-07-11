package config

// HttpConfig Http服务器配置
type HttpConfig struct {
	Enable bool   `toml:"enable" json:"enable"` // 是否启用Http服务器
	Host   string `toml:"host" json:"host"`     // 监听域名、IP
	Port   uint16 `toml:"port" json:"port"`     // 监听端口
	Static string `toml:"static" json:"static"` // 静态资源目录
}

var httpConf *HttpConfig

// GetHttpConfig 获取Http服务器配置
func GetHttpConfig() (*HttpConfig, error) {
	key := "http"

	// 初始化配置
	if httpConf == nil {
		_conf := &HttpConfig{
			Host: "0.0.0.0",
			Port: 8080,
		}

		// 如果配置不存在，则创建默认配置
		if !Exists(key) {
			err := SetHttpConfig(_conf)
			if err != nil {
				return nil, err
			}
		}

		err := Unmarshal(key, _conf)
		if err != nil {
			return nil, err
		}
		httpConf = _conf
	}

	return httpConf, nil
}

// SetHttpConfig 设置Http服务器配置
func SetHttpConfig(conf *HttpConfig) error {
	key := "http"
	httpConf = conf
	return Marshal(key, conf)
}
