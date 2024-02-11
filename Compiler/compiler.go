package Compiler

import (
	"ProtoCompiler/Compiler/Protobuf"
	"regexp"
	"strings"
)

const (
	PROTOBUF = "protobuf"
)

type Compiler struct {
	sourceCode string // source code
	sourceType string // type of the source code
}

func NewCompiler(sourceCode string, sourceType string) *Compiler {
	return &Compiler{sourceCode: sourceCode, sourceType: sourceType}
}

func (c *Compiler) ParseProtobuf() *Protobuf.Protobuf {
	// first iterate over the source code and add \n in front and back of '}' to make sure that all char '}' is in an independent line.
	// this will be used in later parsing process.
	var source strings.Builder
	for _, char := range c.sourceCode {
		if char == '}' {
			source.WriteRune('\n')
		}
		source.WriteRune(char)
		if char == '}' {
			source.WriteRune('\n')
		}
	}

	lines := regexp.MustCompile(`\r?\n`).Split(source.String(), -1)
	return Protobuf.ParseProtobuf(lines)
}

// Parse the source code based on the sourceType
func (c *Compiler) Parse() interface{} {
	if c.sourceType == PROTOBUF {
		return c.ParseProtobuf()
	}

	return nil
}
