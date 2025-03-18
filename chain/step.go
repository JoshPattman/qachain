package chain

type Step interface {
	Do(ctx *Context) ([]Step, error)
}

func Run(steps []Step, ctx *Context) error {
	for _, s := range steps {
		next, err := s.Do(ctx)
		if err != nil {
			return err
		}
		Run(next, ctx)
	}
	return nil
}
