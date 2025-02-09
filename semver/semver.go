//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/
//
// manage semantic versions for release.

package semver

import (
	"embed"
	"fmt"
)

// linker will override this variable. We parse it at run time
// See the semver package readme for details.
var SemVerInfo = LinkerDataDefault

const LinkerDataDefault = "hash|date|suffix|app|title"

// logic to valid the loading of the Info struct & linker data
func Initialise(fs embed.FS, filePath string) error {
	if err := getEmbeddedHistory(fs, filePath); err != nil {
		return err
	}

	if err := cleanLinkerData(); err != nil {
		return err
	}

	// set up the Short & Long responses from the components
	inf.Short = inf.Version + inf.SuffixShort

	//see https://semver.org/
	inf.Long = fmt.Sprintf("%s%s (%s:%s:%s:%s:%s)",
		inf.Version,
		inf.SuffixLong,
		inf.CodeName,
		inf.Date,
		inf.OS,
		inf.ARCH,
		inf.Note)
	return nil
}

func History() []ReleaseHistory {
	return history
}

func Info() VersionInfo {
	return inf
}
