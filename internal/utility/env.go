package utility

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type EnvValue interface {
	int64 | string | bool | uint64 | float64
}

func fetchEnv[T EnvValue](k string) (T, error) {
	var (
		parsedVal T
		err       error = nil
	)

	v, ok := os.LookupEnv(k)
	if !ok {
		return parsedVal, errors.New(fmt.Sprintf("Environment variable \"%s\" is not defined", k))
	}

	switch p := any(&parsedVal).(type) {
	case *string:
		*p = v
	case *int64:
		*p, err = strconv.ParseInt(v, 10, 64)
	case *bool:
		*p, err = strconv.ParseBool(v)
	case *uint64:
		*p, err = strconv.ParseUint(v, 10, 64)
	case *float64:
		*p, err = strconv.ParseFloat(v, 64)
	default:
		panic("Unknown type environment value!")
	}

	if err != nil {
		return parsedVal, fmt.Errorf("Failed to parse Environment variable \"%s\"->%v: %s", k, v, err.Error())
	}
	return parsedVal, nil
}

func MustGetEnv[T EnvValue](k string) T {
	v, err := fetchEnv[T](k)
	if err != nil {
		panic(err)
	}
	return v
}

func GetEnv[T EnvValue](k string, fallback T) T {
	v, err := fetchEnv[T](k)
	if err != nil {
		return fallback
	}
	return v
}
