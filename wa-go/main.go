package main

import (
	"fmt"
	"strings"
	"unsafe"
)

func main() {
	fmt.Println("Hello Webassembly")
}

//export Fv1
func Fv1(n int32) string {
	fmt.Printf("Wasi (kinda) rocks: %d\n", n)
	return "ok"
}

//export Fv2
func Fv2(s string) string {
	fmt.Printf("Wasi (kinda) rocks: %s\n", s)
	return "ok"
}

//export Fv3
func Fv3(ptr *int32) *byte {
	var buff strings.Builder
	startPoint := uintptr(unsafe.Pointer(ptr))
	step := 0
	for {
		s := *(*int32)(unsafe.Pointer(startPoint + uintptr(step)))
		if s == 0 {
			break
		}

		buff.WriteByte(byte(s))
		step++
	}

	param := buff.String()
	output := "Helloka " + param

	return &(([]byte)(output)[0])
}

//export Fv4
func Fv4(n int32, p *byte) *byte {
	ret := make([]byte, n) // return and input size are the same

	inp := unsafe.Slice(p, n)
	fmt.Println(inp)

	copy(ret, inp)

	return &ret[0]
}
