GoKeto: Marketo REST API Client
===============================
[![GoDoc](https://godoc.org/github.com/FrenchBen/goketo?status.svg)](https://godoc.org/github.com/FrenchBen/goketo)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/FrenchBen/goketo/master/LICENSE)
[![CircleCI](https://circleci.com/gh/FrenchBen/goketo.svg?style=shield)](https://circleci.com/gh/FrenchBen/goketo)
[![TravisCI](https://travis-ci.org/FrenchBen/goketo.svg?branch=master)](https://travis-ci.org/FrenchBen/goketo)
[![Go Report Card](https://goreportcard.com/badge/github.com/FrenchBen/goketo)](https://goreportcard.com/report/github.com/FrenchBen/goketo)
[![Badge Badge](http://doyouevenbadge.com/github.com/FrenchBen/goketo)](http://doyouevenbadge.com)

<p align="center">
  <a href="http://golang.org" target="_blank"><img alt="Go package" src="https://golang.org/doc/gopher/pencil/gopherhat.jpg" width="20%" /></a>
  <a href="https://www.marketo.com/" target="_blank"><img src="https://raw.githubusercontent.com/FrenchBen/go-marketo/master/doc/Marketo-logo.jpg" alt="Marketo Logo"/></a>
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
  listID, _ := strconv.Atoi(auth.LeadID)
	leadRequest := &goketo.LeadRequest{ID: listID}
	leads, err := goketo.Leads(marketo, leadRequest)
	if err != nil {
		logrus.Error("Couldn't get leads: ", err)
	}  
  results := []goketo.LeadResult{}
	err = json.Unmarshal(leads.Result, &results)
  logrus.Infof("My leads: %v", results)


  // Get user by lead ID
  leadID, _ := results[0].ID
	leadRequest := &goketo.LeadRequest{ID: leadID}
	lead, err := goketo.Lead(marketo, leadRequest)
	if err != nil {
		logrus.Error("Couldn't get lead: ", err)
	}
  result := []goketo.LeadResult{}
	err = json.Unmarshal(leads.Result, &result)
  logrus.Infof("My lead from ID: %v", result)
}
```

To view more the token and fields sent with the request, set your log level to debug:
`logrus.SetLevel(logrus.DebugLevel)`

For information on usage, please see the [GoDoc](https://godoc.org/github.com/FrenchBen/goketo).

License
----------------
This source is licensed under an MIT License, see the LICENSE file for full details. If you use this code, it would be great to hear from you.
