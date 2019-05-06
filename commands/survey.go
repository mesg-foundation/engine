package commands

import survey "gopkg.in/AlecAivazis/survey.v1"

// askPass wraps survey.AskOne password.
func askPass(message string, response interface{}) error {
	return survey.AskOne(&survey.Password{Message: message}, response, survey.MinLength(1))
}

// askInput wraps survey.AskOne input.
func askInput(message string, response interface{}) error {
	return survey.AskOne(&survey.Input{Message: message}, response, survey.MinLength(1))
}

// askSelect wraps survey.AskOne select.
func askSelect(message string, options []string, response interface{}) error {
	return survey.AskOne(&survey.Select{
		Message: message,
		Options: options,
	}, response, nil)
}
