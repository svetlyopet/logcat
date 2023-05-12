package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// PrintHelp prints out to stdout help information about this program and exits
func PrintHelp() {
	fmt.Println("Usage: logcat -file [FILEPATH] -outdir [DIRECTORY]")
	fmt.Println("Example: logcat -dir /opt/artifactory/apps/artifactory/var/log/artifactory-request.log -out /tmp")
	os.Exit(1)
}

// IsAbsolutePath checks if input string is an absolute path
func IsAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}

// AppendSlashIfMissing checks if string ends with slash and if not it appends it
func AppendSlashIfMissing(path string) string {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	return path
}
