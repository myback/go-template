package main

import "C"
import (
	"bytes"

	tplFunc "go-template/go/template-functions"
	"go-template/go/utils"
)

//export Render
func Render(srcFile, valuesData *C.char) *C.char {
	out := new(bytes.Buffer)

	file := C.GoString(srcFile)
	data := C.GoString(valuesData)

	values := utils.MustUnmarshalValues(data)
	utils.CheckErr(tplFunc.Render(file, out, values))

	return C.CString(out.String())
}

//export RenderFomValuesFile
func RenderFomValuesFile(srcFile, valuesFile *C.char) *C.char {
	values, err := utils.ReadValuesFile(C.GoString(valuesFile))
	utils.CheckErr(err)

	out := new(bytes.Buffer)
	utils.CheckErr(tplFunc.Render(C.GoString(srcFile), out, values))

	return C.CString(out.String())
}

func main() {}
