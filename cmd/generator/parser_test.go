package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const corruptFile = `
package math

type Calc interface {
    Add(a int, b int) int
`

const testFile = `
package math

type Calc1 interface {
    Add(a int, b double) (c int)
}

type Calc2 interface {
    Add(int, double) (int, error)
}
`

func TestParserError(t *testing.T) {
	r := strings.NewReader(corruptFile)
	_, err := Parse(r)
	assert.Error(t, err)
}

func TestParser(t *testing.T) {
	r := strings.NewReader(testFile)
	d, err := Parse(r)
	assert.NoError(t, err)

	assert.Equal(t, "math", d.Package)

	iface := d.interfaces[0]
	assert.Equal(t, "Calc1", iface.Name)
	method := iface.methods[0]
	assert.Equal(t, "Add", method.Name)
	assert.Equal(t, "a", method.Params[0].Name)
	assert.Equal(t, "int", method.Params[0].Type)
	assert.Equal(t, "b", method.Params[1].Name)
	assert.Equal(t, "double", method.Params[1].Type)
	assert.Equal(t, "c", method.Results[0].Name)
	assert.Equal(t, "int", method.Results[0].Type)

	iface = d.interfaces[1]
	assert.Equal(t, "Calc2", iface.Name)
	method = iface.methods[0]
	assert.Equal(t, "Add", method.Name)
	assert.Equal(t, "param1", method.Params[0].Name)
	assert.Equal(t, "int", method.Params[0].Type)
	assert.Equal(t, "param2", method.Params[1].Name)
	assert.Equal(t, "double", method.Params[1].Type)
	assert.Equal(t, "res1", method.Results[0].Name)
	assert.Equal(t, "int", method.Results[0].Type)
	assert.Equal(t, "res2", method.Results[1].Name)
	assert.Equal(t, "error", method.Results[1].Type)
}
