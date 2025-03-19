package chain

import "fmt"

type Input interface {
	Source() string
	Set(val any) error
}

type valInput[T any] struct {
	Src string
	Ptr *T
}

func I[T any](key string, ptr *T) Input {
	return &valInput[T]{
		Src: key,
		Ptr: ptr,
	}
}

func (v *valInput[T]) Source() string {
	return v.Src
}

func (v *valInput[T]) Set(val any) error {
	vt, ok := val.(T)
	if !ok {
		return fmt.Errorf("cannot convert %T to %T", val, vt)
	}
	*v.Ptr = vt
	return nil
}
