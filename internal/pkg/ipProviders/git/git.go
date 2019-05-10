package git

import (
	"errors"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	yaml "gopkg.in/yaml.v2"

	"github.com/mitchellh/mapstructure"

	"github.com/stakater/Whitelister/internal/pkg/utils"
)

var pullOptions *git.PullOptions = &git.PullOptions{RemoteName: "origin"}
var path = "/tmp/whitelister-config"

// Git Ip provider class implementing the IpProvider interface
type Git struct {
	AccessToken string
	URL         string
	Config      string
	repository  *git.Repository
	workingTree *git.Worktree
}

// Equal Compares Git objects
func (git1 *Git) Equal(git2 *Git) bool {
	if git1.URL != git2.URL ||
		git1.AccessToken != git2.AccessToken ||
		git1.Config != git2.Config {
		return false
	}
	return true
}

//Config stores IpPermissions read from config file.
type Config struct {
	IpPermissions []utils.IpPermission `yaml:"ipPermissions"`
}

// GetName returns the name of IP Provider
func (g *Git) GetName() string {
	return "Git"
}

// Init initializes the Git Configuration
func (g *Git) Init(params map[interface{}]interface{}) error {
	err := mapstructure.Decode(params, &g) //Converts the params to Git struct fields
	if err != nil {
		return err
	}

	if g.AccessToken == "" {
		return errors.New("Missing Git Access Token")
	}

	if g.URL == "" {
		return errors.New("Missing Git URL")
	}

	if g.Config == "" {
		g.Config = "config.yaml"
	}

	return g.cloneRepository()
}

// GetIPPermissions - Get List of IP addresses to whitelist
func (g *Git) GetIPPermissions() ([]utils.IpPermission, error) {
	err := g.pullRepository()

	if err != nil {
		return nil, err
	}

	conf, err := g.readConfig()

	if err != nil {
		return nil, err
	}

	return conf.IpPermissions, nil
}

func (g *Git) cloneRepository() error {
	var err error
	// Clone the given repository, creating the remote, the local branches
	// and fetching the objects, exactly as:
	logrus.Infof("Cloning Repo %s at %s", g.URL, path)

	options := &git.CloneOptions{
		URL: g.URL,
		Auth: &http.BasicAuth{
			Username: "Stakater", //can be anything except empty string
			Password: g.AccessToken,
		},
	}
	g.repository, err = git.PlainClone(path, false, options)

	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			logrus.Infof("Repo already cloned. Cotinuing")
			g.repository, err = git.PlainOpen(path)
			if err != nil {
				return err
			}
		} else {
			if err == git.ErrRepositoryAlreadyExists {
				return err
			}
		}
	}

	// Get the working directory for the repository
	g.workingTree, err = g.repository.Worktree()
	if err != nil {
		return err
	}

	return nil
}

func (g *Git) printLatestCommit() {
	ref, err := g.repository.Head()
	if err != nil {
		logrus.Errorf("Unable to get head : %v", err)
		return
	}

	commit, err := g.repository.CommitObject(ref.Hash())

	if err != nil {
		logrus.Errorf("Unable to get commit : %v", err)
		return
	}

	logrus.Infof("Pulled %v", commit)

}

func (g *Git) pullRepository() error {
	err := g.workingTree.Pull(pullOptions)

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			logrus.Info("No changes to fetch from git")
		} else {
			return err
		}
	} else {
		// Print the latest commit that was just pulled
		g.printLatestCommit()
	}

	return nil
}

func (g *Git) readConfig() (Config, error) {
	var config Config
	// Read YML file
	source, err := ioutil.ReadFile(path + "/" + g.Config)
	if err != nil {
		return config, err
	}

	// Unmarshall
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
