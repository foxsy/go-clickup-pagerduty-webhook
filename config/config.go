package config

import (
    "gopkg.in/yaml.v3"
    "io/ioutil"
    "log"
)

type Condition struct {
    Key   string `yaml:"key"`
    Value string `yaml:"value"`
}

type Rule struct {
    Event     string    `yaml:"event"`
    Condition Condition `yaml:"condition"`
    Action    string    `yaml:"action"`
    Group     string    `yaml:"group"`
    List	  string	`yaml:"list"`
    Space	  string	`yaml:"space"`
}

type Config struct {
    Rules []Rule `yaml:"rules"`
}

var AppConfig Config

func LoadConfig(path string) {
    yamlFile, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatalf("Error reading YAML file: %v", err)
    }
    err = yaml.Unmarshal(yamlFile, &AppConfig)
    if err != nil {
        log.Fatalf("Error parsing YAML file: %v", err)
    }
}
