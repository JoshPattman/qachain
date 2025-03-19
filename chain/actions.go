package chain

type action struct {
	key      string
	setValue any
	delete   bool
}
type Actions struct {
	actions []action
}

func (a *Actions) Set(key string, val any) {
	a.actions = append(a.actions, action{
		key:      key,
		setValue: val,
	})
}

func (a *Actions) Delete(key string) {
	a.actions = append(a.actions, action{
		key:    key,
		delete: true,
	})
}
