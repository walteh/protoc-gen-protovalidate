package main

import (
	"encoding/json"
	"strings"
)

const (
	DefaultBufValidateFile        = "buf/validate/validate.proto"
	DefaultProtoValidateGoRef     = "_local"
	DefaultProtoValidatePythonRef = "_local"
	DefaultProtoValidateCCRef     = "_local"
	DefaultProtoValidateJavaRef   = "_local"
)

type Config struct {
	BufValidateFile        *string `json:"buf_validate_file,omitempty"`
	ProtoValidateGoRef     *string `json:"protovalidate_go_ref,omitempty"`
	ProtoValidatePythonRef *string `json:"protovalidate_python_ref,omitempty"`
	ProtoValidateCCRef     *string `json:"protovalidate_cc_ref,omitempty"`
	ProtoValidateJavaRef   *string `json:"protovalidate_java_ref,omitempty"`
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

func (c *Config) GetProtoValidateGoRef() string {
	if c.ProtoValidateGoRef == nil {
		return DefaultProtoValidateGoRef
	}
	return *c.ProtoValidateGoRef
}

func (c *Config) GetProtoValidatePythonRef() string {
	if c.ProtoValidatePythonRef == nil {
		return DefaultProtoValidatePythonRef
	}
	return *c.ProtoValidatePythonRef
}

func (c *Config) GetProtoValidateCCRef() string {
	if c.ProtoValidateCCRef == nil {
		return DefaultProtoValidateCCRef
	}
	return *c.ProtoValidateCCRef
}

func (c *Config) GetProtoValidateJavaRef() string {
	if c.ProtoValidateJavaRef == nil {
		return DefaultProtoValidateJavaRef
	}
	return *c.ProtoValidateJavaRef
}
