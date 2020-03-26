package utils

type DBTypeEnum int

const (
	DBType_MONGODB = iota
)

type EnvTypeEnum int

const (
	PRODUCTION = iota
	TESTING
)
