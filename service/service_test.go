package service

import (
	"testing"

	"github.com/gogo/protobuf/types"
	"github.com/mesg-foundation/engine/protobuf/convert"
	"github.com/stretchr/testify/require"
)

func TestGetEvent(t *testing.T) {
	var (
		eventKey = "1"
		s        = &Service{
			Events: []*Service_Event{
				{Key: eventKey},
			},
		}
	)

	e, err := s.GetEvent(eventKey)
	require.NoError(t, err)
	require.Equal(t, eventKey, e.Key)
}

func TestGetEventNonExistent(t *testing.T) {
	var (
		serviceName = "1"
		eventKey    = "2"
		s           = &Service{
			Name: serviceName,
			Events: []*Service_Event{
				{Key: "3"},
			},
		}
	)

	e, err := s.GetEvent(eventKey)
	require.Zero(t, e)
	require.Equal(t, &EventNotFoundError{
		EventKey:    eventKey,
		ServiceName: serviceName,
	}, err)
}

func TestEventValidateAndRequireData(t *testing.T) {
	var (
		serviceName    = "1"
		eventKey       = "2"
		paramKey       = "3"
		validEventData = map[string]interface{}{
			paramKey: "4",
		}
		inValidEventData = map[string]interface{}{
			paramKey: 4,
		}
		s = &Service{
			Name: serviceName,
			Events: []*Service_Event{
				{
					Key: eventKey,
					Data: []*Service_Parameter{
						{
							Key:  paramKey,
							Type: "String",
						},
					},
				},
			},
		}
	)

	warnings, err := s.ValidateEventData(eventKey, validEventData)
	require.NoError(t, err)
	require.Len(t, warnings, 0)

	data := &types.Struct{}
	convert.Unmarshal(validEventData, data)
	require.NoError(t, err)

	err = s.RequireEventData(eventKey, data)
	require.NoError(t, err)

	warnings, err = s.ValidateEventData(eventKey, inValidEventData)
	require.NoError(t, err)
	require.Len(t, warnings, 1)
	require.Equal(t, paramKey, warnings[0].Key)

	convert.Unmarshal(inValidEventData, data)
	require.NoError(t, err)

	err = s.RequireEventData(eventKey, data)
	require.Equal(t, &InvalidEventDataError{
		EventKey:    eventKey,
		ServiceName: serviceName,
		Warnings:    warnings,
	}, err)
}
