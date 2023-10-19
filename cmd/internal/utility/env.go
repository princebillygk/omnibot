package utility

import (
	"os"
	"strconv"
)

type EnvValue interface {
	int64 | string | bool | uint64 | float64
}

func GetEnv[T EnvValue](k string, fallback T) T {
	v, ok := os.LookupEnv(k)
	if !ok {
		return fallback
	}

	var (
		pv  T
		err error = nil
	)

	switch p := any(&pv).(type) {
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
		return fallback
	}
	return pv
}
