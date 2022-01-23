package util

import (
	"io/ioutil"
	"log"
)

func WriteValuesToFile(yaml string, output string) {
	err := ioutil.WriteFile(output, []byte(yaml), 0600)
	if err != nil {
		log.Fatal(err)
	}
}
