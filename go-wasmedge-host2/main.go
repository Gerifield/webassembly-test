package main

import (
	"flag"
	"fmt"
	"os"
	"unsafe"

	"github.com/second-state/WasmEdge-go/wasmedge"
)

func main() {
	fileName := flag.String("wasm", "module.wasm", "WASM module to load")
	flag.Parse()
	wasmedge.SetLogErrorLevel()
	conf := wasmedge.NewConfigure(wasmedge.WASI)
	// vm := wasmedge.NewVMWithConfig(conf)
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

	//vm.RunWasmFile(*fileName, "_start")

	//res, err := vm.Execute("Fv1", int32(22), int32(2))

	// Nope, error
	// res, err := vm.Execute("Fv3", int32(22))
	// fmt.Println(res[0].(int32))

	n := int32(10)

	// Get some memory
	ptr, err := vm.Execute("malloc", n)
	if err != nil {
		fmt.Println("malloc failed:", err)
	}
	// Mem. address
	fmt.Printf("inputPtr memory at: %p\n", unsafe.Pointer((uintptr)(ptr[0].(int32))))

	// Add some data to the memory
	inputPointer := ptr[0].(int32)
	// Write the subject into the memory.
	mem := store.FindMemory("memory")
	memData, _ := mem.GetData(uint(inputPointer), uint(n))
	copy(memData, []byte{22, 33, 44, 55, 66})

	// Execute function
	fv4RetPtr, err := vm.Execute("Fv4", n, ptr[0])
	if err != nil {
		fmt.Println("fv4 call failed:", err)
	} else {

		// Read back the returned memory address (same size?)
		fmt.Printf("fv4RetPtr memory at: %p\n", unsafe.Pointer((uintptr)(fv4RetPtr[0].(int32))))
		mem := store.FindMemory("memory")
		if mem != nil {
			// int32 occupies 4 bytes

			fv4ReturnMemory, err := mem.GetData(uint(fv4RetPtr[0].(int32)), uint(n))
			if err == nil && fv4ReturnMemory != nil {
				fmt.Println("fv4ReturnMemory:", fv4ReturnMemory)
			}
		}
	}

	_, err = vm.Execute("free", ptr...)
	if err != nil {
		fmt.Println("free failed:", err)
	}

	// Question: Do we need to free fv4RetPtr???
	_, err = vm.Execute("free", fv4RetPtr...)
	if err != nil {
		fmt.Println("free failed:", err)
	}

	exitcode := wasi.WasiGetExitCode()
	if exitcode != 0 {
		fmt.Println("Go: Running wasm failed, exit code:", exitcode)
	}

	vm.Release()
	store.Release()
	conf.Release()
}
