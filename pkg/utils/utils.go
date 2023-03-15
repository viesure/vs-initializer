package utils

import (
	"os"
)

// GetEnvOrDefault retrieves the value of the environment variable named by
// `key`. If the variable is present in the environment,
// the value (which may be empty) is returned. Otherwise the value specified by
// `defaultValue` will be returned.
func GetEnvOrDefault(key string, defaultValue string) string {

	value, exist := os.LookupEnv(key)

	if !exist {
		return defaultValue
	}

	return value

}

// ReadFileOrEmpty retrieves a filename string given by `filename`. If the
// gicen file exists, the functions returns the content, otherwise an empty
// string.
func ReadFileOrEmpty(filename string) string {

	content, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}

	return string(content)
}

// NormalizeDirectoryPath retrieves a path string given by `path`. If the
// variable ends with a path separator, this function returns the path without
// the last path separator.
func NormalizeDirectoryPath(path string) string {

	if len(path) > 0 && path[len(path)-1:] == string(os.PathSeparator) {
		return path[:len(path)-1]
	}

	return path

}
