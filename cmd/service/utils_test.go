package service

import (
	"os"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
	git "gopkg.in/src-d/go-git.v4"
)

func TestDefaultPath(t *testing.T) {
	require.Equal(t, defaultPath([]string{}), "./")
	require.Equal(t, defaultPath([]string{"foo"}), "foo")
	require.Equal(t, defaultPath([]string{"foo", "bar"}), "foo")
}

func TestBuildDockerImagePathDoNotExist(t *testing.T) {
	_, err := buildDockerImage("/doNotExist")
	require.NotNil(t, err)
}

func TestGitCloneRepositoryDoNotExist(t *testing.T) {
	path, _ := createTempFolder()
	defer os.RemoveAll(path)
	err := gitClone("/doNotExist", path, "testing...")
	require.NotNil(t, err)
}

func TestGitCloneWithoutURLSchema(t *testing.T) {
	path, _ := createTempFolder()
	defer os.RemoveAll(path)
	err := gitClone("github.com/mesg-foundation/awesome.git", path, "testing...")
	require.Nil(t, err)
}

func TestGitCloneCustomBranch(t *testing.T) {
	branchName := "5-generic-service"
	path, _ := createTempFolder()
	defer os.RemoveAll(path)
	err := gitClone("github.com/mesg-foundation/service-ethereum-erc20#"+branchName, path, "testing...")
	require.Nil(t, err)
	repo, err := git.PlainOpen(path)
	require.Nil(t, err)
	branch, err := repo.Branch(branchName)
	require.Nil(t, err)
	require.NotNil(t, branch)
}

func TestDownloadServiceIfNeededAbsolutePath(t *testing.T) {
	path := "/users/paul/service-js-function"
	newPath, didDownload, err := downloadServiceIfNeeded(path)
	require.Nil(t, err)
	require.Equal(t, path, newPath)
	require.Equal(t, false, didDownload)
}

func TestDownloadServiceIfNeededRelativePath(t *testing.T) {
	path := "./service-js-function"
	newPath, didDownload, err := downloadServiceIfNeeded(path)
	require.Nil(t, err)
	require.Equal(t, path, newPath)
	require.Equal(t, false, didDownload)
}

func TestDownloadServiceIfNeededUrl(t *testing.T) {
	path := "https://github.com/mesg-foundation/awesome.git"
	newPath, didDownload, err := downloadServiceIfNeeded(path)
	defer os.RemoveAll(newPath)
	require.Nil(t, err)
	require.NotEqual(t, path, newPath)
	require.Equal(t, true, didDownload)
}

func TestCreateTempFolder(t *testing.T) {
	path, err := createTempFolder()
	defer os.RemoveAll(path)
	require.Nil(t, err)
	require.NotEqual(t, "", path)
}

func TestRemoveTempFolder(t *testing.T) {
	path, _ := createTempFolder()
	err := os.RemoveAll(path)
	require.Nil(t, err)
}

func TestInjectConfigurationInDependencies(t *testing.T) {
	s := &service.Service{}
	injectConfigurationInDependencies(s, "TestInjectConfigurationInDependencies")
	require.Equal(t, s.Dependencies["service"].Image, "TestInjectConfigurationInDependencies")
}

func TestInjectConfigurationInDependenciesWithConfig(t *testing.T) {
	s := &service.Service{
		Configuration: &service.Dependency{
			Command: "xxx",
			Image:   "yyy",
		},
	}
	injectConfigurationInDependencies(s, "TestInjectConfigurationInDependenciesWithConfig")
	require.Equal(t, s.Dependencies["service"].Image, "TestInjectConfigurationInDependenciesWithConfig")
	require.Equal(t, s.Dependencies["service"].Command, "xxx")
}

func TestInjectConfigurationInDependenciesWithDependency(t *testing.T) {
	s := &service.Service{
		Dependencies: map[string]*service.Dependency{
			"test": {
				Image: "xxx",
			},
		},
	}
	injectConfigurationInDependencies(s, "TestInjectConfigurationInDependenciesWithDependency")
	require.Equal(t, s.Dependencies["service"].Image, "TestInjectConfigurationInDependenciesWithDependency")
	require.Equal(t, s.Dependencies["test"].Image, "xxx")
}

func TestInjectConfigurationInDependenciesWithDependencyOverride(t *testing.T) {
	s := &service.Service{
		Dependencies: map[string]*service.Dependency{
			"service": {
				Image: "xxx",
			},
		},
	}
	injectConfigurationInDependencies(s, "TestInjectConfigurationInDependenciesWithDependencyOverride")
	require.Equal(t, s.Dependencies["service"].Image, "TestInjectConfigurationInDependenciesWithDependencyOverride")
}
