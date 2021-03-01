package main

import (
	"bufio"
	"bytes"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/textproto"
	"os"
	"path/filepath"
)

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	SMTPAddr string `yaml:"smtp_addr"`
}

func main() {
	// read the config file
	conf, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	// read mail from stdin
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	// parse mail headers
	tp := textproto.NewReader(bufio.NewReader(bytes.NewReader(b))) // TODO remove this :(
	header, err := tp.ReadMIMEHeader()
	if err != nil {
		log.Fatal(err)
	}

	// Read `From`
	from := conf.Username
	if val := header.Get("From"); val != "" {
		val = from
	}

	// Read `To`
	var to []string
	if val := header["To"]; len(val) > 0 {
		to = val
	}

	auth := sasl.NewLoginClient(conf.Username, conf.Password)
	if err := smtp.SendMail(conf.SMTPAddr, auth, from, to, bytes.NewReader(b)); err != nil {
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
