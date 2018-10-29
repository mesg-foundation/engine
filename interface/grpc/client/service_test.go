package client

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServices(t *testing.T) {
	wf := &Workflow{
		OnEvent:  &Event{ServiceID: "xxx"},
		OnResult: &Result{ServiceID: "yyy"},
		Execute:  &Task{ServiceID: "zzz"},
	}
	services := wf.services()
	require.Equal(t, len(services), 3)
	require.Equal(t, services[0], "xxx")
	require.Equal(t, services[1], "yyy")
	require.Equal(t, services[2], "zzz")
}

func TestServicesDuplicate(t *testing.T) {
	wf := &Workflow{
		OnEvent:  &Event{ServiceID: "xxx"},
		OnResult: &Result{ServiceID: "yyy"},
		Execute:  &Task{ServiceID: "xxx"},
	}
	services := wf.services()
	require.Equal(t, len(services), 2)
	require.Equal(t, services[0], "xxx")
	require.Equal(t, services[1], "yyy")
}

func TestIterateService(t *testing.T) {
	wf := &Workflow{
		OnEvent:  &Event{ServiceID: "xxx"},
		OnResult: &Result{ServiceID: "yyy"},
		Execute:  &Task{ServiceID: "zzz"},
	}
	cpt := 0
	err := iterateService(wf, func(ID string) error {
		cpt++
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, cpt, 3)
}

func TestIterateServiceWithError(t *testing.T) {
	wf := &Workflow{
		OnEvent:  &Event{ServiceID: "xxx"},
		OnResult: &Result{ServiceID: "yyy"},
		Execute:  &Task{ServiceID: "zzz"},
	}
	err := iterateService(wf, func(ID string) error {
		return errors.New("test error")
	})
	require.Error(t, err)
}
