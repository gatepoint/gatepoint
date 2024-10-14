package config

func ValueOrDefault[T comparable](val, defaultVal T) T {
	if val == *new(T) {
		return defaultVal
	}
	return val
}
