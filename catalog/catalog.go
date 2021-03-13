// Copyright 2014-2016 Fraunhofer Institute for Applied Information Technology FIT

package catalog

import (
	"fmt"

	"github.com/linksmart/thing-directory/wot"
)

type ThingDescription = map[string]interface{}

const (
	ResponseContextURL = "https://linksmart.eu/thing-directory/context.jsonld"
	ResponseType       = "Catalog"
	ResponseMediaType  = "application/ld+json"
	// DNS-SD
	DNSSDServiceType    = "_wot._tcp"
	DNSSDServiceSubtype = "_directory" // _directory._sub._wot._tcp
	// Storage backend types
	BackendMemory  = "memory"
	BackendLevelDB = "leveldb"
	// TD keys used internally
	_id       = "id"
	_created  = "created"
	_modified = "modified"
	_ttl      = "ttl"
)

func validateThingDescription(td map[string]interface{}) ([]wot.ValidationError, error) {
	issues, err := wot.ValidateMap(&td)
	if err != nil {
		return nil, fmt.Errorf("error validating with JSON schema: %s", err)
	}

	if td[_ttl] != nil {
		_, ok := td[_ttl].(float64)
		if !ok {
			issues = append(issues, wot.ValidationError{
				Name:   _ttl,
				Reason: fmt.Sprintf("Invalid type. Expected float64, given: %T", td[_ttl]),
			})
		}
	}

	return issues, nil
}

// Controller interface
type CatalogController interface {
	add(d ThingDescription) (string, error)
	get(id string) (ThingDescription, error)
	update(id string, d ThingDescription) error
	patch(id string, d ThingDescription) error
	delete(id string) error
	list(page, perPage int) ([]ThingDescription, int, error)
	filterJSONPath(path string, page, perPage int) ([]interface{}, int, error)
	filterXPath(path string, page, perPage int) ([]interface{}, int, error)
	total() (int, error)
	cleanExpired()

	Stop()
}

// Storage interface
type Storage interface {
	add(id string, td ThingDescription) error
	update(id string, td ThingDescription) error
	delete(id string) error
	get(id string) (ThingDescription, error)
	list(page, perPage int) ([]ThingDescription, int, error)
	total() (int, error)
	iterator() <-chan ThingDescription
	Close()
}
