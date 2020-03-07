package parse

import "strings"

type DefinitionField struct {
	Name string
	Type string
}

func (df *DefinitionField) GetName() string {
	name := df.Name
	if df.Name == "" {
		t := df.Type
		if t[0] == '*' {
			t = t[1:]
		}
		if idx := strings.LastIndex(t, "."); idx != -1 {
			t = t[idx+1:]
		}
		name = t
	}
	return name
}

func (df *DefinitionField) IsPrivate() bool {
	name := df.GetName()
	return 'a' <= name[0] && name[0] <= 'z'
}

type Definition struct {
	Wither bool
	Getter bool
	Name   string
	Fields []DefinitionField
}
