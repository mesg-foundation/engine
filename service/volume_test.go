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

func TestDeleteVolumes(t *testing.T) {
	var (
		dependencyKey1 = "1"
		dependencyKey2 = "2"
		volumeA        = "a"
		volumeB        = "b"
		s, mc          = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestCreateVolumes",
			Dependencies: []*Dependency{
				{
					Key:     dependencyKey1,
					Image:   "1",
					Volumes: []string{volumeA, volumeB},
				},
				{
					Key:         dependencyKey2,
					Image:       "1",
					VolumesFrom: []string{dependencyKey1},
				},
			},
		})
	)

	var (
		d1, _    = s.getDependency(dependencyKey1)
		volumes1 = d1.extractVolumes()
	)

	mc.On("DeleteVolume", volumes1[0].Source).Once().Return(nil)
	mc.On("DeleteVolume", volumes1[1].Source).Once().Return(nil)

	require.NoError(t, s.DeleteVolumes())

	mc.AssertExpectations(t)
}
