package internal

import (
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// NewCache returns a new cache for Moirai
func NewCache() (cache.Cache, error) {
	c, err := cache.New(config.GetConfigOrDie(), cache.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}
