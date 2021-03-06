package wasmedge

/*
#include <wasmedge/wasmedge.h>

typedef void (*wasmedgego_HostFuncWrapper)
  (void *, void *, WasmEdge_MemoryInstanceContext *, const WasmEdge_Value *, const uint32_t, WasmEdge_Value *, const uint32_t);

WasmEdge_Result wasmedgego_HostFuncInvoke(void *Func, void *Data,
                                  WasmEdge_MemoryInstanceContext *MemCxt,
                                  const WasmEdge_Value *Params, const uint32_t ParamLen,
                                  WasmEdge_Value *Returns, const uint32_t ReturnLen);
*/
import "C"
import (
	"errors"
	"reflect"
	"unsafe"
)

type Function struct {
	_inner *C.WasmEdge_FunctionInstanceContext
	_index uint
	_own   bool
}

type Table struct {
	_inner *C.WasmEdge_TableInstanceContext
	_own   bool
}

type Memory struct {
	_inner *C.WasmEdge_MemoryInstanceContext
	_own   bool
}

type Global struct {
	_inner *C.WasmEdge_GlobalInstanceContext
	_own   bool
}

func NewFunction(ftype *FunctionType, fn hostFunctionSignature, additional interface{}, cost uint) *Function {
	if ftype == nil {
		return nil
	}

	index := hostfuncMgr.add(fn, additional)
	function := C.WasmEdge_FunctionInstanceCreateBinding(
		ftype._inner,
		C.wasmedgego_HostFuncWrapper(C.wasmedgego_HostFuncInvoke),
		unsafe.Pointer(uintptr(index)),
		nil,
		C.uint64_t(cost))
	if function == nil {
		hostfuncMgr.del(index)
		return nil
	}
	return &Function{_inner: function, _index: index, _own: true}
}

func (self *Function) GetFunctionType() *FunctionType {
	return &FunctionType{
		_inner: C.WasmEdge_FunctionInstanceGetFunctionType(self._inner),
		_own:   false,
	}
}

func (self *Function) Release() {
	if self._own && self._inner != nil {
		C.WasmEdge_FunctionInstanceDelete(self._inner)
		hostfuncMgr.del(self._index)
	}
	self._inner = nil
	self._own = false
}

func NewTable(ttype *TableType) *Table {
	if ttype == nil {
		return nil
	}
	table := C.WasmEdge_TableInstanceCreate(ttype._inner)
	if table == nil {
		return nil
	}
	return &Table{_inner: table, _own: true}
}

func (self *Table) GetTableType() *TableType {
	return &TableType{
		_inner: C.WasmEdge_TableInstanceGetTableType(self._inner),
		_own:   false,
	}
}

func (self *Table) GetData(off uint) (interface{}, error) {
	cval := C.WasmEdge_Value{}
	res := C.WasmEdge_TableInstanceGetData(self._inner, &cval, C.uint32_t(off))
	if !C.WasmEdge_ResultOK(res) {
		return nil, newError(res)
	}
	return fromWasmEdgeValue(cval), nil
}

func (self *Table) SetData(data interface{}, off uint) error {
	cval := toWasmEdgeValue(data)
	res := C.WasmEdge_TableInstanceSetData(self._inner, cval, C.uint32_t(off))
	if !C.WasmEdge_ResultOK(res) {
		return newError(res)
	}
	return nil
}

func (self *Table) GetSize() uint {
	return uint(C.WasmEdge_TableInstanceGetSize(self._inner))
}

func (self *Table) Grow(size uint) error {
	res := C.WasmEdge_TableInstanceGrow(self._inner, C.uint32_t(size))
	if !C.WasmEdge_ResultOK(res) {
		return newError(res)
	}
	return nil
}

func (self *Table) Release() {
	if self._own {
		C.WasmEdge_TableInstanceDelete(self._inner)
	}
	self._inner = nil
	self._own = false
}

func NewMemory(mtype *MemoryType) *Memory {
	if mtype == nil {
		return nil
	}
	memory := C.WasmEdge_MemoryInstanceCreate(mtype._inner)
	if memory == nil {
		return nil
	}
	return &Memory{_inner: memory, _own: true}
}

func (self *Memory) GetMemoryType() *MemoryType {
	return &MemoryType{
		_inner: C.WasmEdge_MemoryInstanceGetMemoryType(self._inner),
		_own:   false,
	}
}

func (self *Memory) GetData(off uint, length uint) ([]byte, error) {
	p := C.WasmEdge_MemoryInstanceGetPointer(self._inner, C.uint32_t(off), C.uint32_t(length))
	if p == nil {
		return nil, errors.New("Failed get data pointer")
	}
	// Use SliceHeader to wrap the slice from cgo
	var r []byte
	s := (*reflect.SliceHeader)(unsafe.Pointer(&r))
	s.Cap = int(length)
	s.Len = int(length)
	s.Data = uintptr(unsafe.Pointer(p))
	return r, nil
}

func (self *Memory) SetData(data []byte, off uint, length uint) error {
	var ptrdata *C.uint8_t = nil
	if len(data) > 0 {
		ptrdata = (*C.uint8_t)(unsafe.Pointer(&data[0]))
	}
	res := C.WasmEdge_MemoryInstanceSetData(self._inner, ptrdata, C.uint32_t(off), C.uint32_t(length))
	if !C.WasmEdge_ResultOK(res) {
		return newError(res)
	}
	return nil
}

func (self *Memory) GetPageSize() uint {
	return uint(C.WasmEdge_MemoryInstanceGetPageSize(self._inner))
}

func (self *Memory) GrowPage(size uint) error {
	res := C.WasmEdge_MemoryInstanceGrowPage(self._inner, C.uint32_t(size))
	if !C.WasmEdge_ResultOK(res) {
		return newError(res)
	}
	return nil
}

func (self *Memory) Release() {
	if self._own {
		C.WasmEdge_MemoryInstanceDelete(self._inner)
	}
	self._inner = nil
	self._own = false
}

func NewGlobal(gtype *GlobalType, val interface{}) *Global {
	if gtype == nil {
		return nil
	}
	cval := toWasmEdgeValue(val)
	global := C.WasmEdge_GlobalInstanceCreate(gtype._inner, cval)
	if global == nil {
		return nil
	}
	return &Global{_inner: global, _own: true}
}

func (self *Global) GetGlobalType() *GlobalType {
	return &GlobalType{
		_inner: C.WasmEdge_GlobalInstanceGetGlobalType(self._inner),
		_own:   false,
	}
}

func (self *Global) GetValue() interface{} {
	cval := C.WasmEdge_GlobalInstanceGetValue(self._inner)
	return fromWasmEdgeValue(cval)
}

func (self *Global) SetValue(val interface{}) {
	C.WasmEdge_GlobalInstanceSetValue(self._inner, toWasmEdgeValue(val))
}

func (self *Global) Release() {
	if self._own {
		C.WasmEdge_GlobalInstanceDelete(self._inner)
	}
	self._inner = nil
	self._own = false
}
