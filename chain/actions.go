package chain

type action struct {
	key      string
	setValue any
	delete   bool
}

// Actions remembers a series of actions to perform to the context after a step has executed Do.
type Actions struct {
	actions []action
}

// Set will set the element of the context to the value.
func (a *Actions) Set(key string, val any) {
	a.actions = append(a.actions, action{
		key:      key,
		setValue: val,
	})
}

// Delete will remove the element of the context.
func (a *Actions) Delete(key string) {
	a.actions = append(a.actions, action{
		key:    key,
		delete: true,
	})
}
