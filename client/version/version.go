package version

type Version struct {
	Commit string `json:"commit"`
	Tag    string `json:"tag"`
}
