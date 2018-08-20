package api

import (
	"os"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
	git "gopkg.in/src-d/go-git.v4"
)

func TestGitCloneRepositoryDoNotExist(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	deployer := newServiceDeployer(a)

	path, _ := deployer.createTempDir()
	defer os.RemoveAll(path)
	err := deployer.gitClone("/doNotExist", path)
	require.NotNil(t, err)
}

func TestGitCloneWithoutURLSchema(t *testing.T) {
	m, _ := newAPIAndDockerTest(t)
	deployer := newServiceDeployer(m)

	path, _ := deployer.createTempDir()
	defer os.RemoveAll(path)
	err := deployer.gitClone("github.com/mesg-foundation/awesome.git", path)
	require.Nil(t, err)
}

func TestGitCloneCustomBranch(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	deployer := newServiceDeployer(a)

	branchName := "5-generic-service"
	path, _ := deployer.createTempDir()
	defer os.RemoveAll(path)
	err := deployer.gitClone("github.com/mesg-foundation/service-ethereum-erc20#"+branchName, path)
	require.Nil(t, err)
	repo, err := git.PlainOpen(path)
	require.Nil(t, err)
	branch, err := repo.Branch(branchName)
	require.Nil(t, err)
	require.NotNil(t, branch)
}

func TestCreateTempFolder(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	deployer := newServiceDeployer(a)

	path, err := deployer.createTempDir()
	defer os.RemoveAll(path)
	require.Nil(t, err)
	require.NotEqual(t, "", path)
}

func TestRemoveTempFolder(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	deployer := newServiceDeployer(a)

	path, _ := deployer.createTempDir()
	err := os.RemoveAll(path)
	require.Nil(t, err)
}

func TestInjectConfigurationInDependencies(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	deployer := newServiceDeployer(a)

	s := &service.Service{}
	deployer.injectConfigurationInDependencies(s, "TestInjectConfigurationInDependencies")
	require.Equal(t, s.Dependencies["service"].Image, "TestInjectConfigurationInDependencies")
}

func TestInjectConfigurationInDependenciesWithConfig(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	deployer := newServiceDeployer(a)

	s := &service.Service{
		Configuration: &service.Dependency{
			Command: "xxx",
			Image:   "yyy",
		},
	}
	deployer.injectConfigurationInDependencies(s, "TestInjectConfigurationInDependenciesWithConfig")
	require.Equal(t, s.Dependencies["service"].Image, "TestInjectConfigurationInDependenciesWithConfig")
	require.Equal(t, s.Dependencies["service"].Command, "xxx")
}

func TestInjectConfigurationInDependenciesWithDependency(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	deployer := newServiceDeployer(a)

	s := &service.Service{
		Dependencies: map[string]*service.Dependency{
			"test": {
				Image: "xxx",
			},
		},
	}
	deployer.injectConfigurationInDependencies(s, "TestInjectConfigurationInDependenciesWithDependency")
	require.Equal(t, s.Dependencies["service"].Image, "TestInjectConfigurationInDependenciesWithDependency")
	require.Equal(t, s.Dependencies["test"].Image, "xxx")
}

func TestInjectConfigurationInDependenciesWithDependencyOverride(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	deployer := newServiceDeployer(a)

	s := &service.Service{
		Dependencies: map[string]*service.Dependency{
			"service": {
				Image: "xxx",
			},
		},
	}
	deployer.injectConfigurationInDependencies(s, "TestInjectConfigurationInDependenciesWithDependencyOverride")
	require.Equal(t, s.Dependencies["service"].Image, "TestInjectConfigurationInDependenciesWithDependencyOverride")
}
