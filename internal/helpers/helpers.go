package helpers

import (
	"bytes"
	"encoding/gob"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
)

type DataWrapper struct {
	Data any
}

// CheckErrAndExit checks if an error is not nil and exits the program with exit code 1
func CheckErrAndExit(err error) {
	if err != nil {
		os.Exit(1)
	}
}

func buildMapFromStruct(v any) map[string]any {
	m := make(map[string]any)
	t := reflect.TypeOf(v)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := reflect.ValueOf(v).FieldByName(f.Name)
		m[f.Name] = fv.Interface()
	}

	return m
}

func serializeStructMap(m map[string]any) ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(m); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func Serialize(v any) ([]byte, error) {
	m := buildMapFromStruct(v)
	return serializeStructMap(m)
}

func Deserialize(data []byte, v map[string]any) error {
	buffer := bytes.Buffer{}
	buffer.Write(data)
	decoder := gob.NewDecoder(&buffer)
	if err := decoder.Decode(&v); err != nil {
		return err
	}
	return nil
}

func Closer(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Println(err)
	}
}

func Binarization(v any) []byte {
	var b []byte
	m := make(map[string]any)
	t := reflect.TypeOf(v)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := reflect.ValueOf(v).FieldByName(f.Name)
		m[f.Name] = fv.Interface()
	}

	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(m); err != nil {
		log.Println(err)
	}

	b = buffer.Bytes()
	return b
}

// ArgsToString converts args to string
func ArgsToString(args []string) string {
	var msg string
	for _, arg := range args {
		msg += arg + " "
	}

	// trim last space
	msg = strings.TrimRight(msg, " ")
	return msg
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetString viper get string and check if it is empty
func GetString(key string) (string, error) {
	val := viper.GetString(key)
	if val == "" {
		return "", nil
	}
	return val, nil
}
