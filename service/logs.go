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
	"io"

	"github.com/mesg-foundation/core/x/xstrings"
)

// Log holds log streams of dependency.
type Log struct {
	Dependency string
	Standard   io.ReadCloser
	Error      io.ReadCloser
}

// Logs gives service's logs and applies dependencies filter to filter logs.
// if dependencies has a length of zero all dependency logs will be provided.
func (s *Service) Logs(dependencies ...string) ([]*Log, error) {
	var (
		logs       []*Log
		isNoFilter = len(dependencies) == 0
	)
	for _, dep := range s.Dependencies {
		if isNoFilter || xstrings.SliceContains(dependencies, dep.Key) {
			rstd, rerr, err := dep.Logs()
			if err != nil {
				return nil, err
			}
			logs = append(logs, &Log{
				Dependency: dep.Key,
				Standard:   rstd,
				Error:      rerr,
			})
		}
	}
	return logs, nil
}
