package config

import (
    "gopkg.in/yaml.v3"
    "io/ioutil"
    "log"
)

type Group struct {
    Name               string `yaml:"name"`
    EscalationPolicyID string `yaml:"escalation_policy_id"`
}

type GroupConfig struct {
    Groups []Group `yaml:"groups"`
}

var GroupAppConfig GroupConfig

func LoadGroupConfig(file string) {
    yamlFile, err := ioutil.ReadFile(file)
    if err != nil {
        log.Fatalf("Error reading groups YAML file: %v", err)
    }
    err = yaml.Unmarshal(yamlFile, &GroupAppConfig)
    if err != nil {
        log.Fatalf("Error parsing groups YAML file: %v", err)
    }
}
