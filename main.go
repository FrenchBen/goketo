package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/FrenchBen/go-marketo/goketo"
	"github.com/Sirupsen/logrus"
)

type marketoAuthConfig struct {
	ClientID       string `json:"client_id" yaml:"client_id"`
	ClientSecret   string `json:"client_secret" yaml:"client_secret"`
	ClientEndpoint string `json:"endpoint" yaml:"endpoint"`
}

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
	goketo, err := goketo.NewAuthClient(auth.ClientID, auth.ClientSecret, auth.ClientEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	logrus.Infof("We have client! %#v", goketo)
}
