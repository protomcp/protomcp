package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var flags flag.FlagSet
	interfaces := flags.Bool("interfaces", true, "generate interface files")
	services := flags.Bool("services", true, "generate service interfaces")
	pattern := flags.String("interface_pattern", "I%",
		"pattern for interface names (e.g., 'I%' for prefix, '%Interface' for suffix)")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		// Declare support for proto3 optional fields
		plugin.SupportedFeatures = SupportedFeatures

		return run(plugin, &GeneratorOptions{
			GenerateInterfaces: *interfaces,
			GenerateServices:   *services,
			InterfacePattern:   *pattern,
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
