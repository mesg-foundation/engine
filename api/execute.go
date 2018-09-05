package api

// ExecuteTask executes a task tasKey with inputData and tags for service serviceID.
func (a *API) ExecuteTask(serviceID, taskKey string, inputData map[string]interface{},
	tags []string) (executionID string, err error) {
	return newTaskExecutor(a).Execute(serviceID, taskKey, inputData, tags)
}
