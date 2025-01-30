//  Copyright ©2019-2024  Mr MXF   info@mrmxf.com
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
var SemVerInfo = LinkerSemverDefault

// logic to valid the loading of the Info struct & linker data
func Initialise(fs embed.FS, filePath string) error {
	if err := getEmbeddedHistoryData(fs, filePath); err != nil {
		return err
	}

	if err := cleanLinkerData(); err != nil {
		return err
	}

	// set up the Short & Long responses from the components
	Info.Short = Info.Version + Info.SuffixShort

	//see https://semver.org/
	Info.Long = fmt.Sprintf("%s%s (%s:%s:%s:%s:%s)",
		Info.Version,
		Info.SuffixLong,
		Info.CodeName,
		Info.Date,
		Info.OS,
		Info.ARCH,
		"\""+Info.Note+"\"")
	return nil
}
