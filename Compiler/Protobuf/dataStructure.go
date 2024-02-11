package Protobuf

type Protobuf struct {
	Messages []*Message // top messages
	Enums    []*Enum    // top enums
	Syntax   string
	Package  string
}

type Message struct {
	Name           string     // message name
	Fields         []*Field   // message fields
	NestedMessages []*Message // nested defined message
	NestedEnums    []*Enum    // nested defined enum
}

type Field struct {
	Type        string
	Name        string // field name
	Tag         int    // field tag: like 1, 2, 3
	IsBaseField bool   // true: like int32 etc		 false: message
	IsRepeated  bool
}

type Enum struct {
	Name   string // enum name
	Values []*EnumValue
}

type EnumValue struct {
	Name  string // enum name
	Value int    // enum value
}
