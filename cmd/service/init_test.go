package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTemplate(t *testing.T) {
	templates, err := getTemplates(templatesURL)
	require.Nil(t, err)
	require.True(t, len(templates) > 0)
}

func TestGetTemplateBadURL(t *testing.T) {
	_, err := getTemplates("...")
	require.NotNil(t, err)
}

func TestTemplatesToOption(t *testing.T) {
	templates := []*templateStruct{
		{Name: "xxx", URL: "https://..."},
		{Name: "yyy", URL: "https://..."},
	}
	options := templatesToOption(templates)
	require.Equal(t, len(options), len(templates)+2)
	require.Equal(t, options[0], "xxx (https://...)")
	require.Equal(t, options[1], "yyy (https://...)")
	require.Equal(t, options[2], custom)
	require.Equal(t, options[3], addMyOwn)
}

func TestGetTemplateResultAddMyOwn(t *testing.T) {
	require.Nil(t, getTemplateResult(addMyOwn, []*templateStruct{}))
}

func TestGetTemplateResultAdd(t *testing.T) {
	templates := []*templateStruct{
		{Name: "xxx", URL: "https://..."},
		{Name: "yyy", URL: "https://..."},
	}
	result := "yyy (https://...)"
	tmpl := getTemplateResult(result, templates)
	require.NotNil(t, tmpl)
	require.Equal(t, tmpl.Name, templates[1].Name)
}

func TestDownloadTemplate(t *testing.T) {
	tmpl := &templateStruct{
		URL: "https://github.com/mesg-foundation/template-service-javascript",
	}
	path, err := downloadTemplate(tmpl)
	defer os.RemoveAll(path)
	require.Nil(t, err)
	require.NotNil(t, path)
}
