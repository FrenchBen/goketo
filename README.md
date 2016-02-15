GoKeto: Marketo REST API Client
===============================
<p align="center">
  <img src="https://raw.githubusercontent.com/FrenchBen/go-marketo/master/doc/Marketo-logo.jpg" alt="Marketo Logo"/>
</p>

About
----------------
Unofficial Golang client for the Marketo.com REST API: http://developers.marketo.com/documentation/rest/.
Inspired by the `VojtechVitek/go-trello` implementation

Requires Go `1.5.3`

Installation
----------------
The recommended way of installing the client is via `go get`. Simply run the following command to add the package.

    go get github.com/FrenchBen/goketo/

Usage
----------------
Below is an example of how to use this library

```
package main

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/FrenchBen/goketo"
	"github.com/Sirupsen/logrus"
)


func main() {
	// Get config auth
	authFile, _ := filepath.Abs("./auth.yaml")
	data, err := ioutil.ReadFile(authFile)
	if err != nil {
		logrus.Errorf("error reading auth file %q: %v", "auth.yaml", err)
	}
	auth := &marketoAuthConfig{}
	if err = yaml.Unmarshal(data, auth); err != nil {
		logrus.Errorf("Error during Yaml: %v", err)
	}
	// New Marketo Client
	marketo, err := goketo.NewAuthClient(auth.ClientID, auth.ClientSecret, auth.ClientEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	// Get leads
	leads, err := marketo.Leads(auth.LeadID)
	if err != nil {
		logrus.Error("Couldn't get leads: ", err)
	}
  logrus.Infof("My leads: %v", leads)

  // Get user by lead ID
	lead, err := marketo.Lead(leads.Result[0].ID)
	if err != nil {
		logrus.Error("Couldn't get lead: ", err)
	}
  logrus.Infof("My lead from ID: %v", lead)
}
```


License
----------------
This source is licensed under an MIT License, see the LICENSE file for full details. If you use this code, it would be great to hear from you.
