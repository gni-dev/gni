package generator

import (
	"html/template"
	"io"
)

const golangTmp = `
package {{.Package}}
`

func Golang(d *Definition, w io.Writer) error {
	t, err := template.New("golang").Parse(golangTmp)
	if err != nil {
		return err
	}

	return t.Execute(w, d)
}
