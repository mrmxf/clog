// Generic user interface - can be switched from CLI to WEB
package ui

import (
	"fmt"
	"github.com/mrmxf/clog/config"
	"github.com/mrmxf/clog/ux"
	"github.com/mrmxf/clog/ux/cli"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

var mode = ux.CLI

// return the start Menu item. This might not be the home menu e.g. the user
// types `clog Core` means that we should start at the Core submenu. The user
// can then navigate upwards if `Core` was not the menu intended. In golang,
// comparing pointers with `==` returns true if both pointers refer to the same
// object.
func findStartMenu(theCmd *cobra.Command) (*ux.MenuForm, error) {
	menu := ux.HomeMenu

	menu, err := matchMenuForm(theCmd, menu)
	if err == nil {
		return menu, nil
	} else {
		return nil, err
	}
}

func matchMenuForm(theCmd *cobra.Command, menu *ux.MenuForm) (*ux.MenuForm, error) {
	//if the desired command matches the current menu (e.g. home) then return
	cmdId := -1
	if menu.Cmd != nil && menu.Cmd.Annotations != nil {
		_, err := fmt.Sscanf(theCmd.Annotations["menu-id"], "%d", &cmdId)
		if err != nil {
			slog.Error("fatal error parsing menus in matchMenuForm")
			os.Exit(1)
		}
	}
	if cmdId == menu.Id {
		return menu, nil
	}
	for _, MenuForm := range menu.Children {
		// walk the tree depth first
		matchItem, err := matchMenuForm(theCmd, MenuForm)
		//something went wrong - propagate back up
		if err != nil {
			return nil, err
		}
		//got a match at a lower level - propagate up
		if matchItem != nil {
			return matchItem, nil
		}
		//test this MenuForm
		if MenuForm.Cmd == theCmd {
			return MenuForm, nil
		}
	}
	//no match and no errors
	return nil, nil
}

// show a menu and return the command that was run and the invoking choice
func HomeMenu(theCmd *cobra.Command, dummyList ...string) (*cobra.Command, error) {
	// if we got here then we are in interactive mode
	config.Cfg().Set("isInteractive", true)
	switch mode {
	case ux.CLI:
		startMenu, err := findStartMenu(theCmd)
		if err != nil {
			return nil, err
		}
		chosenMenu, err := cli.ShowCliMenu(startMenu)

		if chosenMenu == nil || chosenMenu.Cmd == nil {
			return nil, err
		}
		chosenMenu.Cmd.Run(chosenMenu.Cmd, []string{})

		//return the command that was run
		return chosenMenu.Cmd, nil
	}
	return nil, nil
}

// show a custom menu and return the "command" chosen child menu
func CustomMenu(rootMenu *ux.MenuForm) (*ux.MenuForm, error) {
	switch mode {
	case ux.CLI:
		chosenMenu, _ := cli.ShowCliMenu(rootMenu)
		return chosenMenu, nil
	}
	return nil, nil
}

// Action form displays a title, some text and a set of choices.
//
// the return value is the index of the choice list.
// If Options is nil then a no (0) and yes (1) option is synthesised
func ActionForm(formData *ux.ActionFormData) (int, error) {
	var noYes = []ux.Option{{Key: "n", Caption: "No"}, {Key: "y", Caption: "Yes"}}
	if formData.Options == nil {
		formData.Options = noYes
	}
	switch mode {
	case ux.CLI:
		chosen, _ := cli.ActionForm(formData)
		return chosen, nil
	}
	return 0, nil
}
