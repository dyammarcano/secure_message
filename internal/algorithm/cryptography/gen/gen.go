//go:build generate

package main

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/caarlos0/log"
	"os"
	"text/template"
)

//go:generate go run gen.go

const (
	keysLength = 44
	keysNumber = 512
	fileName   = "generated_file.go"
	tmpl       = `// Package crypto this package was created by with 'go generate ./...' command

package cryptography

const (
	keysLenght = {{.KeysLenght}}
)

var stringKeys = map[int]string{
	{{- range .Keys}}
	{{.Key}}: "{{.Value}}",{{end}}
}`
)

// Generate a master key
func generateKeys(size int) string {
	masterKey := make([]byte, size)
	_, _ = rand.Read(masterKey)
	return hex.EncodeToString(masterKey)
}

type keysStruct struct {
	Keys       []KeyValue
	KeysLenght int
}

type KeyValue struct {
	Key   int
	Value string
}

func main() {
	newTmpl, err := template.New("keys").Parse(tmpl)
	if err != nil {
		log.Fatalf("error parsing template: %s", err.Error())
	}

	keys := keysStruct{
		Keys:       []KeyValue{},
		KeysLenght: keysLength,
	}

	for i := 0; i < keysNumber; i++ {
		keys.Keys = append(keys.Keys, KeyValue{Key: i, Value: generateKeys(keysLength)})
	}

	//change directory
	if err = os.Chdir(".."); err != nil {
		log.Fatalf("error changing directory: %s", err.Error())
	}

	outputFile, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("error creating file: %s", err.Error())
	}
	defer func(outputFile *os.File) {
		if err = outputFile.Close(); err != nil {
			log.Fatalf("error closing file: %s", err.Error())
		}
	}(outputFile)

	if err = newTmpl.Execute(outputFile, keys); err != nil {
		log.Fatalf("error executing template: %s", err.Error())
	}

	log.Infof("generated file: %s", fileName)
}
