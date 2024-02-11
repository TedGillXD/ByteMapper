package Protobuf

import (
	"regexp"
	"strconv"
	"strings"
)

func isBaseField(fieldType string) bool {
	baseTypes := []string{
		"double", "float",
		"int32", "int64", "uint32", "uint64",
		"sint32", "sint64", "fixed32", "fixed64",
		"sfixed32", "sfixed64", "bool", "string", "bytes",
	}

	for _, t := range baseTypes {
		if fieldType == t {
			return true
		}
	}

	return false
}

func parseSyntax(lines []string) (int, string) {
	syntaxRegex := regexp.MustCompile(`^syntax\s*=\s*"(proto2|proto3)";`)
	for index, line := range lines {
		if syntaxMatch := syntaxRegex.FindStringSubmatch(line); syntaxMatch != nil {
			if syntaxMatch[1] == "proto2" { // TODO: add proto2 support
				panic("this protobuf compiler only support proto3 yet.")
				return 0, ""
			}
			return index, syntaxMatch[1]
		}
	}

	panic("can not find the syntax from this .proto file")
	return 0, ""
}

func parsePackage(lines []string) (int, string) {
	packageRegex := regexp.MustCompile(`^package\s+([\w.]+);`)
	for index, line := range lines {
		if packageMatch := packageRegex.FindStringSubmatch(line); packageMatch != nil {
			return index, packageMatch[1]
		}
	}

	panic("can not find the package name from this .proto file")
	return 0, ""
}

func parseField(lines []string, current int) (int, *Field) {
	fieldRegex := regexp.MustCompile(`^(repeated\s+)?(\w+)\s+(\w+)\s+=\s+(\d+);`)
	match := fieldRegex.FindStringSubmatch(lines[current])

	if match != nil {
		fieldType := match[2]
		fieldName := match[3]
		fieldTag, _ := strconv.Atoi(match[4])

		field := Field{
			Type:        fieldType,
			Name:        fieldName,
			Tag:         fieldTag,
			IsBaseField: isBaseField(fieldType),
			IsRepeated:  match[1] != "",
		}
		return current + 1, &field
	}

	return current + 1, nil
}

func parseEnumValue(lines []string, current int) (int, *EnumValue) {
	enumValueRegex := regexp.MustCompile(`^(\w+)\s*=\s*(\d+);`)
	matches := enumValueRegex.FindStringSubmatch(lines[current])

	if len(matches) > 2 {
		valueName := matches[1]
		value, err := strconv.Atoi(matches[2])
		if err != nil {
			panic("the value of a enum value is not a integer!")
			return current + 1, nil
		}

		return current + 1, &EnumValue{Value: value, Name: valueName}
	}

	// add 1 to move to next line
	return current + 1, nil
}

func parseEnum(lines []string, current int) (int, *Enum) {
	enum := &Enum{} // first parse the declaration of enum to get the name
	enumRegex := regexp.MustCompile(`^enum\s+(\w+)`)
	matches := enumRegex.FindStringSubmatch(lines[current])
	if len(matches) < 2 {
		panic("cannot parse the declaration of enum in line " + lines[current])
	}
	enum.Name = matches[1]
	current++ // move to next line

	// this loop will exit when it encounters the right brace
	for current < len(lines) && !strings.HasPrefix(lines[current], "}") {
		// parsing definition of enum
		var enumValue *EnumValue
		current, enumValue = parseEnumValue(lines, current)
		if enumValue != nil {
			enum.Values = append(enum.Values, enumValue)
		}
	}

	// add 1 to jump over the right brace
	return current + 1, enum
}

func parseMessage(lines []string, current int) (int, *Message) {
	// first parse the message declaration line to get its name
	msg := &Message{}
	messageRegex := regexp.MustCompile(`^message\s+(\w+)`)
	matches := messageRegex.FindStringSubmatch(strings.TrimSpace(lines[current]))
	if len(matches) < 2 {
		panic("cannot parse the declaration of message in line " + lines[current])
	}
	msg.Name = matches[1]
	current++ // move to next line for parsing field or enum definition

	for current < len(lines) && !strings.HasPrefix(lines[current], "}") {
		line := strings.TrimSpace(lines[current])
		if strings.HasPrefix(line, "message") {
			var nestedMsg *Message
			current, nestedMsg = parseMessage(lines, current)
			msg.NestedMessages = append(msg.NestedMessages, nestedMsg)
		} else if strings.HasPrefix(line, "enum") {
			var enumDef *Enum
			current, enumDef = parseEnum(lines, current)
			if enumDef != nil {
				msg.NestedEnums = append(msg.NestedEnums, enumDef)
			}
		} else {
			var field *Field
			current, field = parseField(lines, current)
			if field != nil {
				msg.Fields = append(msg.Fields, field)
			}
		}
	}
	return current + 1, msg
}

func parseTopMessagesAndEnums(lines []string, current int) ([]*Message, []*Enum) {
	var messages []*Message
	var enums []*Enum
	for current < len(lines) {
		if strings.HasPrefix(lines[current], "message") {
			var msg *Message
			current, msg = parseMessage(lines, current) // 跳过消息名称行
			messages = append(messages, msg)
		} else if strings.HasPrefix(lines[current], "enum") {
			var enum *Enum
			current, enum = parseEnum(lines, current)
			enums = append(enums, enum)
		} else {
			current++
		}
	}
	return messages, enums
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func ParseProtobuf(lines []string) *Protobuf {
	for index, line := range lines { // trim space for later use
		lines[index] = strings.TrimSpace(line)
	}

	protobufObj := &Protobuf{}
	current1 := 0
	current2 := 0
	current1, protobufObj.Syntax = parseSyntax(lines)
	current2, protobufObj.Package = parsePackage(lines)
	protobufObj.Messages, protobufObj.Enums = parseTopMessagesAndEnums(lines, maxInt(current1, current2))

	return protobufObj
}
