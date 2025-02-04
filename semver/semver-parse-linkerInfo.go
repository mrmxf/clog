//This simple package manages the version number and name.
//
// semver.Info struct is exported for use in an application
//
// The ParseLinkerJson() function initialises the Info struct

package semver

import (
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// dummy linker string
const ( // iota is reset to 0
	lHASH     = iota
	lDATE     = iota
	lSUFFIX   = iota
	lAPPNAME  = iota
	lAPPTITLE = iota
)

// history is exported via a function
var history []ReleaseHistory // read from releases.yaml
var inf VersionInfo

// read the linker data and take appropriate cleaning actions
func cleanLinkerData() error {
	slog.Debug("Linker string is (" + SemVerInfo + ")")

	defaultInfo := strings.Split(LinkerDataDefault, "|")
	linkerInfo := strings.Split(SemVerInfo, "|")
	//slog.Debug(" linkerInfo is ", "array", linkerInfo)
	//slog.Debug("defaultInfo is ", "array", defaultInfo)

	if len(linkerInfo) != len(defaultInfo) {
		msg := fmt.Sprintf("ldflags SemVerInfo string should have %v fragments,, %v found", len(defaultInfo), len(linkerInfo))
		return errors.New(msg)
	}

	// ---commit hash -----------------------------------------------------------
	bashHash := "$(git rev-list -1 HEAD)"

	if len(linkerInfo[lHASH]) == 0 {
		msg := fmt.Sprintf("ldflags %s string fragment is empty - use %s", defaultInfo[lHASH], bashHash)
		return errors.New(msg)
	}

	if linkerInfo[lHASH] == defaultInfo[lHASH] {
		inf.CommitId = "xxxx^xxxx|xxxx^xxxx|xxxx^xxxx|xxxx^xxxx|"
	} else {
		inf.CommitId = linkerInfo[lHASH]
	}

	if len(inf.CommitId) < 40 {
		msg := fmt.Sprintf("ldflags %s string fragment should be 40 chars - use %s", defaultInfo[lHASH], bashHash)
		return errors.New(msg)
	}

	// --- date --- create automatically if empty string ------------------------
	now := time.Now().Format("2006-01-02")

	if len(linkerInfo[lDATE]) == 0 || linkerInfo[lDATE] == defaultInfo[lDATE] {
		inf.Date = now
	} else {
		inf.Date = linkerInfo[lDATE]
	}

	// --- app name -------------------------------------------------------------
	if len(linkerInfo[lAPPNAME]) == 0 || linkerInfo[lAPPNAME] == defaultInfo[lAPPNAME] {
		bi, ok := debug.ReadBuildInfo()
		if ok {
			inf.AppName = filepath.Base(bi.Main.Path) // name of the module
		}
	} else {
		inf.AppName = linkerInfo[lAPPNAME]
	}

	// --- app title-------------------------------------------------------------
	if len(linkerInfo[lAPPTITLE]) == 0 || linkerInfo[lAPPTITLE] == defaultInfo[lAPPTITLE] {
		bi, ok := debug.ReadBuildInfo()
		if ok {
			inf.AppTitle = filepath.Base(bi.Main.Path) // name of the module
		}
	} else {
		inf.AppTitle = linkerInfo[lAPPTITLE]
	}

	// --- suffix -------------------------------------------------------------
	suffix := linkerInfo[lSUFFIX]
	if linkerInfo[lSUFFIX] == defaultInfo[lSUFFIX] {
		suffix = "dev"
	}

	//replace underscores with spaces and beautify
	inf.AppTitle = strings.ReplaceAll(inf.AppTitle, "_", " ")
	inf.ARCH = runtime.GOARCH
	inf.OS = runtime.GOOS

	inf.Version = history[0].Version
	inf.CodeName = history[0].CodeName
	inf.Note = history[0].Note

	if len(suffix) > 0 {
		inf.SuffixShort = "-" + suffix
		inf.SuffixLong = "-" + suffix + "." + inf.CommitId[:4]
	} else {
		inf.SuffixShort = ""
		inf.SuffixLong = "+" + inf.CommitId[:4]
	}
	//slog.Debug("semver.Info is ", "struct", Info)
	return nil
}
