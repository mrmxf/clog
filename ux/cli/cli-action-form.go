//  Copyright ©2018-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// Give user an interactive choice of commands from a given parent

package cli

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/mrmxf/clog/ux"
)

// ActionForm shows rows of text and displays options.
//
// CLI mode works best with exactly 2 options. Web mode doesn't care
func ActionForm(formData *ux.ActionFormData) (int, error) {
	var group *huh.Group
	var opt0chosen bool
	var selected int
	if len(formData.Options) == 2 {
		chooser := huh.NewConfirm().
			Title(formData.Title).
			Description(strings.Join(formData.Description, "\n")).
			Affirmative(formData.Options[0].Caption).
			Negative(formData.Options[1].Caption).
			Value(&opt0chosen)
		group = huh.NewGroup(chooser)
	} else {

		var opts []huh.Option[int]
		for i, item := range formData.Options {
			h := huh.NewOption(item.Caption, i)
			opts = append(opts, h)
		}
		selector := huh.NewSelect[int]().
			Title(formData.Title).
			Description(strings.Join(formData.Description, "\n")).
			Options(opts...).
			Value(&selected)
		group = huh.NewGroup(selector)
	}

	// chain together group fields in a new huh form & run it
	form := huh.NewForm(group)
	err := form.Run()
	if err != nil {
		return 0, err
	}
	if len(formData.Options) != 2 {
		return selected, nil
	}
	if opt0chosen {
		return 0, nil
	}
	return 1, nil
}
