package chain

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
