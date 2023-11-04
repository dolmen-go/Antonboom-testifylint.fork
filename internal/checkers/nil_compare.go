package checkers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// NilCompare detects situations like
//
//	assert.Equal(t, value, nil)
//	assert.NotEqual(t, value, nil)
//
// and requires
//
//	assert.Nil(t, value)
//	assert.NotNil(t, value)
type NilCompare struct{}

// NewNilCompare constructs NilCompare checker.
func NewNilCompare() NilCompare { return NilCompare{} }
func (NilCompare) Name() string { return "nil-compare" }

func (checker NilCompare) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if len(call.Args) < 2 {
		return nil
	}

	survivingArg, ok := xorNil(call.Args[0], call.Args[1])
	if !ok {
		return nil
	}

	var proposedFn string

	switch call.Fn.Name {
	case "Equal", "Equalf", "EqualValues", "EqualValuesf", "Exactly", "Exactlyf":
		proposedFn = "Nil"
	case "NotEqual", "NotEqualf", "NotEqualValues", "NotEqualValuesf":
		proposedFn = "NotNil"
	default:
		return nil
	}

	return newUseFunctionDiagnostic(checker.Name(), call, proposedFn,
		newSuggestedFuncReplacement(call, proposedFn, analysis.TextEdit{
			Pos:     call.Args[0].Pos(),
			End:     call.Args[1].End(),
			NewText: analysisutil.NodeBytes(pass.Fset, survivingArg),
		}),
	)
}

func xorNil(first, second ast.Expr) (ast.Expr, bool) {
	a, b := isNil(first), isNil(second)
	if xor(a, b) {
		if a {
			return second, true
		}
		return first, true
	}
	return nil, false
}
