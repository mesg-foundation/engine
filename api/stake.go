package api

import (
	"fmt"

	"github.com/mesg-foundation/core/systemservices/stake"
	"github.com/sirupsen/logrus"
)

// RequireStake returns an error if the given address doesn't have any MESG token
func (a *API) RequireStake(address string) error {
	logrus.WithField("address", address).Debug("Checking stake")
	serviceID, err := a.systemservices.StakeServiceID()
	if err != nil {
		return err
	}
	e, err := a.ExecuteAndListen(serviceID, stake.HasStakeTask, stake.HasStakeInputs(address))
	if err != nil {
		return err
	}
	res, err := stake.HasStakeOutputs(e)
	if err != nil {
		return err
	}
	if !res {
		return fmt.Errorf("MESG tokens are required in order to use this function")
	}
	return nil
}
