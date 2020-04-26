package ingester

import "github.com/timdrysdale/gradexpath"

// This must remain idempotent so we can call it every startup
func EnsureDirectoryStructure() error {
	return gradexpath.SetupGradexPaths()
}
