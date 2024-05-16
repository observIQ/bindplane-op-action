package model

import (
	"strconv"
	"strings"
)

func TrimVersion(resourceKey string) string {
	key, _ := SplitVersion(resourceKey)
	return key
}

func SplitVersion(resourceKey string) (string, Version) {
	return SplitVersionDefault(resourceKey, VersionLatest)
}

func SplitVersionDefault(resourceKey string, defaultVersion Version) (string, Version) {
	parts := strings.SplitN(resourceKey, ":", 2)
	name := parts[0]
	if len(parts) == 1 {
		return name, defaultVersion
	}
	switch parts[1] {
	case "":
		return name, defaultVersion
	case "latest":
		return name, VersionLatest
	case "stable", "current":
		return name, VersionCurrent
	case "pending":
		return name, VersionPending
	}
	version, err := strconv.Atoi(parts[1])
	if err != nil {
		return name, defaultVersion
	}
	return name, Version(version)
}

type Version int

const (
	// VersionPending refers to the pending Version of a resource, which is the version that is currently being rolled
	// out. This is currently only used for Configurations.
	VersionPending Version = -2

	// VersionCurrent refers to the current Version of a resource, which is the last version to be successfully rolled
	// out. This is currently only used for Configurations.
	VersionCurrent Version = -1

	// VersionLatest refers to the latest Version of a resource which is the latest version that has been created.
	VersionLatest Version = 0
)
