package resolver

const (
	// resolveFoundOutputKey is the found output key of resolve task.
	resolveFoundOutputKey = "found"

	// resolveFoundOutputKey is the not found output key of resolve task.
	resolveNotFoundOutputKey = "notFound"
)

// ResolveInputs is the input data of resolve task.
type ResolveInputs struct {
	// ServiceID is the service id that searched on peers.
	ServiceID string `json:"serviceID"`
}

// ResolveFoundOutput is the found output data of resolve task.
type ResolveFoundOutput struct {
	// Address is the IP address of core peer.
	Address string `json:"address"`

	// ServiceID is the service id of found service.
	ServiceID string `json:"serviceID"`
}

// ResolveNotFoundOutput is the not found output data of resolve task.
type ResolveNotFoundOutput struct {
	// ServiceID is the service id of service that couldn't found on peer network.
	ServiceID string `json:"serviceID"`
}

// Resolve is the task that return the address of a core that runs the desired service.
// If a core that run the desired service is found, a message is post on the outputFound chan.
// If no core that run the desired service is found, a message is post on the outputNotFound chan.
// If an error occurred during the execution of the task, a message is post on the outputError chan.
// If an error occurred before the execution of the task, it is return in the err variable.
// The chans are returned even if there is a possible memory leak if they are not listened. So make sure to listen to all of them.
func (r *Resolver) Resolve(inputs *ResolveInputs) (outputFound chan *ResolveFoundOutput, outputNotFound chan *ResolveNotFoundOutput, outputError chan *ErrorOutput, err error) {
	outputFound = make(chan *ResolveFoundOutput)
	outputNotFound = make(chan *ResolveNotFoundOutput)
	outputError = make(chan *ErrorOutput)
	return outputFound, outputNotFound, outputError, nil
}
