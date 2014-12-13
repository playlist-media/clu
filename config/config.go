package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	cconfig "github.com/coreos/coreos-cloudinit/config"
	"gopkg.in/yaml.v2"
)

type Instance struct {
	Name        string              `yaml:"name"`
	Kind        string              `yaml:"kind"`
	MachineType string              `yaml:"machine_type"`
	MachineOpts string              `yaml:"machine_opts"`
	CloudConfig cconfig.CloudConfig `yaml:"cloud_config"`
}

type Config struct {
	Global struct {
		ProjectID    string `yaml:"project_id" env:"PROJECT_ID"`
		DiscoveryURL string `yaml:"discovery_url" env:"DISCOVERY_URL"`
		Zone         string `yaml:"zone" env:"ZONE"`
	} `yaml:"global"`
	Instances map[string]Instance `yaml:"instances"`
	Kinds     map[string]cconfig.CloudConfig
	Units     map[string]cconfig.Unit
}

func NewConfig(contents string) (*Config, error) {
	var cfg Config
	err := yaml.Unmarshal([]byte(contents), &cfg)
	if err != nil {
		return &cfg, err
	}

	et := reflect.TypeOf(cfg.Global)
	ev := reflect.ValueOf(&cfg.Global)

	for i := 0; i < et.NumField(); i++ {
		if key := et.Field(i).Tag.Get("env"); key != "" {
			if env := os.Getenv(key); env != "" {
				ev.Elem().Field(i).SetString(env)
			}
		}
	}

	return &cfg, nil
}

func NewConfigFromFile(filename string) (*Config, error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c, err := NewConfig(string(d))
	return c, err
}

func (c Config) String() string {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return ""
	}

	stringified := string(bytes)
	stringified = fmt.Sprintf("#clu-config\n%s", stringified)

	return stringified
}
