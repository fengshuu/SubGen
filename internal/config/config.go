package config

import (
    "os"
    "strings"

    "gopkg.in/yaml.v3"
)

type Subscription struct {
    Name             string `yaml:"name"`
    URL              string `yaml:"url"`
    Path             string `yaml:"path,omitempty"`
    AdditionalPrefix string `yaml:"additional_prefix,omitempty"`
}

type AppConfig struct {
    BaseConfigURL string         `yaml:"base_config_url"`
    Subscriptions []Subscription `yaml:"subscriptions"`
}

func Load(path string) (*AppConfig, error) {
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var cfg AppConfig
    if err := yaml.Unmarshal(b, &cfg); err != nil {
        return nil, err
    }
    if strings.TrimSpace(cfg.BaseConfigURL) == "" {
        cfg.BaseConfigURL = "https://gist.githubusercontent.com/liuran001/5ca84f7def53c70b554d3f765ff86a33/raw/"
    }
    return &cfg, nil
}