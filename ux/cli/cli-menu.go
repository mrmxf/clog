//  Copyright ©2018-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// Give user an interactive choice of commands from a given parent

package cli

import (
	"fmt"
	"log/slog"
	"runtime"

	"github.com/charmbracelet/huh"
	"github.com/mrmxf/clog/config"
	"github.com/mrmxf/clog/ux"
)

var parentMenu *ux.MenuForm

func theMenuWithName(name string) *ux.MenuForm {
	for _, menu := range parentMenu.Children {
		if name == menu.Name {
			return menu
		}
	}
	return nil
}

// use the huh library to show an elegant multi-select form
func ShowCliMenu(menu *ux.MenuForm) (*ux.MenuForm, error) {
	parentMenu = menu
	vStr := config.Cfg().GetString("ver")
	var selected string

	var opts []huh.Option[string]
	opts = append(opts, huh.NewOption("exit", "exit"))
	for _, item := range parentMenu.Children {
		h := huh.NewOption(item.Name, item.Name)
		opts = append(opts, h)
	}
	opts = append(opts, huh.NewOption("help", "help"))
	selector := huh.NewSelect[string]()
	selector.
		Title(fmt.Sprintf("%s (clog@%s)", parentMenu.Name, vStr)).
		Options(opts...).
		Value(&selected)

	theForm := huh.NewForm(huh.NewGroup(selector))
	err := theForm.Run()
	if err != nil {
		slog.Error("Error running user interface form " + err.Error())
		return nil, err
	}
	return theMenuWithName(selected), nil
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
