//go:build generate && linux

package main

import (
	"encoding/hex"
	"fmt"
	"github.com/dyammarcano/secure_message/internal/algorithm/cryptography"
	"os"
	"strings"
)

//go:generate go run generate_keys.go

func main() {
	fileName := "generated_file.go"

	var builder strings.Builder
	builder.WriteString(`// Package crypto this package was created by code generator

package cryptography

import "encoding/hex"

var stringKeys = map[int]string{
`)

	for i := 0; i < 256; i++ {
		data, err := cryptography.GenerateKeys(44)
		if err != nil {
			panic(err)
		}
		builder.WriteString(fmt.Sprintf("\t%d: \"%s\",\n", i, hex.EncodeToString(data)))
	}

	builder.WriteString(`}

func Keys() map[int][]byte {
	keys := make(map[int][]byte)
	for i, key := range stringKeys {
		data, err := hex.DecodeString(key)
		if err != nil {
			panic(err)
		}
		keys[i] = data
	}
	return keys
}
`)

	//change directory
	if err := os.Chdir(".."); err != nil {
		panic(err)
	}

	outputFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	if _, err := outputFile.WriteString(builder.String()); err != nil {
		panic(err)
	}

	fmt.Printf("Generated file: %s", fileName)
}
