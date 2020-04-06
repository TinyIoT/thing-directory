// Copyright 2014-2016 Fraunhofer Institute for Applied Information Technology FIT

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/kelseyhightower/envconfig"
	"github.com/linksmart/go-sec/auth/obtainer"
	"github.com/linksmart/go-sec/auth/validator"
)

const (
	DNSSDServiceType = "_linksmart-td._tcp"
	BackendMemory    = "memory"
	BackendLevelDB   = "leveldb"
)

type Config struct {
	ServiceID      string         `json:"serviceID"`
	Description    string         `json:"description"`
	PublicEndpoint string         `json:"publicEndpoint"`
	BindAddr       string         `json:"bindAddr"`
	BindPort       int            `json:"bindPort"`
	DnssdEnabled   bool           `json:"dnssdEnabled"`
	Storage        StorageConfig  `json:"storage"`
	ServiceCatalog ServiceCatalog `json:"serviceCatalog"`
	Auth           validator.Conf `json:"auth"`
}

type ServiceCatalog struct {
	Enabled  bool          `json:"enabled"`
	Discover bool          `json:"discover"`
	Endpoint string        `json:"endpoint"`
	Ttl      int           `json:"ttl"`
	Auth     obtainer.Conf `json:"auth"`
}

type StorageConfig struct {
	Type string `json:"type"`
	DSN  string `json:"dsn"`
}

var supportedBackends = map[string]bool{
	BackendMemory:  false,
	BackendLevelDB: true,
}

func (c *Config) Validate() error {
	if c.BindAddr == "" || c.BindPort == 0 || c.PublicEndpoint == "" {
		return fmt.Errorf("BindAddr, BindPort, and PublicEndpoint have to be defined")
	}
	_, err := url.Parse(c.PublicEndpoint)
	if err != nil {
		return fmt.Errorf("PublicEndpoint should be a valid URL")
	}
	_, err = url.Parse(c.Storage.DSN)
	if err != nil {
		return fmt.Errorf("storage DSN should be a valid URL")
	}
	if !supportedBackends[c.Storage.Type] {
		return fmt.Errorf("unsupported storage backend")
	}

	if c.ServiceCatalog.Enabled {
		if c.ServiceCatalog.Endpoint == "" && c.ServiceCatalog.Discover == false {
			return fmt.Errorf("Service Catalog must have either endpoint or set discovery flag")
		}
		if c.ServiceCatalog.Ttl <= 0 {
			return fmt.Errorf("Service Catalog must have TTL >= 0")
		}
		if c.ServiceCatalog.Auth.Enabled {
			// Validate ticket obtainer config
			err = c.ServiceCatalog.Auth.Validate()
			if err != nil {
				return fmt.Errorf("invalid Service Catalog auth: %s", err)
			}
		}
	}

	if c.Auth.Enabled {
		// Validate ticket validator config
		err = c.Auth.Validate()
		if err != nil {
			return fmt.Errorf("invalid auth: %s", err)
		}
	}

	return err
}

func loadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	// Override loaded values with environment variables
	err = envconfig.Process("td", &config)
	if err != nil {
		return nil, err
	}

	if err = config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %s", err)
	}
	return &config, nil
}
