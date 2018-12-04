package commands

import survey "gopkg.in/AlecAivazis/survey.v1"

// Survey is a command line survey.
type Survey interface {
	AskOne(p survey.Prompt, response interface{}, v survey.Validator, opts ...survey.AskOpt) error
}

// defaultSurvey is default implementation of Survey.
type defaultSurvey struct{}

func newSurvey() *defaultSurvey {
	return &defaultSurvey{}
}

// AskOne prompts a Yes/No question.
func (*defaultSurvey) AskOne(p survey.Prompt, response interface{}, v survey.Validator,
	opts ...survey.AskOpt) error {
	return survey.AskOne(p, response, v, opts...)
}
