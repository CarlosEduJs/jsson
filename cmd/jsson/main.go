package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"jsson/internal/transpiler"
	"os"
	"path/filepath"
)

func main() {
	inputPtr := flag.String("i", "", "Input JSSON file")
	// include merge mode: keep (default), overwrite, error
	mergeMode := flag.String("include-merge", "keep", "Include merge strategy: keep|overwrite|error")
	flag.Parse()

	if *inputPtr == "" {
		fmt.Println("Please provide an input file with -i")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(*inputPtr)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Resolve the absolute path of the input file and pass its directory as baseDir
	absInput, err := filepath.Abs(*inputPtr)
	if err != nil {
		fmt.Printf("Error resolving input path: %v\n", err)
		os.Exit(1)
	}
	baseDir := filepath.Dir(absInput)

	l := lexer.New(string(data))
	l.SetSourceFile(absInput)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, msg := range p.Errors() {
			fmt.Println("\t" + msg)
		}
		os.Exit(1)
	}

	t := transpiler.New(program, baseDir, *mergeMode, absInput)
	jsonOutput, err := t.Transpile()
	if err != nil {
		fmt.Printf("Transpilation error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonOutput))
}
