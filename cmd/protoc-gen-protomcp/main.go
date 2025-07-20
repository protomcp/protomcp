package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var flags flag.FlagSet
	interfaces := flags.Bool("interfaces", DefaultGenerateInterfaces, "generate interface files")
	services := flags.Bool("services", DefaultGenerateServices, "generate service interfaces")
	pattern := flags.String("interface_pattern", DefaultInterfacePattern,
		"pattern for interface names (e.g., 'I%' for prefix, '%Interface' for suffix)")
	noImpl := flags.Bool("noimpl", DefaultGenerateNoImpl, "generate NoImpl structs for interfaces")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		// Declare support for proto3 optional fields
		plugin.SupportedFeatures = SupportedFeatures

		return run(plugin, &GeneratorOptions{
			GenerateInterfaces: *interfaces,
			GenerateServices:   *services,
			InterfacePattern:   *pattern,
			GenerateNoImpl:     *noImpl,
		})
	})
}

func run(plugin *protogen.Plugin, opts *GeneratorOptions) error {
	gen := NewGenerator(plugin)

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}

		if err := gen.GenerateFile(file, opts); err != nil {
			return fmt.Errorf("failed to generate file %s: %w", file.Desc.Path(), err)
		}
	}

	return nil
}
