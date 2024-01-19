package deltatocumulativeprocessor

import (
	"go.opentelemetry.io/collector/component"
)

var _ component.ConfigValidator = (*Config)(nil)

type Config struct{}

func (c *Config) Validate() error {
	return nil
}
