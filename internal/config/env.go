package config

import "strings"

type Env string

const (
	Development Env = "development"
	Production  Env = "production"
)

var currentEnv = Production

func SetEnv(env string) {
	if Env(strings.ToLower(env)) == Development {
		currentEnv = Development
	} else {
		currentEnv = Production
	}
}

func IsDev() bool {
	return currentEnv == Development
}

func IsProd() bool {
	return currentEnv == Production
}
