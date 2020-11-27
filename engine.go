package dsl

import (
	"regexp"
	"sync"

	"github.com/Knetic/govaluate"
)

var (
	defaultEngine *Engine
	RegexStore    = make(map[string]*regexp.Regexp)
	regexmux      sync.RWMutex
)

type Engine struct {
	HelperFunctions map[string]govaluate.ExpressionFunction
	ExpressionStore map[string]*govaluate.EvaluableExpression
	exprmux         sync.RWMutex
}

func NewEngine() (*Engine, error) {
	engine := &Engine{
		HelperFunctions: DefaultHelperFunctions,
		ExpressionStore: make(map[string]*govaluate.EvaluableExpression),
	}
	return engine, nil
}

func (e *Engine) EvalExpr(expr string, vars map[string]interface{}) (interface{}, error) {
	e.exprmux.Lock()
	defer e.exprmux.Unlock()
	compiled, err := govaluate.NewEvaluableExpressionWithFunctions(expr, e.HelperFunctions)
	if err != nil {
		return nil, err
	}

	e.ExpressionStore[expr] = compiled

	return compiled.Evaluate(vars)
}

func (e *Engine) EvalExprFromCache(expr string, vars map[string]interface{}) (interface{}, error) {
	compiled, ok := e.ExpressionStore[expr]
	if !ok {
		return e.EvalExpr(expr, vars)
	}

	return compiled.Evaluate(vars)
}

func EvalExpr(expr string, vars map[string]interface{}) (interface{}, error) {
	if defaultEngine == nil {
		var err error
		defaultEngine, err = NewEngine()
		if err != nil {
			return nil, err
		}
	}

	return defaultEngine.EvalExprFromCache(expr, vars)
}

func Regex(regxp string) (*regexp.Regexp, error) {
	if compiled, ok := RegexStore[regxp]; ok {
		return compiled, nil
	}

	regexmux.Lock()
	defer regexmux.Unlock()
	compiled, err := regexp.Compile(regxp)
	if err != nil {
		return nil, err
	}
	RegexStore[regxp] = compiled

	return compiled, nil
}
