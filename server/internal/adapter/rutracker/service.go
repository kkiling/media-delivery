package rutracker

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/kkiling/goplatform/log"
)

const (
	cookeFile = "rutracker_cookies.gob"
	apiUrl    = "https://rutracker.org/forum/"
)

type Api struct {
	username   string
	password   string
	cookiesDir string
	baseAPIUrl *url.URL
	httpClient *http.Client
	logger     log.Logger
}

type Config struct {
	Username   string
	Password   string
	CookiesDir string
	ProxyURL   *string
}

func NewApi(logger log.Logger, cfg Config) (*Api, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookiejar.New: %w", err)
	}
	baseAPIUrl, err := url.Parse(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}

	// Создаем базовый HTTP клиент
	httpClient := &http.Client{Jar: jar}
	logThis := logger.Named("rutracker")
	// Если передан прокси, настраиваем его
	if cfg.ProxyURL != nil {
		proxy, err := url.Parse(*cfg.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		httpClient.Transport = transport
		logThis.Infof("proxy configured")
	}

	return &Api{
		username:   cfg.Username,
		password:   cfg.Password,
		cookiesDir: cfg.CookiesDir,
		baseAPIUrl: baseAPIUrl,
		httpClient: &http.Client{Jar: jar},
		logger:     logThis,
	}, nil
}
