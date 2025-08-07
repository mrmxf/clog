package kfg

import (
	"errors"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

// remember the position of the release slice for the convenience funtion
var releaseSliceCache *[]AppRelease

func ReleasesPath() string {
	if Raw == nil {
		return ""
	}
	return Raw.String(KongifReleasesPathKey)
}

func CurrentRelease() *AppRelease {
	// if Konfigure has not been called or releases not loaded
	if Raw == nil || releaseSliceCache == nil || len(*releaseSliceCache) == 0 {
		return nil
	}
	// return a pointer to index 0 of the cached loaded releases
	return &(*releaseSliceCache)[0]
}

// LoadReleases populates a slice of type []AppRelease.
//
// How it works:
//   - 1. Validates that Konfigure() has been initialized
//   - 2. Reads the releases.yaml file from the filesystem
//   - 3. Parses the YAML content into []AppRelease
//   - 4. Populates AppConfig.Releases with the parsed data
//   - 5. Logs the success or failure using slog.Debug
//
// Returns an error if the file cannot be read or parsed, nil if successful.
func LoadReleases(destination *[]AppRelease) error {
	// Check if Konfigure() has been called first
	if Raw == nil {
		return errors.New("configuration not initialized: call Konfigure() before using LoadReleases()")
	}

	// Check if there is work to do
	if destination == nil {
		slog.Debug("nil destination in kfg.LoadReleases, nothing to do.")
		return nil
	}

	// Check if the releases path has been initialized
	if len(ReleasesPath()) < 6 {
		slog.Debug("len(" + KongifReleasesPathKey + ") value in Konfig is too short")
		return nil
	}

	slog.Debug("LoadReleases: attempting to load releases data", "file_path", ReleasesPath())

	// Read the releases.yaml file from the filesystem
	data, err := os.ReadFile(ReleasesPath())
	if err != nil {
		slog.Debug("LoadReleases: failed to read releases file", "file_path", ReleasesPath(), "error", err)
		return err
	}

	// Parse the YAML content into []AppRelease
	err = yaml.Unmarshal(data, destination)
	if err != nil {
		slog.Debug("LoadReleases: failed to parse releases YAML", "file_path", ReleasesPath(), "error", err)
		return err
	}

	// we managed to load the releases so cache the latest one
	releaseSliceCache = destination
	slog.Debug("LoadReleases: successfully loaded releases", "file_path", ReleasesPath(), "count", len(*destination))
	return nil
}
