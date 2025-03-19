package steps

import "chain"

var _ chain.Step = &conditionalStep{}

type conditionalStep struct {
	key     string
	ifTrue  []chain.Step
	ifFalse []chain.Step
	value   bool
}

// NewConditionalStep creates a new step that runs ine of it's child options depending on the boolean variable.
func NewConditionalStep(key string, ifTrue, ifFalse []chain.Step) chain.Step {
	return &conditionalStep{
		key:     key,
		ifTrue:  ifTrue,
		ifFalse: ifFalse,
	}
}

// Do implements Step.
func (f *conditionalStep) Do(actions *chain.Actions) ([]chain.Step, error) {
	if f.value {
		return f.ifTrue, nil
	} else {
		return f.ifFalse, nil
	}
}

// Inputs implements Step.
func (f *conditionalStep) Inputs() []chain.Input {
	return []chain.Input{
		chain.I(f.key, &f.value),
	}
}
