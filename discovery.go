package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/grandcat/zeroconf"
	"github.com/tinyiot/thing-directory/wot"
)

// escape special characters as recommended by https://tools.ietf.org/html/rfc6763#section-4.3
func escapeDNSSDServiceInstance(instance string) (escaped string) {
	// replace \ by \\
	escaped = strings.ReplaceAll(instance, "\\", "\\\\")
	// replace . by \.
	escaped = strings.ReplaceAll(escaped, ".", "\\.")
	return escaped
}

// register as a DNS-SD Service
func registerDNSSDService(conf *Config) (func(), error) {
	instance := escapeDNSSDServiceInstance(conf.DNSSD.Publish.Instance)

	log.Printf("DNS-SD: registering as \"%s.%s.%s\", subtype: %s",
		instance, wot.DNSSDServiceType, conf.DNSSD.Publish.Domain, wot.DNSSDServiceSubtypeDirectory)

	var ifs []net.Interface

	for _, name := range conf.DNSSD.Publish.Interfaces {
		iface, err := net.InterfaceByName(name)
		if err != nil {
			return nil, fmt.Errorf("error finding interface %s: %s", name, err)
		}
		if (iface.Flags & net.FlagMulticast) > 0 {
			ifs = append(ifs, *iface)
		} else {
			return nil, fmt.Errorf("interface %s does not support multicast", name)
		}
		log.Printf("DNS-SD: will register to interface: %s", name)
	}

	if len(ifs) == 0 {
		log.Println("DNS-SD: publish interfaces not set. Will register to all interfaces with multicast support.")
	}

	sd, err := zeroconf.Register(
		instance,
		wot.DNSSDServiceType+","+wot.DNSSDServiceSubtypeDirectory,
		conf.DNSSD.Publish.Domain,
		conf.HTTP.BindPort,
		[]string{"td=/td", "version=" + Version},
		ifs,
	)
	if err != nil {
		return sd.Shutdown, err
	}

	return sd.Shutdown, nil
}
