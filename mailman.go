package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("/var/log/mailman.log", b, 0640); err != nil {
		log.Fatal(err)
	}
}
