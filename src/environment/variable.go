package environment

import (
	"os"
	"strconv"
	"time"
)

type Variable struct {
	key string
}

func NewVariable(key string) *Variable {
	return &Variable{key}
}

func (variable *Variable) GetStringValue(defaultValue string) string {
	if value, exists := os.LookupEnv(variable.key); exists {
		return value
	}
	return defaultValue
}

func (variable *Variable) GetUint64Value(defaultValue uint64) uint64 {
	value, exists := os.LookupEnv(variable.key)
	if !exists {
		return defaultValue
	}
	parsedValue, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return parsedValue
}

func (variable *Variable) GetDurationValue(defaultValue time.Duration) time.Duration {
	value, exists := os.LookupEnv(variable.key)
	if !exists {
		return defaultValue
	}
	parsedValue, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return parsedValue
}
