package cmd

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version information",
	Run: func(cmd *cobra.Command, args []string) {
		buf := bytes.NewBuffer(nil)
		_, _ = fmt.Fprintf(buf, "Appname      %s\n", name())
		_, _ = fmt.Fprintf(buf, "Version      %s\n", versionInfo())
		_, _ = fmt.Fprintf(buf, "\n")
		if trace() != unavailable {
			_, _ = fmt.Fprintf(buf, "GitTrace     %s\n", trace())
			_, _ = fmt.Fprintf(buf, "GitBranch    %s\n", branch())
			_, _ = fmt.Fprintf(buf, "GitHash      %s @ %s\n", hash(), timeString())
			_, _ = fmt.Fprintf(buf, "GitRepo      %s\n", repo())
			_, _ = fmt.Fprintf(buf, "\n")
		}
		_, _ = fmt.Fprintf(buf, "BuildHash    %s\n", buildHash())
		_, _ = fmt.Fprintf(buf, "BuildInfo    go-%s-%s @ %s\n",
			versionGo(), arch(), buildTimeString(),
		)
		fmt.Println(buf.String())
	},
}

const unavailable = "N/A"

var (
	appname   string
	version   string
	goVersion string

	gitRepo         string
	gitBranch       string
	gitHash         string
	gitNumber       string
	gitStatusNumber string
	gitStatusHash   string

	buildRand      string
	buildIndicator string
	buildTime      string
)

func name() string {
	if len(appname) > 0 {
		return appname
	}
	return unavailable
}

func versionInfo() string {
	if len(version) > 0 {
		if version[0] != 'v' {
			return version
		}
		if len(version) > 1 {
			return version[1:]
		}
	}
	if len(appname) > 0 {
		return "0.0.0"
	}
	return unavailable
}

func repo() string {
	if len(gitHash) >= 40 && len(gitRepo) > 0 {
		return gitRepo
	}
	return unavailable
}

func branch() string {
	if len(gitHash) >= 40 && len(gitBranch) > 0 {
		return gitBranch
	}
	return unavailable
}

func trace() string {
	if number() <= 0 {
		return unavailable
	}
	var devFlag string
	if len(gitStatusNumber) > 0 && len(gitStatusHash) >= 40 &&
		gitStatusHash != "da39a3ee5e6b4b0d3255bfef95601890afd80709" { // empty string sha1sum
		devFlag = fmt.Sprintf(" # %s.%s", gitStatusNumber, gitStatusHash[:7])
	}
	return fmt.Sprintf("%d.%s%s", number(), shortHash(), devFlag)
}

func hash() string {
	if len(gitHash) >= 40 {
		return gitHash[:40]
	}
	return unavailable
}

func timeString() string {
	if len(gitHash) <= 41 {
		return unavailable
	}
	t1, err := strconv.ParseInt(gitHash[41:], 10, 64)
	if err != nil {
		return unavailable
	}
	t := time.Unix(t1, 0)
	if !t.IsZero() {
		return t.Format("2006-01-02 15:04:05")
	}
	return unavailable
}

func shortHash() string {
	if len(gitHash) >= 40 {
		return gitHash[:7]
	}
	return unavailable
}

func number() uint64 {
	n, err := strconv.ParseUint(gitNumber, 10, 64)
	if err == nil {
		return n
	}
	return 0
}

func versionGo() string {
	tmp := strings.Split(goVersion, " ")
	if len(tmp) == 4 && len(tmp[2]) > 2 {
		return tmp[2][2:]
	}
	return unavailable
}

func arch() string {
	tmp := strings.Split(goVersion, " ")
	if len(tmp) == 4 {
		return tmp[3]
	}
	return unavailable
}

func buildHash() string {
	if len(buildRand) <= 0 {
		return unavailable
	}
	raw := strings.Join([]string{
		appname,
		version,
		goVersion,

		gitRepo,
		gitBranch,
		gitHash,
		gitNumber,
		gitStatusNumber,
		gitStatusHash,

		buildRand,
		buildIndicator,
		buildTime,
	}, "\x00")
	return fmt.Sprintf("%x", sha1.Sum([]byte(raw)))
}

func buildTimeString() string {
	t, err := strconv.ParseInt(buildTime, 10, 64)
	if err != nil {
		return unavailable
	}
	_t := time.Unix(t, 0)
	if !_t.IsZero() {
		return _t.Format("2006-01-02 15:04:05")
	}
	return unavailable
}
