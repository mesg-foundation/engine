package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/sdk"
	workflowvm "github.com/mesg-foundation/core/workflow/vm"
)

func setupWorkflow(vm *workflowvm.VM, sdk *sdk.SDK) (closer func()) {
	var (
		execLis   = sdk.Execution.GetStream(nil)
		eventList = sdk.Event.GetStream(nil)
	)
	go func() {
		for exec := range execLis.C {
			go func(exec *execution.Execution) {
				fmt.Println(" exec", exec.TaskKey, exec.ServiceHash)
				if exec.Status == execution.Completed {
					fmt.Println(" exec completed", exec.TaskKey)
					vm.Process(workflowvm.Event{
						InstanceHash: exec.ServiceHash,
						ParentHash:   exec.ParentHash,
						Key:          "executionFinished",
						TaskKey:      exec.TaskKey,
						Data:         exec.Outputs,
						Secret:       getSecret(exec.Tags),
					})
				}
			}(exec)
		}
	}()
	go func() {
		for ev := range eventList.C {
			go func(ev *event.Event) {
				fmt.Println(" event", ev.Key, ev.Instance.Hash)
				vm.Process(workflowvm.Event{
					InstanceHash: ev.Instance.Hash,
					Key:          ev.Key,
					Data:         ev.Data.(map[string]interface{}),
				})
			}(ev)
		}
	}()
	go func() {
		for execReq := range vm.ExecuctionRequests {
			go func(execReq *workflowvm.Execution) {
				fmt.Println("execReq", execReq.InstanceHash, execReq.TaskKey, execReq.Inputs)
				if _, err := sdk.Execution.Execute(execReq.InstanceHash, execReq.TaskKey, execReq.Inputs,
					[]string{buildSecretTag(execReq.Secret)}); err != nil {
					log.Println("err while executing workflow task:", err)
				}
			}(execReq)
		}
	}()
	return func() {
		execLis.Close()
		eventList.Close()
		vm.Shutdown()
	}
}

const secretPrefix = "workflow,id:"

func getSecret(tags []string) string {
	for _, tag := range tags {
		if strings.HasPrefix(tag, secretPrefix) {
			return strings.TrimPrefix(tag, secretPrefix)
		}
	}
	return ""
}

func buildSecretTag(secret string) string {
	return secretPrefix + secret
}
