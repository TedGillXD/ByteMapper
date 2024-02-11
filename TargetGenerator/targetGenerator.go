package TargetGenerator

import (
	"ProtoCompiler/Compiler/Protobuf"
	"ProtoCompiler/TargetGenerator/Translators"
)

type TargetGenerator struct {
	parsedObj interface{}
}

func NewTargetGenerator(parsedObj interface{}) *TargetGenerator {
	switch parsedObj.(type) {
	case *Protobuf.Protobuf:
	// nothing happened
	default:
		panic("passed a wrong data type into TargetGenerator!")
	}
	return &TargetGenerator{parsedObj: parsedObj}
}

func (t *TargetGenerator) GetCppFile() (header string, source string) {
	switch t.parsedObj.(type) {
	case *Protobuf.Protobuf:
		return Translators.GetCppFromProtobuf(t.parsedObj.(*Protobuf.Protobuf))

	}

	return "", ""
}
