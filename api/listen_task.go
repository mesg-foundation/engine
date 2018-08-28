package api

// ListenTask listens tasks on service token.
func (a *API) ListenTask(token string) (*TaskListener, error) {
	l := newTaskListener(a)
	return l, l.listen(token)
}
