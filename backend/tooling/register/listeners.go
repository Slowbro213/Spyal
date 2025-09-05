package register

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"spyal/tooling/utils"
)

func GetListenerConstructors(dir string) ([]string, error) {
	c := cases.Title(language.English)

	files, err := utils.ListAllFiles(dir)
	if err != nil {
		return nil, err
	}

	suffix := c.String(strings.TrimSuffix(dir, "s"))
	pattern := `^func (New.*` + suffix + `)\(`

	var constructors []string
	for _, path := range files {
		constructor, err := utils.FindMatchingLine(path, pattern)
		if err != nil || constructor == "" {
			continue
		}
		constructors = append(constructors, constructor)
	}

	return constructors, nil
}

//gosec:disable
func GenerateListenerRegistry(dir string, constructors []string) error {
	f, err := os.Create(filepath.Join(dir, "registry.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	singular := strings.TrimSuffix(dir, "s")
	c := cases.Title(language.English)
	capitalized := c.String(singular)

	// Set type and var name
	varName := capitalized + "Registry"
	interfaceName := capitalized

	// Template string with dynamic package and variable
	tmpl := GeneratedCodeMsg + `
package {{ .Package }}

import "spyal/contracts"

var {{ .VarName }} = []contracts.New{{.Interface}}Func {
{{- range .Constructors }}
	{{ . }},
{{- end }}
}
`

	// Render it
	t := template.Must(template.New("registry").Parse(tmpl))
	return t.Execute(f, struct {
		Package      string
		VarName      string
		Interface    string
		Constructors []string
	}{
		Package:      dir,
		VarName:      varName,
		Interface:    interfaceName,
		Constructors: constructors,
	})
}
