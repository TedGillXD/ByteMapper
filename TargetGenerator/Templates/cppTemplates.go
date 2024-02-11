package Templates

type CppEnum struct {
	Name       string
	EnumValues []CppEnumValue
}

type CppEnumValue struct {
	Name  string
	Value string
}

type CppMessage struct {
	Name                   string
	NestedMessagesAndEnums []string
	MessageFields          []CppMessageField
}

type CppMessageField struct {
	Name         string
	Type         string
	IsRepeatable bool
}

type CppHeader struct {
	HeaderGuard      string
	Package          string
	Includes         []string // the header file that should include in the cpp header file
	MessagesAndEnums []string // all cppHeaderMessageTemplate and cppHeaderEnumTemplate strings store in this slice
}

const CppHeaderTemplate = `
#ifndef {{.HeaderGuard}}_H
#define {{.HeaderGuard}}_H

#include <string>
#include <vector>
{{range .Includes}}
#include <{{.}}>
{{end}}

namespace {{.Package}} {
	
	// top messages and enums declarations
	{{range .MessagesAndEnums}}
	{{.}}
	{{end}}

} // namespace {{.HeaderGuard}}_H

#endif // {{.HeaderGuard}}_H
`

const CppHeaderMessageTemplate = `
class {{.Name}} {
// constructor and destructor
public:
	{{.Name}}() = default;
	~{{.Name}}() = default;

// nested message definition
public:
	{{range .NestedMessagesAndEnums}}
	{{.}}
	{{end}}

public:
	void parseFromStream(std::vector<uint8_t>& stream);

// getter and setter
public:
    {{range .MessageFields}}
    // Getter for {{.Name}}
    const {{if .IsRepeatable}}std::vector<{{.Type}}>{{else}}{{.Type}}{{end}}& get_{{.Name}}() const;

    // Setter for {{.Name}}
    void set_{{.Name}}(const {{if .IsRepeatable}}std::vector<{{.Type}}>&{{else}}{{.Type}}{{end}} value);
    {{end}}

private:
	{{range .MessageFields}}
	{{.Type}} {{.Name}};
	{{end}}
};
`

const CppHeaderEnumTemplate = `
enum class {{.Name}} {
	{{range .EnumValues}}
	{{.Name}} = {{.Value}},
	{{end}}
};
`
