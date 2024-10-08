package program

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	registry "github.com/junqirao/simple-registry"
)

const (
	lineSeparator              = ";"
	envKeyNewResultWrapper     = "newResultWrapper"
	envKeyExprMultilineWrapper = "exprMultilineWrapper"
	envKeyGlobalVariable       = "global"
	envKeySetGlobalVariable    = "set_global"
	envKeyLogger               = "logger"
	envKeyRequest              = "request"
	envKeyResponse             = "response"
	envKeyUpstream             = "upstream"
	envKeyCtx                  = "ctx"
	envKeyTerminateIf          = "terminate_if"
	envKeyIPGEO                = "ipgeo"
	envKeyJWT                  = "jwt"
)

type (
	// Program expression program
	Program struct {
		Name string
		p    *vm.Program
		md5  string
	}
	// Programs sets
	Programs struct {
		mu sync.RWMutex
		ps []*Program
	}
	resultWrapper struct {
		v    any
		expr string
	}
)

func newResultWrapper(v any, expr string) *resultWrapper {
	return &resultWrapper{v: v, expr: expr}
}

func (r *resultWrapper) Ok() (b bool, reason string) {
	if r.v != nil {
		if v, ok := r.v.(bool); ok {
			return v, r.expr
		}
		if v, ok := r.v.(error); ok && v != nil {
			return false, v.Error()
		}
	}
	return true, r.expr
}

// NewProgram create a new program
// statements: separated by ';' e.g.
//
//	1+1;
//	function1();
//	function2();
func NewProgram(name, statements string) (*Program, error) {
	p := &Program{
		Name: name,
	}
	return p, p.build(statements)
}

func (p *Program) Exec(_ context.Context, env ...map[string]interface{}) error {
	// prepare
	var e map[string]interface{}
	if len(env) > 0 && env[0] != nil {
		e = env[0]
	} else {
		e = make(map[string]interface{})
	}

	// run
	res, err := expr.Run(p.p, e)
	if err != nil {
		return err
	}

	// check result
	if s, ok := res.(string); ok && s != "" {
		return errors.New(s)
	}
	return nil
}

func (p *Program) Update(statements string) error {
	return p.build(statements)
}

func (p *Program) build(statements string) (err error) {
	md5 := gmd5.MustEncryptString(statements)
	if p.md5 == md5 {
		return
	}

	p.md5 = md5

	// pre process support multiline statements,
	// till expr-lang supports https://github.com/expr-lang/expr/issues/697

	var parts []string
	for _, s := range breakLine(statements) {
		if s != "" {
			parts = append(parts, fmt.Sprintf("%s(%s,`%s`)", envKeyNewResultWrapper, s, s))
		}
	}
	statements = fmt.Sprintf(
		"exprMultilineWrapper(%s)",
		strings.Join(parts, ","))
	program, err := expr.Compile(statements)
	if err != nil {
		return
	}
	p.p = program
	return
}

func buildTest(ctx context.Context, statements string) error {
	env := BuildEnvFromRequest(ctx, nil, registry.Instance{
		Id:          "program_build_test_id",
		Host:        "program_build_test_host",
		HostName:    "program_build_test_hostname",
		Port:        0,
		ServiceName: "program_build_test_service_name",
		Meta:        make(map[string]interface{}),
	})
	for i, line := range breakLine(statements) {
		p, err := expr.Compile(line)
		if err != nil {
			return fmt.Errorf("syntax error at line %d: %v", i+1, err)
		}
		if _, err = expr.Run(p, env); err != nil {
			return fmt.Errorf("compile error at line %d: %v", i+1, err)
		}
	}
	return nil
}

func breakLine(statements string) []string {
	// replace all line wrappers
	statements = strings.ReplaceAll(statements, "\r", "")
	statements = strings.ReplaceAll(statements, "\n", "")
	statements = strings.ReplaceAll(statements, "\t", "")
	statements = strings.TrimSpace(statements)
	// split by lineSeparator
	var lines []string
	for _, line := range strings.Split(statements, lineSeparator) {
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func (p *Programs) Create(info *Info) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	name := info.Name
	program, err := NewProgram(name, info.TryDecode(context.Background()))
	if err != nil {
		return err
	}

	for i, prog := range p.ps {
		if prog.Name == name {
			p.ps[i] = program
			return nil
		}
	}

	p.ps = append(p.ps, program)

	g.Log().Infof(context.Background(), "service %s build programs success: name=%s, count=%v", info.ServiceName, info.Name, len(p.ps))
	return nil
}

func (p *Programs) Exec(ctx context.Context, env map[string]interface{}) (lastExec string, err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, program := range p.ps {
		lastExec = program.Name
		if err = program.Exec(ctx, env); err != nil {
			return
		}
	}
	return
}

func (p *Programs) Delete(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, program := range p.ps {
		if program.Name == name {
			p.ps = append(p.ps[:i], p.ps[i+1:]...)
			return
		}
	}
}
