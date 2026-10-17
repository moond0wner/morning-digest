package todo

type step int

const (
	_ step = iota
	choiceButton
	inputTitle
	inputDescription
	choiceTaskForDone
)

type taskDialogState struct {
	Step step
	Task *Task
}

func (h *handler) getState(userID int64) (*taskDialogState, bool) {
	h.mtx.RLock()
	defer h.mtx.RUnlock()
	state, ok := h.states[userID]
	return state, ok
}

func (h *handler) setState(userID int64, state *taskDialogState) {
	h.mtx.Lock()
	h.states[userID] = state
	h.mtx.Unlock()
}

func (h *handler) deleteState(userID int64) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	delete(h.states, userID)
}
