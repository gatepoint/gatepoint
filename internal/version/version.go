package version

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"

	"sigs.k8s.io/yaml"
)

var (
	binVersion  string
	gitBranch   string
	gitCommitID string
	buildDate   string
)

type VersionInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`

	GitBranch string `json:"gitBranch"`
	Compiler  string `json:"compiler"`
	Platform  string `json:"platform"`
}

// Version returns the version information about the application
func Version() *VersionInfo {
	return &VersionInfo{
		Version:   binVersion,
		GitCommit: gitCommitID,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),

		GitBranch: gitBranch,
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func Print(w io.Writer, format string) error {
	v := Version()
	switch format {
	case "json":
		if m, err := json.MarshalIndent(v, "", "  "); err == nil {
			_, _ = fmt.Fprintln(w, string(m))
		}
	case "yaml":
		if m, err := yaml.Marshal(v); err == nil {
			_, _ = fmt.Fprintln(w, string(m))
		}
	default:
		if m, err := json.MarshalIndent(v, "", "  "); err == nil {
			_, _ = fmt.Fprintln(w, string(m))
		}
	}
	return nil
}
