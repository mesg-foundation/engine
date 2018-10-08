package resolver

const (
	// addPeersSuccessOutputKey is the success output key of add peers task.
	addPeersSuccessOutputKey = "success"
)

// AddPeersInputs is the input data of add peers task.
type AddPeersInputs struct {
	// Addresses are the addresses of core peers.
	Addresses []string `json:"addresses"`
}

// AddPeersSuccessOutput is the success output data of add peers task.
type AddPeersSuccessOutput struct {
	// Addresses are the addresses of core peers.
	Addresses []string `json:"addresses"`
}

// AddPeers is the task that actually new add peers to the service.
// If the task execute successfully, a message is post on the outputSuccess chan.
// If an error occurred during the execution of the task, a message is post on the outputError chan.
// If an error occurred before the execution of the task, it is return in the err variable.
// The chans are returned even if there is a possible memory leak if they are not listened. So make sure to listen to all of them.
func (r *Resolver) AddPeers(inputs *AddPeersInputs) (outputSuccess chan *AddPeersSuccessOutput, outputError chan *ErrorOutput, err error) {
	outputSuccess = make(chan *AddPeersSuccessOutput)
	outputError = make(chan *ErrorOutput)
	return outputSuccess, outputError, nil
}
