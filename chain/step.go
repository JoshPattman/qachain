package chain

import "fmt"

// Step is a specific piece of code that may execute during running a chain.
type Step interface {
	// Inputs defines which inputs this step requires to run.
	Inputs() []Input
	// Do performs this step, making any changes to the state through the actions.
	// The changes only take effect after do has completed.
	Do(actions *Actions) ([]Step, error)
}

// Run the provided steps and their follow up steps recursively.
// Steps are run in a depth-first manner (all of a step's follow ups will be run before its next sibling).
// The context both provides the inputs to the step, and collects the outputs from them.
func Run(steps []Step, ctx *Context) error {
	for _, s := range steps {
		inputs := s.Inputs()
		for _, i := range inputs {
			val, ok := ctx.vars[i.Source()]
			if !ok {
				return fmt.Errorf("step required a key that did not exist")
			}
			err := i.Set(val)
			if err != nil {
				return err
			}
		}
		actions := &Actions{}
		next, err := s.Do(actions)
		if err != nil {
			return err
		}
		for _, ac := range actions.actions {
			if ac.delete {
				delete(ctx.vars, ac.key)
			} else {
				ctx.vars[ac.key] = ac.setValue
			}
		}
		if err := Run(next, ctx); err != nil {
			return err
		}
	}
	return nil
}
