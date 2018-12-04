package commands

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// confirmSurvey confirms a survey with Yes.
func confirmSurvey() func(args mock.Arguments) {
	return func(args mock.Arguments) {
		*args.Get(1).(*bool) = true
	}
}

// rejectSurvey reject a survey with No.
func rejectSurvey() func(args mock.Arguments) {
	return func(args mock.Arguments) {
		*args.Get(1).(*bool) = false
	}
}

func TestServiceDeleteCmdFlags(t *testing.T) {
	var (
		c     = newServiceDeleteCmd(nil, nil)
		flags = c.cmd.Flags()
	)

	// check defaults.
	require.False(t, c.yes)
	require.False(t, c.deleteAllServices)
	require.False(t, c.keepData)

	// check shortlands.
	require.Equal(t, flags.ShorthandLookup("y"), flags.Lookup("yes"))

	// check set flags.
	flags.Set("yes", "true")
	flags.Set("all", "true")
	flags.Set("keep-data", "true")
	require.True(t, c.yes)
	require.True(t, c.deleteAllServices)
	require.True(t, c.keepData)
}

func TestServiceDeleteWithoutFlagNorArgument(t *testing.T) {
	var (
		sm = newMockSurvey()
		me = newMockExecutor()
		c  = newServiceDeleteCmd(me, sm)
	)

	require.Contains(t, "at least one service id must be provided (or run with --all flag)", c.cmd.Execute().Error())

	sm.AssertExpectations(t)
	me.AssertExpectations(t)
}

func TestServiceDeleteWithAllFlagWithServiceConfirmation(t *testing.T) {
	var (
		sm = newMockSurvey()
		me = newMockExecutor()
		c  = newServiceDeleteCmd(me, sm)
	)

	c.cmd.Flags().Set("all", "true")

	sm.On("AskOne", &survey.Confirm{
		Message: "Are you sure to delete all services?",
		Default: false,
	}, mock.Anything, mock.Anything).Once().Return(nil).Run(confirmSurvey())

	sm.On("AskOne", &survey.Confirm{
		Message: "Do you want to remove service(s)' persistent data as well?",
		Default: false,
	}, mock.Anything, mock.Anything).Once().Return(nil)

	me.On("ServiceDeleteAll", false).Once().Return(nil)

	require.NoError(t, c.cmd.Execute())

	sm.AssertExpectations(t)
	me.AssertExpectations(t)
}

func TestServiceDeleteWithAllFlagWithServiceAndDataConfirmation(t *testing.T) {
	var (
		sm = newMockSurvey()
		me = newMockExecutor()
		c  = newServiceDeleteCmd(me, sm)
	)

	c.cmd.Flags().Set("all", "true")

	sm.On("AskOne", &survey.Confirm{
		Message: "Are you sure to delete all services?",
		Default: false,
	}, mock.Anything, mock.Anything).Once().Return(nil).Run(confirmSurvey())

	sm.On("AskOne", &survey.Confirm{
		Message: "Do you want to remove service(s)' persistent data as well?",
		Default: false,
	}, mock.Anything, mock.Anything).Once().Return(nil).Run(confirmSurvey())

	me.On("ServiceDeleteAll", true).Once().Return(nil)

	require.NoError(t, c.cmd.Execute())

	sm.AssertExpectations(t)
	me.AssertExpectations(t)
}

func TestServiceDeleteWithAllFlagWithoutDeleteServiceConfirmation(t *testing.T) {
	var (
		sm = newMockSurvey()
		me = newMockExecutor()
		c  = newServiceDeleteCmd(me, sm)
	)

	c.cmd.Flags().Set("all", "true")

	sm.On("AskOne", &survey.Confirm{
		Message: "Are you sure to delete all services?",
		Default: false,
	}, mock.Anything, mock.Anything).Once().Return(nil).Run(rejectSurvey())

	require.Contains(t, errConfirmationNeeded.Error(), c.cmd.Execute().Error())

	sm.AssertExpectations(t)
	me.AssertExpectations(t)
}

func TestServiceDeleteWithServiceID(t *testing.T) {
	var (
		serviceID = "1"
		sm        = newMockSurvey()
		me        = newMockExecutor()
		c         = newServiceDeleteCmd(me, sm)
	)

	c.cmd.SetArgs([]string{serviceID})

	sm.On("AskOne", &survey.Confirm{
		Message: "Are you sure to delete service(s)?",
		Default: false,
	}, mock.Anything, mock.Anything).Once().Return(nil).Run(confirmSurvey())

	sm.On("AskOne", &survey.Confirm{
		Message: "Do you want to remove service(s)' persistent data as well?",
		Default: false,
	}, mock.Anything, mock.Anything).Once().Return(nil)

	me.On("ServiceDelete", false, serviceID).Once().Return(nil)

	require.NoError(t, c.cmd.Execute())

	sm.AssertExpectations(t)
	me.AssertExpectations(t)
}

func TestServiceDeleteWithServiceIDAndYesFlag(t *testing.T) {
	var (
		serviceID = "1"
		sm        = newMockSurvey()
		me        = newMockExecutor()
		c         = newServiceDeleteCmd(me, sm)
	)

	c.cmd.SetArgs([]string{serviceID})
	c.cmd.Flags().Set("yes", "true")

	me.On("ServiceDelete", true, serviceID).Once().Return(nil)

	require.NoError(t, c.cmd.Execute())

	sm.AssertExpectations(t)
	me.AssertExpectations(t)
}

func TestServiceDeleteWithAllAndYesFlags(t *testing.T) {
	var (
		serviceID = "1"
		sm        = newMockSurvey()
		me        = newMockExecutor()
		c         = newServiceDeleteCmd(me, sm)
	)

	c.cmd.SetArgs([]string{serviceID})
	c.cmd.Flags().Set("all", "true")
	c.cmd.Flags().Set("yes", "true")

	me.On("ServiceDeleteAll", true).Once().Return(nil)

	require.NoError(t, c.cmd.Execute())

	sm.AssertExpectations(t)
	me.AssertExpectations(t)
}
