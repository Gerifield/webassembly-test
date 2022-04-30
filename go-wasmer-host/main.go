package main

import (
	"fmt"
	"io/ioutil"
	"log"

	wasmer "github.com/wasmerio/wasmer-go/wasmer"
)

func main() {
	log.Println("Loading file")
	wasmBytes, err := ioutil.ReadFile("simple.wasm")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Init engine")
	engine := wasmer.NewEngine()
	log.Println("Init store")
	store := wasmer.NewStore(engine)

	// Compiles the module
	log.Println("Compile module")
	module, err := wasmer.NewModule(store, wasmBytes)
	if err != nil {
		log.Println(err)
		return
	}

	// Instantiates the module
	log.Println("Import objects")
	// importObject := wasmer.NewImportObject()
	wasiEnv, err := wasmer.NewWasiStateBuilder("wasi-program").
		// Choose according to your actual situation
		// Argument("--foo").
		// Environment("ABC", "DEF").
		// MapDirectory("./", ".").
		Finalize()
	if err != nil {
		log.Println(err)
		return
	}

	importObject, err := wasiEnv.GenerateImportObject(store, module)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Create instance")
	instance, err := wasmer.NewInstance(module, importObject)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(importObject)

	// Gets the `sum` exported function from the WebAssembly instance.
	log.Println("Looking for function")
	fv1Fn, err := instance.Exports.GetRawFunction("Fv1")
	if err != nil {
		log.Println(err)
		return
	}

	// for _, p := range fv1Fn.Type().Params() {
	// 	fmt.Println(p.Kind())
	// }

	// fmt.Println(fv1Fn.Type().Results())
	// fmt.Println(fv1Fn.ParameterArity())
	// fmt.Println(fv1Fn.ResultArity())

	// Calls that exported function with Go standard values. The WebAssembly
	// types are inferred and values are casted automatically.
	log.Println("Calling engine")
	// var ptr int32
	// result, err := fv1Fn.Call(ptr, 21) // first param is something elese, pointer maybe to a string response?
	result, err := fv1Fn.Native()(0, 21)
	if err != nil {
		log.Println(err)
		return
	}

	// CONCLUSION: If we use string with tinygo it generates an i32 first param which is bad probably

	fmt.Println(result)
}
