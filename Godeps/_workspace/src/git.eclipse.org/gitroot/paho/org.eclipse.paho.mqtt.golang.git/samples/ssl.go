/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

/*
To run this sample, The following certificates
must be created:

  rootCA-crt.pem - root certificate authority that is used
                   to sign and verify the client and server
                   certificates.
  rootCA-key.pem - keyfile for the rootCA.

  server-crt.pem - server certificate signed by the CA.
  server-key.pem - keyfile for the server certificate.

  client-crt.pem - client certificate signed by the CA.
  client-key.pem - keyfile for the client certificate.

  CAfile.pem     - file containing concatenated CA certificates
                   if there is more than 1 in the chain.
                   (e.g. root CA -> intermediate CA -> server cert)

  Instead of creating CAfile.pem, rootCA-crt.pem can be added
  to the default openssl CA certificate bundle. To find the
  default CA bundle used, check:
  $GO_ROOT/src/pks/crypto/x509/root_unix.go
  To use this CA bundle, just set tls.Config.RootCAs = nil.
*/

package main

import "io/ioutil"
import "fmt"
import "time"
import "crypto/tls"
import "crypto/x509"
import MQTT "linksmart.eu/lc/core/Godeps/_workspace/src/git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"

func NewTlsConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("samplecerts/CAfile.pem")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair("samplecerts/client-crt.pem", "samplecerts/client-key.pem")
	if err != nil {
		panic(err)
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(cert.Leaf)

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}

var f MQTT.MessageHandler = func(client *MQTT.MqttClient, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	tlsconfig := NewTlsConfig()

	opts := MQTT.NewClientOptions()
	opts.AddBroker("ssl://hushbox.net:17004")
	opts.SetClientId("ssl-sample").SetTlsConfig(tlsconfig)
	opts.SetDefaultPublishHandler(f)

	// Start the connection
	c := MQTT.NewClient(opts)
	_, err := c.Start()
	if err != nil {
		panic(err)
	}

	filter, _ := MQTT.NewTopicFilter("/go-mqtt/sample", 0)
	c.StartSubscription(nil, filter)

	i := 0
	for _ = range time.Tick(time.Duration(1) * time.Second) {
		if i == 5 {
			break
		}
		text := fmt.Sprintf("this is msg #%d!", i)
		c.Publish(MQTT.QOS_ZERO, "/go-mqtt/sample", []byte(text))
		i++
	}

	c.Disconnect(250)
}
