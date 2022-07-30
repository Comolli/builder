package utils

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"unsafe"
)

//get bytes from string safely
func UnsafeGetBytes(s string) []byte {
	return (*[0x7fff0000]byte)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data),
	)[:len(s):len(s)]
}

//byte2string
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		fmt.Printf("painic error stack:%s:%s", p, Bytes2String(debug.Stack()))
	}
}
