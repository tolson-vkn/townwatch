package version

// Set at build
// go build -ldflags "-X github.com/tolson-vkn/townwatch/version.Version=$(git describe --abbrev=0)
// commit: $(git rev-parse HEAD)
var (
	Version   string = "UnkownVER"
	GitCommit string = "UnkownSHA"
)

func VersionString() string {
	return Version
}

func GitCommitString() string {
	return GitCommit
}
