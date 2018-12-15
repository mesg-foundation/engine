package core

import "github.com/mesg-foundation/core/config"

// ListenEvent listens events matches with eventFilter on serviceID.
func (s *Server) requireStake() error {
	c, err := config.Global()
	if err != nil {
		return err
	}
	k, err := c.GetKey()
	if err != nil {
		return err
	}
	return s.api.RequireStake(k.Address)
}
