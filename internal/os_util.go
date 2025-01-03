package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	ErrUnsupportedPlatform = fmt.Errorf("unsupported platform")
)

func GetAppDataDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("AppData")
		if dir == "" {
			return "", fmt.Errorf("AppData environment variable not set")
		}
		dir = filepath.Join(dir, "lexApp")
	case "linux":
		d, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(d, ".local", "share", "lexApp")
	default:
		return "", ErrUnsupportedPlatform
	}

	return dir, nil
}
