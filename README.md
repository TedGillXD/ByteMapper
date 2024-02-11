# ByteMapper Library
The **ByteMapper** library offers a powerful toolset for developers working with Protobuf data structures, enabling the generation of metadata from .proto files. This capability allows users to craft custom Protobuf parsers tailored for specific optimization goals, unlocking new levels of efficiency and performance in data handling and processing.

## Key Features
* **Metadata Generation:** ByteMapper reads .proto files and extracts the data structure, transforming it into rich metadata. This metadata serves as the foundation upon which users can build specialized parsers, offering a pathway to optimizations that are not readily achievable with generic parsing solutions.

* **Custom Parser Development:** Armed with the generated metadata, developers have the flexibility to design and implement Protobuf parsers that are optimized for their specific use cases. Whether the focus is on reducing memory usage, increasing parsing speed, or enhancing data validation, ByteMapper provides the necessary groundwork.

* **C++ Source File Generation:** As an additional feature (currently in development), ByteMapper aims to provide simple C++ source file generation. This will further ease the process of integrating optimized Protobuf parsing into existing C++ projects, bridging the gap between Protobuf schemas and performant, type-safe code.

## Potential Use Cases
* **High-Performance Applications:** Ideal for scenarios where processing speed and efficiency are paramount, such as in high-frequency trading platforms, real-time data analytics, and game development.

* **Resource-Constrained Environments:** Enables optimizations that reduce memory footprint and CPU usage, perfect for embedded systems, IoT devices, and mobile applications where resources are limited.

* **Custom Data Processing Pipelines:** Facilitates the creation of bespoke data processing and validation pipelines that are finely tuned to the specific requirements of complex projects.

## Installation
```shell
go get github.com/TedGillXD/ByteMapper
```

## Getting Started
* To generate metadata of protobuf data struct
    ```go
    compiler := Compiler.NewCompiler(sourceCode, Compiler.PROTOBUF)
    // The struct of protobufObj is defined in Compiler/Protobuf/dataStructure.go
    protobufObj := compiler.ParseProtobuf()
    ```

* To cpp source file from metadata
    ```go
    generator := TargetGenerator.NewTargetGenerator(protobufObj)
    // get source and header in string
	source, header := generator.GetCppFile()
    ```

* To change the template
    ```go
    // In TargetGenerator/Templates/cppTemplates.go
    const CppHeaderTemplate = `change the template to whatever you want!`
    const CppHeaderMessageTemplate = `change the template to whatever you want!`
    const CppHeaderEnumTemplate = `change the template to whatever you want!`
    ```
    `targetGenerator.go` will fill the template based on thoes three templates, the metadata will be provided by the structures shows below:
    ```go
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
        Includes         []string
        MessagesAndEnums []string
    }
    ```
    This program use `text/template` for parsing. So your template will work well if you follow the instruction of the `text/template` module.