package steps

import "chain"

var _ chain.Step = &unaryStep[any]{}

func NewUnaryStep[T any](srcKey string, do func(T, *chain.Actions) error) chain.Step {
	return &unaryStep[T]{
		srcKey: srcKey,
		fn:     do,
	}
}

func NewGTEStep(srcKey, tarKey string, threshold int) chain.Step {
	return NewUnaryStep(srcKey, func(inp int, actions *chain.Actions) error {
		result := inp > threshold
		actions.Set(tarKey, result)
		return nil
	})
}

func NewEQStep[T comparable](srcKey, tarKey string, compareTo T) chain.Step {
	return NewUnaryStep(srcKey, func(inp T, actions *chain.Actions) error {
		result := inp == compareTo
		actions.Set(tarKey, result)
		return nil
	})
}

type unaryStep[T any] struct {
	srcKey string
	fn     func(T, *chain.Actions) error
	x      T
}

// Do implements chain.Step.
func (b *unaryStep[T]) Do(actions *chain.Actions) ([]chain.Step, error) {
	return nil, b.fn(b.x, actions)
}

// Inputs implements chain.Step.
func (b *unaryStep[T]) Inputs() []chain.Input {
	return []chain.Input{
		chain.I(b.srcKey, &b.x),
	}
}
