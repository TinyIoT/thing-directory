package catalog

import (
	"context"
	"fmt"

	"github.com/tinyiot/thing-directory/wot"
)

type ThingDescription = map[string]interface{}

const (
	// Storage backend types
	BackendMemory  = "memory"
	BackendLevelDB = "leveldb"
)

func validateThingDescription(td map[string]interface{}) ([]wot.ValidationError, error) {
	result, err := wot.ValidateTD(&td)
	if err != nil {
		return nil, fmt.Errorf("error validating with JSON Schemas: %s", err)
	}
	return result, nil
}

// Controller interface
type CatalogController interface {
	add(d ThingDescription) (string, error)
	get(id string) (ThingDescription, error)
	update(id string, d ThingDescription) error
	patch(id string, d ThingDescription) error
	delete(id string) error
	listPaginate(offset, limit int) ([]ThingDescription, error)
	filterJSONPathBytes(query string) ([]byte, error)
	iterateBytes(ctx context.Context) <-chan []byte
	cleanExpired()
	Stop()
	AddSubscriber(listener EventListener)
}

// Storage interface
type Storage interface {
	add(id string, td ThingDescription) error
	update(id string, td ThingDescription) error
	delete(id string) error
	get(id string) (ThingDescription, error)
	listPaginate(offset, limit int) ([]ThingDescription, error)
	listAllBytes() ([]byte, error)
	iterate() <-chan ThingDescription
	iterateBytes(ctx context.Context) <-chan []byte
	Close()
}
