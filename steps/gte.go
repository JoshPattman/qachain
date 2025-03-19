package steps

import "chain"

var _ chain.Step = &gte{}

func NewGTEStep(srcKey, tarKey string, threshold int) chain.Step {
	return &gte{
		srcKey:    srcKey,
		tarKey:    tarKey,
		threshold: threshold,
	}
}

type gte struct {
	srcKey    string
	tarKey    string
	threshold int
	srvVal    int
}

// Do implements chain.Step.
func (t *gte) Do(actions *chain.Actions) ([]chain.Step, error) {
	val := false
	if t.srvVal > t.threshold {
		val = true
	}
	actions.Set(t.tarKey, val)
	return nil, nil
}

// Inputs implements chain.Step.
func (t *gte) Inputs() []chain.Input {
	return []chain.Input{
		chain.I(t.srcKey, &t.srvVal),
	}
}
