package utils

import "strings"

func WinExt(fullpath string) string {
	parts := strings.Split(fullpath, ".")
	return parts[len(parts)-1]
}

func WinBase(fullpath string) string {
	parts := strings.Split(fullpath, "\\")
	return parts[len(parts)-1]
}

func WinDir(fullpath string) string {
	parts := strings.Split(fullpath, "\\")
	return strings.Join(parts[:len(parts)-1], "\\")
}
