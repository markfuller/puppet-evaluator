package eval

import (
	"github.com/lyraproj/issue/issue"
	"github.com/lyraproj/puppet-evaluator/threadlocal"
	"github.com/lyraproj/puppet-parser/parser"
)

// An Evaluator is responsible for evaluating an Abstract Syntax Tree, typically produced by
// the parser. An implementation must be re-entrant.
type Evaluator interface {
	Context

	// Eval should be considered internal. The only reason it is public is to allow
	// the evaluator to be extended. This is subject to change. Don't use
	Eval(expression parser.Expression) Value
}

type ParserExtension interface {
	Evaluate(e Evaluator) Value
}

// Go calls the given function in a new go routine. The CurrentContext is forked and becomes
// the CurrentContext for that routine.
func Go(f ContextDoer) {
	Fork(CurrentContext(), f)
}

// Fork calls the given function in a new go routine. The given context is forked and becomes
// the CurrentContext for that routine.
func Fork(c Context, doer ContextDoer) {
	go func() {
		defer threadlocal.Cleanup()
		threadlocal.Init()
		cf := c.Fork()
		threadlocal.Set(PuppetContextKey, cf)
		doer(cf)
	}()
}

func LogWarning(issueCode issue.Code, args issue.H) {
	CurrentContext().Logger().LogIssue(Warning(issueCode, args))
}

// Error creates a Reported with the given issue code, location from stack top, and arguments
// Typical use is to panic with the returned value
var Error func(issueCode issue.Code, args issue.H) issue.Reported

// Error2 creates a Reported with the given issue code, location from stack top, and arguments
// Typical use is to panic with the returned value
var Error2 func(location issue.Location, issueCode issue.Code, args issue.H) issue.Reported

// Warning creates a Reported with the given issue code, location from stack top, and arguments
// and logs it on the currently active logger
var Warning func(issueCode issue.Code, args issue.H) issue.Reported
