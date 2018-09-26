package main

import (
	"github.com/kelseyhightower/envconfig"
)

// Specification are env variables used by ems
// split_words is used by envconfig
type Specification struct {
	B2ApplicationID         string `split_words:"true"`
	B2ApplicationKey        string `split_words:"true"`
	B2BucketName            string `split_words:"true"`
	B2Object                string `split_words:"true"`
	B2ObjectID              string `split_words:"true"`
	B2DownloadRetryInterval int    `split_words:"true"`
}

// ParseENV pareses env for variables
func ParseENV() (*Specification, error) {
	s := Specification{}

	err := envconfig.Process("EBS", &s)

	if err != nil {
		return nil, err
	}

	return &s, nil
}
