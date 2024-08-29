package pokeapi

import (
	"net/http"
	"time"

	"github.com/moceviciusda/pokeCLIpse-server/internal/cache"
)

const baseURL = "http://pokeapi.co/api/v2"

type Client struct {
	cache      cache.Cache
	httpClient http.Client
}

func NewClient(cacheInterval, timeout time.Duration) Client {
	return Client{
		cache.NewCache(cacheInterval),
		http.Client{
			Timeout: timeout,
		},
	}
}
