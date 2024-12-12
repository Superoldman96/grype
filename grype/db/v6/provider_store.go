package v6

import (
	"fmt"
	"sort"

	"gorm.io/gorm"

	"github.com/anchore/grype/internal/log"
)

type ProviderStoreReader interface {
	GetProvider(name string) (*Provider, error)
	AllProviders() ([]Provider, error)
}

type providerStore struct {
	db *gorm.DB
}

func newProviderStore(db *gorm.DB) *providerStore {
	return &providerStore{
		db: db,
	}
}

func (s *providerStore) GetProvider(name string) (*Provider, error) {
	log.WithFields("name", name).Trace("fetching provider record")

	var provider Provider
	result := s.db.Where("id = ?", name).First(&provider)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch provider (name=%q): %w", name, result.Error)
	}

	return &provider, nil
}

func (s *providerStore) AllProviders() ([]Provider, error) {
	log.Trace("fetching all provider records")

	var providers []Provider
	result := s.db.Find(&providers)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch all providers: %w", result.Error)
	}

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].ID < providers[j].ID
	})

	return providers, nil
}