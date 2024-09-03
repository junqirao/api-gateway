package program

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	registry "github.com/junqirao/simple-registry"
)

type test struct {
}

func (t test) Print() bool {
	fmt.Println("test-ok")
	return true
}

func TestProgram(t *testing.T) {
	cfg := registry.Config{}
	err := g.Cfg().MustGet(context.Background(), "registry").Struct(&cfg)
	if err != nil {
		panic(err)
		return
	}
	err = registry.Init(context.Background(), cfg)
	if err != nil {
		panic(err)
		return
	}

	Init(context.Background())

	program, err := NewProgram("test", `
		1+1;
		function1();
		function2("aaa");
		setGlobalVariable("test_value", 123);
		print(global.test_value);
		test_struct.Print();
		print(val!=2);
		`)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = program.Exec(context.Background(), map[string]interface{}{
		"function1": func() bool {
			return true
		},
		"function2": func(a string) string {
			return a
		},
		"print": func(s interface{}) bool {
			fmt.Printf("%v\n", s)
			return true
		},
		"test_struct": &test{},
		"val":         1,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
}
