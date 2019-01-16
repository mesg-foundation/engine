// Package ethwallet is a wrapper for system ethereum wallet service to map tasks data.
package ethwallet

// import (
// 	"fmt"

// 	"github.com/mesg-foundation/core/execution"
// )

// // Ethereum wallet tasks.
// const (
// 	AddPeersTask = "addPeers"
// 	ResolveTask  = "resolve"
// )

// // AddPeersInputs map add peer task inputs.
// func AddPeersInputs(addresses []string) map[string]interface{} {
// 	// TODO: this hack is not something that we should do but
// 	// it's needed because *parameterValidator is not able to identify
// 	// string slices for now.
// 	var addressesInterface []interface{}
// 	for _, address := range addresses {
// 		addressesInterface = append(addressesInterface, address)
// 	}
// 	return map[string]interface{}{"addresses": addressesInterface}
// }

// // AddPeersOutputs map add peer task outputs.
// func AddPeersOutputs(e *execution.Execution) error {
// 	switch e.OutputKey {
// 	case "success":
// 		return nil
// 	case "error":
// 		return fmt.Errorf("ethereum wallet: %s", e.OutputData["message"])
// 	}
// 	return fmt.Errorf("ethereum wallet: task add peers has unknown output %s", e.OutputKey)
// }

// // ResolveInputs map resolve task inputs.
// func ResolveInputs(serviceID string) map[string]interface{} {
// 	return map[string]interface{}{"serviceID": serviceID}
// }

// // ResolveOutputs map resolve task outputs.
// func ResolveOutputs(e *execution.Execution) (peerAddress string, err error) {
// 	switch e.OutputKey {
// 	case "found":
// 		return e.OutputData["address"].(string), nil
// 	case "notFound":
// 		return "", fmt.Errorf("ethereum wallet: peer address could not be found for %s service", e.OutputData["serviceID"])
// 	case "error":
// 		return "", fmt.Errorf("ethereum wallet: %s", e.OutputData["message"])
// 	}
// 	return "", fmt.Errorf("ethereum wallet: task resolve has unknown output %s", e.OutputKey)
// }
