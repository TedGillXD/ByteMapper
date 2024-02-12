package Translators

import (
	"ProtoCompiler/Compiler/Protobuf"
	"ProtoCompiler/TargetGenerator/Templates"
	"bytes"
	"strconv"
	"strings"
	"text/template"
)

var cppTypeMapping = map[string]string{
	"double":   "double",
	"float":    "float",
	"int32":    "int32_t",
	"int64":    "int64_t",
	"uint32":   "uint32_t",
	"uint64":   "uint64_t",
	"sint32":   "int32_t",
	"sint64":   "int64_t",
	"fixed32":  "uint32_t",
	"fixed64":  "uint64_t",
	"sfixed32": "int32_t",
	"sfixed64": "int64_t",
	"bool":     "bool",
	"string":   "std::string",
	"bytes":    "std::vector<uint8_t>",
}

func translateFieldToCpp(field *Protobuf.Field) string {
	cppType := ""
	if field.IsBaseField {
		cppType = cppTypeMapping[field.Type]
	} else {
		cppType = field.Type
	}

	if field.IsRepeated {
		cppType = "std::vector<" + cppType + ">"
	}

	return cppType
}

func translateMessageFieldToCpp(msg *Protobuf.Message) {
	for _, field := range msg.Fields {
		cppType := translateFieldToCpp(field)
		field.Type = cppType
	}
	// process nested message
	for _, nestedMsg := range msg.NestedMessages {
		translateMessageFieldToCpp(nestedMsg)
	}
}

func translateTypeToCpp(protobufObj *Protobuf.Protobuf) {
	for _, msg := range protobufObj.Messages {
		translateMessageFieldToCpp(msg)
	}
}

// getCppEnumObject will generate the CppEnum object for template
func getCppEnumObject(enum *Protobuf.Enum) Templates.CppEnum {
	var enumValues []Templates.CppEnumValue
	for _, value := range enum.Values {
		enumValues = append(enumValues, Templates.CppEnumValue{Name: value.Name, Value: strconv.Itoa(value.Value)})
	}

	return Templates.CppEnum{Name: enum.Name, EnumValues: enumValues}
}

func getCppMessageObject(message *Protobuf.Message) Templates.CppMessage {
	// 1. process message fields
	var fields []Templates.CppMessageField
	typesString := ""
	for _, field := range message.Fields {
		fields = append(fields, Templates.CppMessageField{
			Name:         field.Name,
			Type:         field.Type,
			IsRepeatable: field.IsRepeated,
		})
		typesString += field.Type + ","
	}
	typesString = strings.TrimSuffix(typesString, ",")

	var nestedMessagesAndEnums []string

	// 2. process nested enums
	for _, enum := range message.NestedEnums {
		// concat the new string into the end of slice
		nestedMessagesAndEnums = append(nestedMessagesAndEnums, getCppEnumString(getCppEnumObject(enum)))
	}

	// 3. process nested messages
	for _, msg := range message.NestedMessages {
		// concat the new string into the end of slice
		nestedMessagesAndEnums = append(nestedMessagesAndEnums, getCppMessageString(getCppMessageObject(msg)))
	}

	return Templates.CppMessage{Name: message.Name, MessageFields: fields, NestedMessagesAndEnums: nestedMessagesAndEnums, TypesString: typesString}
}

func getCppEnumString(cppEnum Templates.CppEnum) string {
	enumTemplate, err := template.New("enumParser").Parse(Templates.CppHeaderEnumTemplate)
	if err != nil {
		panic("Failed to parse cppEnum" + cppEnum.Name + " " + err.Error())
	}
	var buf bytes.Buffer
	err = enumTemplate.Execute(&buf, cppEnum)
	if err != nil {
		panic("Failed to parse object cppEnum: " + cppEnum.Name + " " + err.Error())
	}

	return buf.String()
}

func getCppMessageString(cppMessage Templates.CppMessage) string {
	messageTemplate, err := template.New("messageParser").Parse(Templates.CppHeaderMessageTemplate)
	if err != nil {
		panic("Failed to parse cppEnum" + cppMessage.Name + " " + err.Error())
	}
	var buf bytes.Buffer
	err = messageTemplate.Execute(&buf, cppMessage)
	if err != nil {
		panic("Failed to parse object cppMessage: " + cppMessage.Name + " " + err.Error())
	}

	return buf.String()
}

func getCppHeaderString(cppHeader Templates.CppHeader) string {
	headerTemplate, err := template.New("headerTemplate").Parse(Templates.CppHeaderTemplate)
	if err != nil {
		panic("Failed to parse cppHeader template: " + err.Error())
	}
	var bufHeader bytes.Buffer
	err = headerTemplate.Execute(&bufHeader, cppHeader)
	if err != nil {
		panic("Failed to parse object cppHeader: " + err.Error())
	}

	return bufHeader.String()
}

func GetCppFromProtobuf(protobufObj *Protobuf.Protobuf, fileName string) (source string, header string) {
	// 1. translate the protobuf object for cpp translation
	translateTypeToCpp(protobufObj)

	// 2. make cpp header file string
	var cppHeader Templates.CppHeader
	cppHeader.HeaderGuard = fileName // using file name as header guard
	cppHeader.Package = protobufObj.Package

	var topMessagesAndEnums []string

	// process top enums
	for _, enum := range protobufObj.Enums {
		topMessagesAndEnums = append(topMessagesAndEnums, getCppEnumString(getCppEnumObject(enum)))
	}

	// process top messages
	for _, msg := range protobufObj.Messages {
		topMessagesAndEnums = append(topMessagesAndEnums, getCppMessageString(getCppMessageObject(msg)))
	}

	cppHeader.MessagesAndEnums = topMessagesAndEnums

	// TODO: 3. make cpp source file string

	return "", getCppHeaderString(cppHeader)
}
