package binary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// Marshal serializes the structure into a byte array
func Marshal(ptr any) ([]byte, error) {
	var buf bytes.Buffer
	if err := serializeStruct(&buf, ptr); err != nil {
		return nil, fmt.Errorf("error occurred while serializing struct: %v", err)
	}
	return buf.Bytes(), nil
}

// writePrimitive writes the primitive data type to buffer
func writePrimitive(buf *bytes.Buffer, v reflect.Value) error {
	return binary.Write(buf, binary.LittleEndian, v.Interface())
}

// writeString writes the string to buffer
func writeString(buf *bytes.Buffer, str string) error {
	if err := binary.Write(buf, binary.LittleEndian, int32(len(str))); err != nil {
		return err
	}
	return binary.Write(buf, binary.LittleEndian, []byte(str))
}

// serializeStruct recursively serializes the structure
func serializeStruct(buf *bytes.Buffer, ptr any) error {
	ref := reflect.ValueOf(ptr)
	switch ref.Kind() {
	case reflect.Ptr:
		return serializeStruct(buf, ref.Elem().Interface())
	case reflect.Struct:
		for i := 0; i < ref.NumField(); i++ {
			if err := serializeStruct(buf, ref.Field(i).Interface()); err != nil {
				return err
			}
		}
	case reflect.Int32, reflect.Float64, reflect.Bool:
		if err := writePrimitive(buf, ref); err != nil {
			return err
		}
	case reflect.String:
		if err := writeString(buf, ref.String()); err != nil {
			return err
		}
	}

	return nil
}

// Unmarshal deserializes the data into the structure
func Unmarshal(data []byte, ptr interface{}) error {
	if reflect.ValueOf(ptr).Kind() != reflect.Ptr || reflect.ValueOf(ptr).IsNil() {
		return fmt.Errorf("expect a non-nil pointer, but got: %T", ptr)
	}

	buf := bytes.NewReader(data)
	if err := deserializeStruct(buf, reflect.ValueOf(ptr).Elem()); err != nil {
		return fmt.Errorf("deserializing struct: %v", err)
	}

	return nil
}

// deserializeStruct recursively deserializes the structure
func deserializeStruct(buf *bytes.Reader, ref reflect.Value) error {
	switch ref.Kind() {
	case reflect.Ptr:
		return deserializeStruct(buf, ref.Elem())
	case reflect.Struct:
		for i := 0; i < ref.NumField(); i++ {
			if err := deserializeStruct(buf, ref.Field(i)); err != nil {
				return err
			}
		}
	case reflect.Int32, reflect.Float64, reflect.Bool:
		return binary.Read(buf, binary.LittleEndian, ref.Addr().Interface())
	case reflect.String:
		var strLen int32
		if err := binary.Read(buf, binary.LittleEndian, &strLen); err != nil {
			return err
		}
		strBytes := make([]byte, strLen)
		if err := binary.Read(buf, binary.LittleEndian, &strBytes); err != nil {
			return err
		}
		ref.SetString(string(strBytes))
	}

	return nil
}
