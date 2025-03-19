package chain

var _ Step = &conditionalStep{}

type conditionalStep struct {
	key     string
	ifTrue  []Step
	ifFalse []Step
	value   bool
}

// NewConditionalStep creates a new step that runs ine of it's child options depending on the boolean variable.
func NewConditionalStep(key string, ifTrue, ifFalse []Step) Step {
	return &conditionalStep{
		key:     key,
		ifTrue:  ifTrue,
		ifFalse: ifFalse,
	}
}

// Do implements Step.
func (f *conditionalStep) Do(actions *Actions) ([]Step, error) {
	if f.value {
		return f.ifTrue, nil
	} else {
		return f.ifFalse, nil
	}
}

// Inputs implements Step.
func (f *conditionalStep) Inputs() []Input {
	return []Input{
		I(f.key, &f.value),
	}
}

var _ Step = &setStep{}

type setStep struct {
	key   string
	value any
}

// NewSetStep creates a step that always sets the provided variable.
func NewSetStep(key string, value any) Step {
	return &setStep{
		key:   key,
		value: value,
	}
}

// Do implements Step.
func (s *setStep) Do(actions *Actions) ([]Step, error) {
	actions.Set(s.key, s.value)
	return nil, nil
}

// Inputs implements Step.
func (f *setStep) Inputs() []Input {
	return []Input{}
}
