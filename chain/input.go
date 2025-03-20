package chain

import (
	"fmt"
)

// Input is a type that knows a source key from the context, and can attempt to set itself with that value.
type Input interface {
	Source() string
	Set(val any) error
}

type valInput[T any] struct {
	Src string
	Ptr *T
}

// I creates a new Input.
func I[T any](key string, ptr *T) Input {
	return &valInput[T]{
		Src: key,
		Ptr: ptr,
	}
}

// Source implements Input.
func (v *valInput[T]) Source() string {
	return v.Src
}

// Set implements Input.
func (v *valInput[T]) Set(val any) error {
	vt, ok := val.(T)
	if !ok {
		return fmt.Errorf("cannot convert %T to %T", val, vt)
	}
	*v.Ptr = vt
	return nil
}

func extractInput(i Input, ctx *Context) error {
	val, ok := ctx.vars[i.Source()]
	if !ok {
		return fmt.Errorf("step required a key that did not exist")
	}
	err := i.Set(val)
	if err != nil {
		return err
	}
	return nil
}
