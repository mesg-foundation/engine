package api

// ListenResultFilter is a filter func for filtering results.
type ListenResultFilter func(*ResultListener)

// ListenResultTaskFilter returns a taskKey filter.
func ListenResultTaskFilter(taskKey string) ListenResultFilter {
	return func(ln *ResultListener) {
		ln.taskKey = taskKey
	}
}

// ListenResultOutputFilter returns a outputKey filter.
func ListenResultOutputFilter(outputKey string) ListenResultFilter {
	return func(ln *ResultListener) {
		ln.outputKey = outputKey
	}
}

// ListenResultTagFilters returns a tags filter.
func ListenResultTagFilters(tags []string) ListenResultFilter {
	return func(ln *ResultListener) {
		ln.tagFilters = tags
	}
}

// ListenResult listens results matches with fiters on serviceID.
func (a *API) ListenResult(serviceID string, filters ...ListenResultFilter) (*ResultListener, error) {
	l := newResultListener(a, filters...)
	return l, l.listen(serviceID)
}
