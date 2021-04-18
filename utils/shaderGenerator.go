package utils

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var fileTemplate = `
// +build {{ .BuildTags }}

package {{ .PackageName }}

// This file was auto-generated

{{ .Shaders }}
`

var shaderTemplate = `const {{ .Name }}Str = "{{ escapeNewLine .Code }}{{ "\\x00" }}"`

type shaderFileData struct {
	BuildTags   string
	PackageName string
	Shaders     string
}

type shaderData struct {
	Name string
	Code string
}

var funcMap = template.FuncMap{
	"escapeNewLine": func(input string) string {
		return strings.ReplaceAll(input, "\n", "\\n")
	},
}

func GenerateShaders(shadersDirPath, targetFilePath, packageName string) {
	fileTmpl := template.Must(template.New("shaders").Parse(fileTemplate))
	shaderTmpl := template.Must(template.New("shader").Funcs(funcMap).Parse(shaderTemplate))

	outFile, err := os.Create(targetFilePath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	sb := strings.Builder{}

	files, err := ioutil.ReadDir(shadersDirPath)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		shader, _ := ioutil.ReadFile(shadersDirPath + "/" + file.Name())

		_ = shaderTmpl.Execute(&sb, shaderData{
			Name: strings.Split(file.Name(), ".")[0],
			Code: string(shader),
		})

		sb.WriteString("\n\n")
	}

	_ = fileTmpl.Execute(outFile, shaderFileData{
		BuildTags:   "darwin linux windows",
		PackageName: packageName,
		Shaders:     sb.String(),
	})
}
