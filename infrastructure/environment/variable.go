package environment

import (
	"os"
	"strconv"
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

func (variable *Variable) GetIntValue(defaultValue int) int {
	value, exists := os.LookupEnv(variable.key)
	if !exists {
		return defaultValue
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsedValue
}
