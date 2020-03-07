package examples

import (
	"github.com/ihaiker/wither/generate"
)

type Demo struct {
	generate.Wither

	Name string
	Age  int
	Sex  bool

	Test Test2

	Name2 *string
}

func (d *Demo) Get() {

}

type Test2 struct {
	Page2  string
	Limit2 int
}

type Test struct {
	Page  string
	Limit int
	*Test2
}
