package service

import (
	"os"
	"testing"

	"github.com/stvp/assert"
)

func TestGetTemplate(t *testing.T) {
	templates, err := getTemplates(templatesURL)
	assert.Nil(t, err)
	assert.True(t, len(templates) > 0)
}

func TestGetTemplateBadURL(t *testing.T) {
	_, err := getTemplates("...")
	assert.NotNil(t, err)
}

func TestTemplatesToOption(t *testing.T) {
	templates := []*template{
		&template{Name: "xxx", URL: "https://..."},
		&template{Name: "yyy", URL: "https://..."},
	}
	options := templatesToOption(templates)
	assert.Equal(t, len(options), len(templates)+2)
	assert.Equal(t, options[0], "xxx (https://...)")
	assert.Equal(t, options[1], "yyy (https://...)")
	assert.Equal(t, options[2], custom)
	assert.Equal(t, options[3], addMyOwn)
}

func TestGetTemplateResultAddMyOwn(t *testing.T) {
	assert.Nil(t, getTemplateResult(addMyOwn, []*template{}))
}

func TestGetTemplateResultAdd(t *testing.T) {
	templates := []*template{
		&template{Name: "xxx", URL: "https://..."},
		&template{Name: "yyy", URL: "https://..."},
	}
	result := "yyy (https://...)"
	tmpl := getTemplateResult(result, templates)
	assert.NotNil(t, tmpl)
	assert.Equal(t, tmpl.Name, templates[1].Name)
}

func TestDownloadTemplate(t *testing.T) {
	tmpl := &template{
		URL: "https://github.com/mesg-foundation/template-service-javascript",
	}
	path, err := downloadTemplate(tmpl)
	defer os.RemoveAll(path)
	assert.Nil(t, err)
	assert.NotNil(t, path)
}
