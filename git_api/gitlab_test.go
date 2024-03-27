package git_api

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"os"
	"testing"
)

func getStubProject() *gitlab.Project {
	return &gitlab.Project{
		ID:   999,
		Name: "TestProject",
	}
}

func getStubGitlabResponse(code int) *gitlab.Response {
	return &gitlab.Response{
		Response: &http.Response{
			StatusCode: code,
		},
	}
}

func getGitlabCodeTestClient(t *testing.T) *GitlabApi {

	err := os.Setenv("CFS_CODE_GITLAB_USER", "testUser")
	assert.Nil(t, err)
	err = os.Setenv("CFS_CODE_GITLAB_PAT", "testPat")
	assert.Nil(t, err)

	return GetGitlabCodeClient("http://localhost:1234/gitlab/test/instance")
}

func TestGetProject_ValidProjectName_ProjectFound(t *testing.T) {

	c := getGitlabCodeTestClient(t)

	c.getProjectFunc = func(c *gitlab.Client, pid interface{}, opt *gitlab.GetProjectOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Project, *gitlab.Response, error) {
		return getStubProject(), getStubGitlabResponse(200), nil
	}

	p, err := c.getProject("TestProject")

	assert.Nil(t, err)
	assert.NotNil(t, p)
}

func TestGetProject_InvalidProjectName_ProjectNotFound(t *testing.T) {

	c := getGitlabCodeTestClient(t)

	c.getProjectFunc = func(c *gitlab.Client, pid interface{}, opt *gitlab.GetProjectOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Project, *gitlab.Response, error) {
		return nil, getStubGitlabResponse(404), errors.New("project not found")
	}

	p, err := c.getProject("TestProject_NonExistent")

	assert.NotNil(t, err)
	assert.Nil(t, p)
}

func TestCreateProject_ValidProjectName_CreatedSuccessfully(t *testing.T) {

	c := getGitlabCodeTestClient(t)

	c.getProjectFunc = func(c *gitlab.Client, pid interface{}, opt *gitlab.GetProjectOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Project, *gitlab.Response, error) {
		return nil, getStubGitlabResponse(404), errors.New("project not found")
	}

	c.createProjectFunc = func(c *gitlab.Client, opt *gitlab.CreateProjectOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Project, *gitlab.Response, error) {
		return getStubProject(), getStubGitlabResponse(201), nil
	}

	p, err := c.createProject("TestProject")

	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.True(t, p.Name == getStubProject().Name)
}

func TestCreateProject_ExistingProjectName_CreationFailed(t *testing.T) {

	c := getGitlabCodeTestClient(t)

	c.getProjectFunc = func(c *gitlab.Client, pid interface{}, opt *gitlab.GetProjectOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Project, *gitlab.Response, error) {
		return getStubProject(), getStubGitlabResponse(200), nil
	}

	p, err := c.createProject("TestProject")

	assert.NotNil(t, err)
	assert.Nil(t, p)
	assert.Contains(t, err.Error(), "found a Gitlab project with the name")
}

func TestValidateNewCode_ExistingProjectName_ProjectFoundError(t *testing.T) {

	c := getGitlabCodeTestClient(t)

	c.getProjectFunc = func(c *gitlab.Client, pid interface{}, opt *gitlab.GetProjectOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Project, *gitlab.Response, error) {
		return getStubProject(), getStubGitlabResponse(200), nil
	}

	err := c.ValidateNewCode("TestProject")

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "found Project when we did not expect one")
}

func TestValidateNewCode_NewProjectName_ProjectNotFound(t *testing.T) {

	c := getGitlabCodeTestClient(t)

	c.getProjectFunc = func(c *gitlab.Client, pid interface{}, opt *gitlab.GetProjectOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Project, *gitlab.Response, error) {
		return nil, getStubGitlabResponse(404), errors.New("project not found")
	}

	err := c.ValidateNewCode("TestProject")

	assert.Nil(t, err)
}

func TestValidateNewCode_NewProjectName_ProjectNotFoundAndUnexpectedError(t *testing.T) {

	c := getGitlabCodeTestClient(t)

	c.getProjectFunc = func(c *gitlab.Client, pid interface{}, opt *gitlab.GetProjectOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Project, *gitlab.Response, error) {
		return nil, getStubGitlabResponse(400), nil
	}

	err := c.ValidateNewCode("TestProject")

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unexpected error returned")
}
