package main

import "fmt"

func main() {
	fmt.Println("Hello Webassembly")
}

//export Fv1
func Fv1(n int32) string {
	fmt.Printf("Wasi (kinda) rocks: %d\n", n)
	return "ok"
}
