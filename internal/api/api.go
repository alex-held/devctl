package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/ghodss/yaml"
	"github.com/google/go-github/github"
	"github.com/rs/zerolog/log"
	
	"github.com/alex-held/dev-env/internal/spec"
	"github.com/alex-held/dev-env/shared"
)

type GithubContentResponse struct {
	repo    *github.RepositoryContent
	content []*github.RepositoryContent
}

type githubAPI struct {
	github.Client
	Owner         string
	Repository    string
	UpdateChannel string
}

type Context struct {
	shared.PathFactory
}

func NewContext(path shared.PathFactory)  *Context {
	return &Context{path}
}

func (ctx *Context) GetSpecFiles() (specFiles []spec.SpecFile)  {
	home := ctx.GetDevEnvHome()
	specRoot := path.Join(home, "spec")
	files:= getFilePaths(specRoot)
	for _, file := range files {
		fileBytes,_ := ioutil.ReadFile(file)
		s := spec.Spec{}
		_ = yaml.Unmarshal(fileBytes, &s)
		specFile := spec.SpecFile{Path: file, Spec: s}
		specFiles = append(specFiles, specFile)
	}
	return specFiles
}

func getFilePaths(root string) (files []string ) {
	fi,_ := ioutil.ReadDir(root)
	for _, info := range fi {
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		subdirFiles := getFilePaths(info.Name())
		files = append(files, subdirFiles...)
	}
	return files
}


func (client *githubAPI) GetPackages(path string) (packages []string, err error) {
	versionMap, err := client.GetPackagesMap(path)
	if err != nil {
		log.Error().Err(err)
		return packages, err
	}
	linq.From(versionMap).
		SelectT(func(kvp linq.KeyValue) string {
			key := kvp.Key.(string)
			return key
		}).ToSlice(&packages)
	return packages, nil
}

func (client *githubAPI) GetPackageVersions(path string) (result []string, err error) {
	versionMap, err := client.GetPackageVersionsMap(path)
	if err != nil {
		log.Error().Err(err)
		return result, err
	}
	linq.From(versionMap).
		SelectT(func(kvp linq.KeyValue) string {
			key := kvp.Key.(string)
			return key
		}).ToSlice(&result)
	return result, nil
}

type GithubAPI interface {
	GetPackagesMap(path string) (result map[string]interface{}, err error)
	GetPackages(path string) (packages []string, err error)
	GetPackageVersionsMap(path string) (result map[string]interface{}, err error)
	GetPackageVersions(path string) (result []string, err error)
}

func (r *GithubContentResponse) FilterAndMap(predicate func(content *github.RepositoryContent) bool) map[string]interface{} {
	result := map[string]interface{}{}
	linq.From(r.content).
		WhereT(predicate).
		SelectT(func(c *github.RepositoryContent) linq.KeyValue {
			kvp := linq.KeyValue{
				Key:   *c.Name,
				Value: *c,
			}
			return kvp
		}).ToMap(&result)
	return result
}

var (
	FilePredicate = func(c *github.RepositoryContent) bool {
		return *c.Type == "file"
	}
	DirectoryPredicate = func(c *github.RepositoryContent) bool {
		return *c.Type == "dir"
	}
)

func (client *githubAPI) GetPackagesMap(path string) (result map[string]interface{}, err error) {
	response, err := client.GetContent(path)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	result = response.FilterAndMap(DirectoryPredicate)
	return result, err
}

func (client *githubAPI) GetPackageVersionsMap(path string) (result map[string]interface{}, err error) {
	response, err := client.GetContent(path)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	content := response.FilterAndMap(FilePredicate)

	result = map[string]interface{}{}
	for key, value := range content {
		newKey := strings.TrimSuffix(key, ".yaml")
		result[newKey] = value
	}

	return result, err
}

func (client *githubAPI) GetContent(path string) (response GithubContentResponse, err error) {
	options := &github.RepositoryContentGetOptions{Ref: client.UpdateChannel}
	repo, directoryContent, _, err := client.Repositories.GetContents(context.Background(), client.Owner, client.Repository, path, options)
	response = GithubContentResponse{repo: repo, content: directoryContent}
	return response, err
}

func NewGithubAPI(client *http.Client) GithubAPI {
	return &githubAPI{
		Client:        *github.NewClient(client),
		Owner:         "alex-held",
		Repository:    "dev-env-registry",
		UpdateChannel: "master",
	}
}

func NewTestGithubAPI(client *github.Client) GithubAPI {
	a := githubAPI{
		Client:        *client,
		Owner:         "alex-held",
		Repository:    "dev-env-registry",
		UpdateChannel: "master",
	}
	return &a
}
