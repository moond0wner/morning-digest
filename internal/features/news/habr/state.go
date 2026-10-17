package habr

type step int

const (
	_ step = iota
	choiceSource
	choiceHub
	choiceLimit
)

type newsDialogState struct {
	Step       step   // 1 - выбор источника, 2 - выбор хаба, 3 - выбор лимита
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

func (h *handler) setState(userID int64, state *newsDialogState) {
	h.mtx.Lock()
	h.states[userID] = state
	h.mtx.Unlock()
}

func (h *handler) deleteState(userID int64) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	delete(h.states, userID)
}
