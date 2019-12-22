package loader

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func BenchmarkA(b *testing.B) {
	s := Bytes("0123456789")
	fmt.Println(String(s))
}

func BenchmarkA2(b *testing.B) {
	s := Bytes2("0123456789")
	fmt.Println(String2(s))
}

func String2(bytes []byte) string {
	return String(bytes)
}

func Bytes2(str string) []byte {
	return []byte(str)
}

func String(bytes []byte) string {
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
	}))
}

func Bytes(str string) []byte {
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
}
