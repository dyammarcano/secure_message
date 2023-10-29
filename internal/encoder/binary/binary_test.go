package binary

//
//import (
//	"fmt"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//type (
//	MyStruct struct {
//		Field1      int32
//		Field2      string
//		OtherStruct OtherStruct
//	}
//
//	OtherStruct struct {
//		FieldA float64
//		FieldB bool
//	}
//)
//
//func TestMarshal(t *testing.T) {
//	data := &MyStruct{
//		Field1: 42,
//		Field2: "Hello, World!",
//		OtherStruct: OtherStruct{
//			FieldA: 3.14159,
//			FieldB: true,
//		},
//	}
//
//	serializedData, err := Marshal(data)
//	assert.Nil(t, err)
//
//	t.Log(fmt.Sprintf("%v", serializedData))
//
//	var deserializedData MyStruct
//	assert.Nil(t, Unmarshal(serializedData, &deserializedData))
//
//	assert.Equal(t, data, deserializedData)
//}
