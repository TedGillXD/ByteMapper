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
	TypesString            string
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
#include <cstddef>
#include <cstdint>
{{range .Includes}}
#include <{{.}}>
{{end}}

namespace {
	template<typename... Types>
	std::vector<std::size_t> getOffsetWithAlignment() {
		std::vector<std::size_t> offsets;
		std::size_t currentOffset = 0;
		std::size_t maxAlignment = 0;
	
		auto addOffset = [&](auto... args) {
			(..., ([&]{
				std::size_t alignment = alignof(decltype(args));
				std::size_t size = sizeof(decltype(args));
				maxAlignment = std::max(maxAlignment, alignment);
	
				std::size_t padding = (alignment - currentOffset % alignment) % alignment;
				currentOffset += padding;
				offsets.push_back(currentOffset);
				currentOffset += size;
			}()));
		};
	
		addOffset(Types()...);
	
		return offsets;
	}
}	// anonymous namespace

namespace {{.Package}} {
	// top messages and enums declarations
	{{range .MessagesAndEnums}}{{.}}
	{{end}}
} // namespace {{.HeaderGuard}}_H
#endif // {{.HeaderGuard}}_H
`

const CppHeaderMessageTemplate = `
class {{.Name}} {
public:
	{{.Name}}() = default;
	~{{.Name}}() = default;
public:
	{{range .NestedMessagesAndEnums}}
	{{.}}
	{{end}}

public:
	void parseFromStream(std::vector<uint8_t>& stream);

public:
    {{range .MessageFields}}
    [[nodiscard]] const {{.Type}}& get_{{.Name}}() const;
    void set_{{.Name}}({{.Type}} value);
    {{end}}

public:
	static std::vector<std::size_t> paramOffset;
	static std::unordered_map<int32_t, std::function<void(MessageBase*, void*)>> fieldHandler;
	std::unordered_map<int32_t, std::function<void (MessageBase *, void *)>> getFieldHandler() override { return {{.Name}}::fieldHandler; }
	std::vector<std::size_t> getFieldOffset() override { return {{.Name}}::paramOffset }

private:
	{{range .MessageFields}}{{.Type}} {{.Name}};
	{{end}}
};
`

const CppHeaderEnumTemplate = `
enum class {{.Name}} {
	{{range .EnumValues}}{{.Name}} = {{.Value}},
	{{end}}
};
`

const CppSourceTemplate = `

`

const CppSourceMessageTemplate = ``
