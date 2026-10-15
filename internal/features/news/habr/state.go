package habr

type step int

const (
	ChoiceSource step = iota
	ChoiceHub
	ChoiceLimit
)

type newsDialogState struct {
	Step       step   // 0 - выбор источника, 1 - выбор хаба, 2 - выбор лимита
	SourceType string // "base" или "hub"
	Hub        string // "go", "cpp", ...
	Limit      int    // 3, 5, 10
}

func (h *handler) getState(userID int64) (*newsDialogState, bool) {
	h.mtx.RLock()
	defer h.mtx.RUnlock()
	state, ok := h.states[userID]
	return state, ok
}

func (h *handler) deleteState(userID int64) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	delete(h.states, userID)
}
