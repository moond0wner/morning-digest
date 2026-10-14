package habr

type NewsDialogState struct {
	Step       step   // 0 - выбор источника, 1 - выбор хаба, 2 - выбор лимита
	SourceType string // "base" или "hub"
	Hub        string // "go", "cpp", ...
	Limit      int    // 3, 5, 10
}

type step int

const (
	ChoiceSource step = iota
	ChoiceHub
	ChoiceLimit
)

func (h *Handler) getState(userID int64) (*NewsDialogState, bool) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	state, ok := h.states[userID]
	return state, ok
}

func (h *Handler) deleteState(userID int64) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	delete(h.states, userID)
}
