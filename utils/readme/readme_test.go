package readme

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookupReadme(t *testing.T) {
	tests := []struct {
		file    string
		present bool
	}{
		{file: "Readme", present: true},
		{file: "Readme.md", present: true},
		{file: "readme", present: true},
		{file: "readme.md", present: true},
		{file: "README", present: true},
		{file: "README.md", present: true},
		{file: "xxx.md", present: false},
	}
	for _, test := range tests {
		dir, err := ioutil.TempDir("./", "mesg")
		require.NoError(t, err)
		defer os.RemoveAll(dir)
		err = ioutil.WriteFile(filepath.Join(dir, test.file), []byte("hello"), os.ModePerm)
		require.NoError(t, err)
		x, err := LookupReadme(dir)
		require.NoError(t, err)
		if test.present {
			require.Equal(t, "hello", x)
		} else {
			require.Equal(t, "", x)
		}
	}
}
