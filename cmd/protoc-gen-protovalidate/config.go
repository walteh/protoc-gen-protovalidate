package main

import (
	"encoding/json"
	"strings"
)

const (
	DefaultBufValidateFile  = "buf/validate/validate.proto"
	DefaultLanguage         = "go"
	DefaultProtoValidateRef = "_local"
)

type Config struct {
	BufValidateFile  *string `json:"buf_validate_file,omitempty"`
	Language         *string `json:"language,omitempty"`
	ProtoValidateRef *string `json:"protovalidate_ref,omitempty"`
}

func parseParams(params string) (*Config, error) {
	paramMap := map[string]string{}

	for _, param := range strings.Split(params, ",") {
		parts := strings.SplitN(param, "=", 2)
		if len(parts) != 2 {
			continue
		}

		paramMap[parts[0]] = parts[1]
	}

	// convert to json
	by, err := json.Marshal(paramMap)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(by, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) GetBufValidateFile() string {
	if c.BufValidateFile == nil {
		return DefaultBufValidateFile
	}
	return *c.BufValidateFile
}

func (c *Config) GetLanguage() string {
	if c.Language == nil {
		return DefaultLanguage
	}
	return *c.Language
}

func (c *Config) GetProtoValidateRef() string {
	if c.ProtoValidateRef == nil {
		return DefaultProtoValidateRef
	}
	return *c.ProtoValidateRef
}
