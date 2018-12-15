package stake

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/mesg-foundation/core/execution"
)

const (
	// HasStakeTask is the name of the task to call to check if the user staked tokens
	HasStakeTask    = "balanceOf"
	contractAddress = "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07"
)

// HasStakeInputs is the inputs mapping for the task
func HasStakeInputs(address string) map[string]interface{} {
	return map[string]interface{}{
		"address":         address,
		"contractAddress": contractAddress,
	}
}

// HasStakeOutputs is function that extracts the parameters
func HasStakeOutputs(e *execution.Execution) (hasStake bool, err error) {
	switch e.OutputKey {
	case "success":
		balance := e.OutputData["balance"].(string)
		val, err := strconv.ParseFloat(balance, 64)
		if err != nil {
			return false, err
		}
		return val > 0, nil
	case "error":
		return false, errors.New(e.OutputData["message"].(string))
	}
	return false, fmt.Errorf("stake: task hasStake has unknown output %s", e.OutputKey)
}
