package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strconv"
)

type Param struct {
	Name string
	Type string
}

type Method struct {
	Name    string
	Params  []Param
	Results []Param
}

type Interface struct {
	Name    string
	methods []Method
}

type Definition struct {
	Package    string
	interfaces []Interface
}

func Parse(r io.Reader) (*Definition, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "", r, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	var d Definition

	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.File:
			d.Package = t.Name.Name
		case *ast.TypeSpec:
			e, ok := t.Type.(*ast.InterfaceType)
			if ok {
				d.interfaces = append(d.interfaces, visitInterface(e, t.Name.Name))
			}
		}
		return true
	})

	return &d, nil
}

func visitInterface(e *ast.InterfaceType, name string) Interface {
	iface := Interface{Name: name}
	for _, m := range e.Methods.List {
		f, ok := m.Type.(*ast.FuncType)
		if ok {
			method := Method{Name: m.Names[0].Name}

			for i, p := range f.Params.List {
				id := p.Type.(*ast.Ident)
				if len(p.Names) == 0 {
					method.Params = append(method.Params, Param{
						Name: "param" + strconv.Itoa(i+1),
						Type: id.Name,
					})
				} else {
					for _, n := range p.Names {
						method.Params = append(method.Params, Param{
							Name: n.Name,
							Type: id.Name,
						})
					}
				}
			}

			for i, p := range f.Results.List {
				id := p.Type.(*ast.Ident)
				if len(p.Names) == 0 {
					method.Results = append(method.Results, Param{
						Name: "res" + strconv.Itoa(i+1),
						Type: id.Name,
					})
				} else {
					for _, n := range p.Names {
						method.Results = append(method.Results, Param{
							Name: n.Name,
							Type: id.Name,
						})
					}
				}
			}
			iface.methods = append(iface.methods, method)
		}
	}
	return iface
}
