package chain

import (
	"fmt"
	"iter"
)

type Context struct {
	document string
	vars     map[string]any
}

func NewContext(doc string) *Context {
	return &Context{
		document: doc,
		vars:     make(map[string]any),
	}
}

func (ctx *Context) Document() string {
	return ctx.document
}

func (ctx *Context) Values() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for key, val := range ctx.vars {
			if !yield(key, val) {
				return
			}
		}
	}
}

func Get[T any](ctx *Context, key string) (T, error) {
	val, ok := ctx.vars[key]
	if !ok {
		return *new(T), fmt.Errorf("key '%s' did not exist in context", key)
	}
	valT, ok := val.(T)
	if !ok {
		return *new(T), fmt.Errorf("key '%s' had type %T but wanted %T", key, val, *new(T))
	}
	return valT, nil
}

func Set[T any](ctx *Context, key string, val T) error {
	ctx.vars[key] = val
	return nil
}

type Step interface {
	Do(ctx *Context) ([]Step, error)
}

func Run(steps []Step, ctx *Context) error {
	for _, s := range steps {
		next, err := s.Do(ctx)
		if err != nil {
			return err
		}
		Run(next, ctx)
	}
	return nil
}

var _ Step = &conditionalStep{}

type conditionalStep struct {
	fn      func(ctx *Context) (bool, error)
	ifTrue  []Step
	ifFalse []Step
}

func NewConditionalStep(condition func(ctx *Context) (bool, error), ifTrue, ifFalse []Step) Step {
	return &conditionalStep{
		fn:      condition,
		ifTrue:  ifTrue,
		ifFalse: ifFalse,
	}
}

// Do implements Step.
func (f *conditionalStep) Do(ctx *Context) ([]Step, error) {
	cond, err := f.fn(ctx)
	if err != nil {
		return nil, err
	}
	if cond {
		return f.ifTrue, nil
	} else {
		return f.ifFalse, nil
	}
}

var _ Step = &setStep{}

type setStep struct {
	key   string
	value any
}

func NewSetStep(key string, value any) Step {
	return &setStep{
		key:   key,
		value: value,
	}
}

// Do implements Step.
func (s *setStep) Do(ctx *Context) ([]Step, error) {
	err := Set(ctx, s.key, s.value)
	return nil, err
}
