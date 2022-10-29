package main

import (
	"bytes"
	"flag"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	packageName  = flag.String("package", "", "name of the package, e.g. 'golang.org/x/text/encoding'")
	typeNames    = flag.String("type", "", "comma-separated list of type names; e.g. 'Decoder")
	outputPrefix = flag.String("prefix", "", "prefix to be added to the output file")
	outputSuffix = flag.String("suffix", "_wasm", "suffix to be added to the output file")
)

func main() {
	flag.Parse()
	if len(*typeNames) == 0 {
		log.Fatalf("the flag -type must be set")
	}
	types := strings.Split(*typeNames, ",")

	if packageName == nil || len(*packageName) < 1 {
		dir := "."
		dir, err := filepath.Abs(dir)
		if err != nil {
			log.Fatalf("unable to determine absolute filepath for requested path %s: %v",
				dir, err)
		}
		packageName = &dir
	}

	pkg, err := parsePackage(*packageName)
	if err != nil {
		log.Fatalf("parsing package: %v", err)
	}

	analysis := analysisData{
		Command:           strings.Join(os.Args[1:], " "),
		PackageName:       pkg.Name,
		PackageFQN:        pkg.Name,
		TypesAndValues:    make(map[string][]string),
		TargetPackageName: "main",
	}

	// Run generate for each type.
	for _, typeName := range types {
		log.Printf("typeName=" + typeName)
		values, err := pkg.ValuesOfType(typeName)
		parameters, err := pkg.findFunctionsAndParameters(typeName)
		if err != nil {
			log.Fatalf("finding functions and their parameters %v: %v", typeName, err)
		}
		log.Printf(strings.Join(values, ","))
		if err != nil {
			log.Fatalf("finding values for type %v: %v", typeName, err)
		}
		analysis.TypesAndValues[typeName] = values
		analysis.TypeName = values[0]
		analysis.WrapperTypeName = analysis.TypeName + "Wrapper"
		analysis.PackageFQN = pkg.FQN

		var buf bytes.Buffer
		if err := typeWrapperTemplate.Execute(&buf, analysis); err != nil {
			log.Fatalf("generating code: %v", err)
		}

		for key, val := range parameters {
			var vals []string
			for _, v := range val {
				vals = append(vals, v.Name)
			}
			data := functionData{
				FunctionName:    key,
				WrapperTypeName: analysis.WrapperTypeName,
				Types:           vals,
			}
			templ := functionWrapperTempl(data)
			buf.WriteString(templ)
		}

		src, err := format.Source(buf.Bytes())
		if err != nil {
			// Should never happen, but can arise when developing this code.
			// The user can compile the output to see the error.
			log.Printf("warning: internal error: invalid Go generated: %s", err)
			log.Printf("warning: compile the package to analyze the error")
			src = buf.Bytes()
		}

		output := strings.ToLower(*outputPrefix + typeName + *outputSuffix + "_js.go")
		outputPath := filepath.Join(".", output)
		if err := os.WriteFile(outputPath, src, 0644); err != nil {
			log.Fatalf("writing output: %s", err)
		}
	}
}
