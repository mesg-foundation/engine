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

package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Init initializes default logger. It panics on invalid format or level.
func Init(format, level string) {
	switch format {
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		panic(fmt.Sprintf("log: %s is not a valid format", format))
	}

	l, err := logrus.ParseLevel(level)
	if err != nil {
		panic(fmt.Sprintf("log: %s is not a valid level", level))
	}
	logrus.SetLevel(l)
}
