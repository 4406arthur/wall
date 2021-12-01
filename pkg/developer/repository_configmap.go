package developer

import (
	"errors"
	"io/ioutil"
	"wall/pkg/entity"

	"gopkg.in/yaml.v2"
)

//ConfigmapRepository
type ConfigmapRepository struct {
	Rules []entity.Developer `yaml:"dev_acct_rules"`
}

//NewMongoRepository create new repository
func NewConfigmapRepository(configPath string) *ConfigmapRepository {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var repo ConfigmapRepository
	err = yaml.Unmarshal(yamlFile, &repo)
	if err != nil {
		panic(err)
	}

	return &repo
}

//Find a Developer
func (r *ConfigmapRepository) Find(name string) (*entity.Developer, error) {
	// [TODO] 2021/11/25 Linear search should avoid, should change the data struct
	for _, val := range *&r.Rules {
		if val.Name == name {
			return &val, nil
		}
	}

	err := errors.New("no such user")
	return nil, err
}
