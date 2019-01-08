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

package commands

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// logFormatValue represents log format flag value.
type logFormatValue string

func (v *logFormatValue) Set(value string) error {
	if value != "text" && value != "json" {
		return fmt.Errorf("%s is not valid log format", value)
	}
	*v = logFormatValue(value)
	return nil
}
func (v *logFormatValue) Type() string   { return "string" }
func (v *logFormatValue) String() string { return string(*v) }

// logLevelValue represents log level flag value.
type logLevelValue string

func (v *logLevelValue) Set(value string) error {
	if _, err := logrus.ParseLevel(value); err != nil {
		return fmt.Errorf("%s is not valid log level", value)
	}
	*v = logLevelValue(value)
	return nil
}
func (v *logLevelValue) Type() string   { return "string" }
func (v *logLevelValue) String() string { return string(*v) }
