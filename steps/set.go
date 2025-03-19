package steps

import "chain"

var _ chain.Step = &setStep{}

type setStep struct {
	key   string
	value any
}

// NewSetStep creates a step that always sets the provided variable.
func NewSetStep(key string, value any) chain.Step {
	return &setStep{
		key:   key,
		value: value,
	}
}

// Do implements Step.
func (s *setStep) Do(actions *chain.Actions) ([]chain.Step, error) {
	actions.Set(s.key, s.value)
	return nil, nil
}

// Inputs implements Step.
func (f *setStep) Inputs() []chain.Input {
	return []chain.Input{}
}
