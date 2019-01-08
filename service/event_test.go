// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetEvent(t *testing.T) {
	var (
		eventKey = "1"
		s, _     = FromService(&Service{
			Events: []*Event{
				{Key: eventKey},
			},
		})
	)
	e, err := s.GetEvent(eventKey)
	require.NoError(t, err)
	require.Equal(t, eventKey, e.Key)
}

func TestGetEventNonExistent(t *testing.T) {
	var (
		serviceName = "1"
		eventKey    = "2"
		s, _        = FromService(&Service{
			Name: serviceName,
			Events: []*Event{
				{Key: "3"},
			},
		})
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
		s, _ = FromService(&Service{
			Name: serviceName,
			Events: []*Event{
				{
					Key: eventKey,
					Data: []*Parameter{
						{
							Key:  paramKey,
							Type: "String",
						},
					},
				},
			},
		})
	)

	e, _ := s.GetEvent(eventKey)

	warnings := e.ValidateData(validEventData)
	require.Len(t, warnings, 0)

	err := e.RequireData(validEventData)
	require.NoError(t, err)

	warnings = e.ValidateData(inValidEventData)
	require.Len(t, warnings, 1)
	require.Equal(t, paramKey, warnings[0].Key)

	err = e.RequireData(inValidEventData)
	require.Equal(t, &InvalidEventDataError{
		EventKey:    eventKey,
		ServiceName: serviceName,
		Warnings:    warnings,
	}, err)
}
