package chain

import (
	"fmt"
	"iter"
	"maps"
)

// Context is a wrapper around a set of key-value pairs
// It is modified indirectly when steps are run, and provides the input data to further steps.
type Context struct {
	vars map[string]any
}

// NewContext creates a new context with the given initial values.
// Themap will be copied so the underlying map is not shared.
func NewContext(initialValues map[string]any) *Context {
	var vals map[string]any
	if initialValues == nil {
		vals = make(map[string]any)
	} else {
		vals = maps.Clone(initialValues)
	}
	return &Context{
		vars: vals,
	}
}

// Values iterates over each key-value pair in the map.
func (ctx *Context) Values() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for key, val := range ctx.vars {
			if !yield(key, val) {
				return
			}
		}
	}
}

// Get extracts a specific typed value from the context.
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
