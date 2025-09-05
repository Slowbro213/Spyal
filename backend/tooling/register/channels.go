package register

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)


//nolint
func GetChannelConstructors(dir string) ([]string, error) {
	fset := token.NewFileSet()
	channelTypes := []string{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") || info.Name() == "registry.go" {
			return err
		}

		node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			return err
		}

		for _, decl := range node.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}

			for _, spec := range gen.Specs {
				ts := spec.(*ast.TypeSpec)
				structType, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}

				for _, field := range structType.Fields.List {
					if len(field.Names) == 0 {
						ident, ok := field.Type.(*ast.SelectorExpr)
						if ok && ident.Sel.Name == "Channel" {
							channelTypes = append(channelTypes, ts.Name.Name)
						}
					}
				}
			}
		}

		return nil
	})

	return channelTypes, err
}

//gosec:disable
func GenerateChannelRegistryFile(dir string, constructors []string) error {
	outFile := filepath.Join(dir, "registry.go")

	out := GeneratedCodeMsg + `
package channels

import (
	"spyal/contracts"
)

//nolint:gochecknoglobals
var Channels = map[string]contracts.Channel{
`

	for _, t := range constructors {
		name := strings.ToLower(t)
		out += fmt.Sprintf("    \"%s\": New%sChannel(),\n", name, t)
	}

	out += "}\n"

	writeFilePermissions := 0600
	return os.WriteFile(outFile, []byte(out), os.FileMode(writeFilePermissions))
}

