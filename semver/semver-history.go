//This simple package manages the version number and name.
//
// semver.Info struct is exported for use in an application
//
// The ParseLinkerJson() function initialises the Info struct

package semver

import (
	"embed"
	_ "embed"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

// read the history and return the latest version string
func getEmbeddedHistoryData(fs embed.FS, filePath string) error {
	yamlBytes, err := fs.ReadFile(filePath)
	if err != nil {
		e := fmt.Sprintf("Cannot read release history (%s)", filePath)
		return errors.New(e)
	}

	err = yaml.Unmarshal(yamlBytes, &Info.History)
	if err != nil {
		e := fmt.Sprintf("Cannot parse embedded history %v\n%v\n", filePath, err)
		return errors.New(e)
	}
	return nil
}
