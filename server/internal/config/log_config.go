package config

// LogConfig 日志配置
type LogConfig struct {
	Level    string `toml:"level" json:"level"`       // 日志记录级别 debug、info、warn、error
	LogFile  string `toml:"log_file" json:"log_file"` // 日志文件
	Console  bool   `toml:"console" json:"console"`   // 是否开启控制台输出
	Encoding string `toml:"encoding" json:"encoding"` // 输出格式 default、json、text
	MaxSize  int    `toml:"max_size" json:"max_size"` // 单个日志文件最大大小，单位：M
}

var logConf *LogConfig

// GetLogConfig 获取Log服务器配置
func GetLogConfig() (*LogConfig, error) {
	key := "log"

	// 初始化配置
	if logConf == nil {
		_conf := &LogConfig{
			Level:    "info",
			LogFile:  "app.log",
			Encoding: "default",
			MaxSize:  10,
		}

		// 如果配置不存在，则创建默认配置
		if !Exists(key) {
			err := SetLogConfig(_conf)
			if err != nil {
				return nil, err
			}
		}

		err := Unmarshal(key, _conf)
		if err != nil {
			return nil, err
		}
		logConf = _conf
	}

	return logConf, nil
}

// SetLogConfig 设置Log服务器配置
func SetLogConfig(conf *LogConfig) error {
	key := "log"
	logConf = conf
	return Marshal(key, conf)
}
