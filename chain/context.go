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
