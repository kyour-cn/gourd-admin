package services

import (
	"app/internal/orm/model"
	"app/internal/orm/query"
	"app/internal/util/cache"
	"context"
	"time"

	json "github.com/json-iterator/go"
)

func NewConfigService(ctx context.Context) *ConfigService {
	return &ConfigService{
		ctx: ctx,
	}
}

type ConfigService struct {
	ctx context.Context
}

func (s *ConfigService) Get(key string) (*model.Config, error) {
	return query.Config.Where(query.Config.Key.Eq(key)).First()
}

type SiteConfig struct {
	AdminCaptchaSwitch bool `json:"admin_captcha_switch"`
}

func (s *ConfigService) GetSite() (*SiteConfig, error) {
	// 从缓存中获取站点配置(5秒缓存)
	return cache.Remember("config_cache_site", time.Second*5, func() (*SiteConfig, error) {
		conf, err := s.Get("site")
		if err != nil {
			return nil, err
		}
		sc := SiteConfig{}
		if conf.Value != nil && *conf.Value != "" {
			err = json.Unmarshal([]byte(*conf.Value), &sc)
			if err != nil {
				return nil, err
			}
		}
		return &sc, nil
	})
}
