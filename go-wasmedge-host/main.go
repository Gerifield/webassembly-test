package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/second-state/WasmEdge-go/wasmedge"
)

func main() {
	fileName := flag.String("wasm", "module.wasm", "WASM module to load")
	flag.Parse()
	wasmedge.SetLogErrorLevel()
	conf := wasmedge.NewConfigure(wasmedge.WASI)

	store := wasmedge.NewStore()
	vm := wasmedge.NewVMWithConfigAndStore(conf, store)

	wasi := vm.GetImportObject(wasmedge.WASI)
	wasi.InitWasi(
		os.Args[1:],
		os.Environ(),
		[]string{".:."},
	)

	err := vm.LoadWasmFile(*fileName)
	if err != nil {
		fmt.Println("failed to load wasm")
	}
	vm.Validate()
	vm.Instantiate()

	subject := "Emberek Webassemblybol!"
	lengthOfSubject := len(subject)

	// Allocate memory for the subject, and get a pointer to it.
	// Include a byte for the NULL terminator we add below.
	allocateResult, _ := vm.Execute("malloc", int32(lengthOfSubject+1))
	inputPointer := allocateResult[0].(int32)

	// Write the subject into the memory.
	mem := store.FindMemory("memory")
	memData, _ := mem.GetData(uint(inputPointer), uint(lengthOfSubject+1))
	copy(memData, subject)
	// C-string terminates by NULL.
	memData[lengthOfSubject] = 0

	// Run the `greet` function. Given the pointer to the subject.
	fv3, _ := vm.Execute("Fv3", inputPointer)
	outputPointer := fv3[0].(int32)

	pageSize := mem.GetPageSize()
	// Read the result of the `greet` function.
	memData, _ = mem.GetData(uint(0), uint(pageSize*65536))

	var output strings.Builder
	nth := 0
	for {
		if memData[int(outputPointer)+nth] == 0 {
			break
		}

		output.WriteByte(memData[int(outputPointer)+nth])
		nth++
	}

	fmt.Println(output.String())

	// Deallocate the subject, and the output.
	vm.Execute("free", inputPointer)
	vm.Execute("free", outputPointer)

	vm.Release()
	store.Release()
	conf.Release()
}
