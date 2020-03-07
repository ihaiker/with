package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func Find(filename string) (string, map[string]string, []Definition, error) {
	if fr, err := os.Open(filename); err != nil {
		return "", nil, nil, err
	} else if file, err := parser.ParseFile(token.NewFileSet(), filepath.Base(filename), fr, 0); err != nil {
		return "", nil, nil, err
	} else {
		ds := make([]Definition, 0)
		imports := map[string]string{}
		for _, spec := range file.Imports {
			if spec.Name == nil {
				imports[spec.Path.Value] = ""
			} else {
				imports[spec.Path.Value] = spec.Name.Name
			}
		}
		for _, decl := range file.Decls {
			switch it := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range it.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if d := find(ts); d.Wither || d.Getter {
							ds = append(ds, d)
						}
					}
				}
			}
		}
		return file.Name.Name, imports, ds, nil
	}

}

func find(ts *ast.TypeSpec) (definition Definition) {
	definition = Definition{
		Wither: false, Getter: false,
		Name:   ts.Name.Name,
		Fields: make([]DefinitionField, 0),
	}

	if st, ok := ts.Type.(*ast.StructType); ok {
		for _, field := range st.Fields.List {
			fieldName := ""
			if len(field.Names) > 0 {
				fieldName = field.Names[0].Name
			}

			switch se := field.Type.(type) {
			case *ast.Ident:
				definition.Fields = append(definition.Fields, DefinitionField{
					Name: fieldName, Type: se.Name,
				})
			case *ast.SelectorExpr:
				pkg := se.X.(*ast.Ident)
				typeName := pkg.Name + "." + se.Sel.Name
				switch typeName {
				case "generate.Wither":
					definition.Wither = true
				case "generate.Getter":
					definition.Getter = true
				default:
					definition.Fields = append(definition.Fields, DefinitionField{
						Name: fieldName, Type: typeName,
					})
				}
			case *ast.StarExpr:
				switch pkg := se.X.(type) {
				case *ast.Ident:
					definition.Fields = append(definition.Fields, DefinitionField{
						Name: fieldName, Type: "*" + pkg.Name,
					})
				case *ast.SelectorExpr:
					p := pkg.X.(*ast.Ident)
					definition.Fields = append(definition.Fields, DefinitionField{
						Name: fieldName, Type: "*" + p.Name + "." + pkg.Sel.Name,
					})
				}
			}
		}
	}
	return
}
