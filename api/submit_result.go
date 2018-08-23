package api

// SubmitResult submits results for executionID.
func (a *API) SubmitResult(executionID, outputKey string, outputData map[string]interface{}) error {
	return newResultSubmitter(a).Submit(executionID, outputKey, outputData)
}
