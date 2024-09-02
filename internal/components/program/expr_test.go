package program

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/expr-lang/expr"
)

func TestExpr(t *testing.T) {
	env := map[string]interface{}{
		"greet":   "Hello, %v!",
		"names":   []string{"world", "you"},
		"sprintf": fmt.Sprintf,
		"do": func(res ...any) string {
			for _, v := range res {
				fmt.Println(reflect.TypeOf(v).Kind().String())
			}
			return "ok"
		},
		"echo": func() bool {
			fmt.Println("hello")
			return true
		},
	}

	code :=
		`echo();
sprintf(greet, names[0]);
1+1;let a=1;a+1;`
	split := strings.Split(strings.ReplaceAll(code, "\n", ""), ";")
	if split[len(split)-1] == "" {
		split = split[:len(split)-1]
	}
	program, err := expr.Compile(fmt.Sprintf("do(%s)", strings.Join(split, ",")), expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}
