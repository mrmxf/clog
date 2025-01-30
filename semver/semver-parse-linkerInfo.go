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
const LinkerSemverDefault = "commithash|date|suffix|appname|apptitle"
const ( // iota is reset to 0
	lHASH     = iota
	lDATE     = iota
	lSUFFIX   = iota
	lAPPNAME  = iota
	lAPPTITLE = iota
)

// read the linker data and take appropriate actions
func cleanLinkerData() error {
	slog.Debug("Linker string is (" + SemVerInfo + ")")

	defaultInfo := strings.Split(LinkerSemverDefault, "|")
	linkerInfo := strings.Split(SemVerInfo, "|")
	// log.Debug(" linkerInfo is ", "array", linkerInfo)
	// log.Debug("defaultInfo is ", "array", defaultInfo)

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
		Info.CommitId = "xxxx^xxxx|xxxx^xxxx|xxxx^xxxx|xxxx^xxxx|"
	} else {
		Info.CommitId = linkerInfo[lHASH]
	}

	if len(Info.CommitId) < 40 {
		msg := fmt.Sprintf("ldflags %s string fragment should be 40 chars - use %s", defaultInfo[lHASH], bashHash)
		return errors.New(msg)
	}

	// --- date --- create automatically if empty string ------------------------
	now := time.Now().Format("2006-01-02")

	if len(linkerInfo[lDATE]) == 0 || linkerInfo[lDATE] == defaultInfo[lDATE] {
		Info.Date = now
	} else {
		Info.Date = linkerInfo[lDATE]
	}

	// --- app name -------------------------------------------------------------
	if len(linkerInfo[lAPPNAME]) == 0 || linkerInfo[lAPPNAME] == defaultInfo[lAPPNAME] {
		bi, ok := debug.ReadBuildInfo()
		if ok {
			Info.AppName = filepath.Base(bi.Main.Path) // name of the module
		}
	} else {
		Info.AppName = linkerInfo[lAPPNAME]
	}

	// --- app title-------------------------------------------------------------
	if len(linkerInfo[lAPPTITLE]) == 0 || linkerInfo[lAPPTITLE] == defaultInfo[lAPPTITLE] {
		bi, ok := debug.ReadBuildInfo()
		if ok {
			Info.AppTitle = filepath.Base(bi.Main.Path) // name of the module
		}
	} else {
		Info.AppTitle = linkerInfo[lAPPTITLE]
	}

	// --- suffix -------------------------------------------------------------
	suffix := linkerInfo[lSUFFIX]
	if linkerInfo[lSUFFIX] == defaultInfo[lSUFFIX] {
		suffix = "dev"
	}

	//replace underscores with spaces and beautify
	Info.AppTitle = highlightTitleCase(strings.ReplaceAll(Info.AppTitle, "_", " "))
	Info.ARCH = runtime.GOARCH
	Info.OS = runtime.GOOS

	Info.Version = Info.History[0].Version
	Info.CodeName = Info.History[0].CodeName
	Info.Note = Info.History[0].Note

	if len(suffix) > 0 {
		Info.SuffixShort = "-" + suffix
		Info.SuffixLong = "-" + suffix + "." + Info.CommitId[:4]
	} else {
		Info.SuffixShort = ""
		Info.SuffixLong = "+" + Info.CommitId[:4]
	}
	// log.Debug("semver.Info is ", "struct", Info)
	return nil
}
