package main

import (
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	SMTPAddr string `yaml:"smtp_addr"`
}

func main() {
	conf, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	auth := sasl.NewLoginClient(conf.Username, conf.Password)

	to := []string{"lunamicard@gmail.com"} // todo parse

	if err := smtp.SendMail(conf.SMTPAddr, auth, conf.Username, to, os.Stdin); err != nil {
		log.Fatal(err)
	}
}

func readConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	f, err := os.Open(filepath.Join(homeDir, ".mailman.yaml"))
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var conf Config
	if err := yaml.NewDecoder(f).Decode(&conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}
