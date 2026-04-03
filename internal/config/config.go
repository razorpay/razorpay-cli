package config

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AppConfig struct {
	AWS AWSConfig `mapstructure:"aws"`
}

type AWSConfig struct {
	Region   string `mapstructure:"region"`
	S3Bucket string `mapstructure:"s3_bucket"`
}

func (c *AppConfig) Validate() error {
	if err := c.AWS.Validate(); err != nil {
		return fmt.Errorf("aws: %w", err)
	}
	return nil
}

func (c AWSConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Region,
			validation.Required,
			validation.In("us-east-1", "us-west-2", "ap-south-1"),
		),
		validation.Field(&c.S3Bucket,
			validation.Required,
			validation.Length(3, 63),
		),
	)
}
