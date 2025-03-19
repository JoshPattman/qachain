package chain

import "fmt"

type Step interface {
	Inputs() []Input
	Do(actions *Actions) ([]Step, error)
}

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
		Run(next, ctx)
	}
	return nil
}
