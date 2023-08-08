package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"text/template"

	"github.com/giftDad/protoc-gen-validate/templates"
	"github.com/giftDad/protoc-gen-validate/templates/rule"

	"google.golang.org/protobuf/compiler/protogen"
)

type validate struct {
	plugin *protogen.Plugin
}

func (t validate) Generate(plugin *protogen.Plugin) error {
	t.plugin = plugin

	for _, f := range t.plugin.Files {
		if len(f.Services) == 0 {
			continue
		}

		t.generateValidate(f)
	}

	return nil
}

func (t *validate) generateValidate(file *protogen.File) {
	fname := file.GeneratedFilenamePrefix + ".validate.go"

	tpl := template.New(fname)
	rule.RegisterFunctions(tpl)
	templates.Register(tpl)

	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, file); err != nil {
		panic(err)
	}

	gf := t.plugin.NewGeneratedFile(fname, file.GoImportPath)
	gf.Write(t.formattedOutput(buf.Bytes()))
}

func (t *validate) formattedOutput(raw []byte) []byte {
	// Reformat generated code.
	fset := token.NewFileSet()
	ast, err := parser.ParseFile(fset, "", raw, parser.ParseComments)
	if err != nil {
		// Print out the bad code with line numbers.
		// This should never happen in practice, but it can while changing generated code,
		// so consider this a debugging aid.
		var src bytes.Buffer
		s := bufio.NewScanner(bytes.NewReader(raw))
		for line := 1; s.Scan(); line++ {
			fmt.Fprintf(&src, "%5d\t%s\n", line, s.Bytes())
		}
		log.Fatal("bad Go source code was generated:", err.Error(), "\n"+src.String())
	}

	out := bytes.NewBuffer(nil)
	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(out, fset, ast)
	if err != nil {
		log.Fatal("generated Go source code could not be reformatted:", err.Error())
	}

	return out.Bytes()
}
