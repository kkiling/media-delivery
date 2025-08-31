package themoviedb

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/kkiling/goplatform/log"
)

const (
	baseApiUrl = "https://api.themoviedb.org/3"
	baseImgUrl = "https://image.tmdb.org/t/p"
	maxPages   = 5
)

type Config struct {
	APIKey   string
	ProxyURL *string
}

type API struct {
	apiKey     string
	baseAPIUrl *url.URL
	baseImgUrl *url.URL
	httpClient *http.Client
	logger     log.Logger
	validate   *validator.Validate
}

func NewApi(logger log.Logger, cfg Config) (*API, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookiejar.New: %w", err)
	}
	urlApi, err := url.Parse(baseApiUrl)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}
	urlImg, err := url.Parse(baseImgUrl)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}

	// Создаем базовый HTTP клиент
	httpClient := &http.Client{Jar: jar}
	logThis := logger.Named("the_movie_db")
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

	return &API{
		apiKey:     cfg.APIKey,
		baseAPIUrl: urlApi,
		baseImgUrl: urlImg,
		httpClient: httpClient,
		logger:     logThis,
		validate:   validator.New(),
	}, nil
}
