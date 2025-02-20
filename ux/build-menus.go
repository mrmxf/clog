package ux

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var HomeMenu *MenuForm
var newMenuIndex = 100

// iterate through the children of the command and fill the array
func populateMenuLevel(cmd *cobra.Command, parent *MenuForm) {
	keyVocabulary := strings.Clone(AliasVocabulary)
	for _, childCmd := range cmd.Commands() {
		//create a new menuLevel for each
		child := &MenuForm{
			Parent: parent,
			Cmd:    childCmd,
			Name:   childCmd.Name(),
			Short:  childCmd.Short,
			Id:     newMenuIndex,
		}
		if childCmd.Annotations == nil {
			childCmd.Annotations = make(map[string]string)
		}
		childCmd.Annotations["menu-id"] = fmt.Sprintf("%d", newMenuIndex)
		newMenuIndex++
		if len(child.Name) > 0 {
			parent.Children = append(parent.Children, child)
		}
	}
	//sort the children into alphabetical order
	sort.SliceStable(parent.Children, func(i, j int) bool {
		return parent.Children[i].Name < parent.Children[j].Name
	})

	//allocate a key & recurse for next level of children
	for _, menuItem := range parent.Children {
		for _, c := range menuItem.Name {
			pos := strings.Index(keyVocabulary, string(c))
			if pos >= 0 {
				menuItem.Key = string(c)
				//trim the vocabulary to prevent reuse
				keyVocabulary = strings.ReplaceAll(keyVocabulary, string(c), "")
				break
			}
		}
		populateMenuLevel(menuItem.Cmd, menuItem)
	}
}

// build the menu hierarchy from the root command
func BuildMenus(rootCommand *cobra.Command) {
	HomeMenu = &MenuForm{
		Parent: nil,
		Cmd:    rootCommand,
		Name:   "home",
		Key:    "/",
		Id:     newMenuIndex,
	}
	if rootCommand.Annotations == nil {
		rootCommand.Annotations = make(map[string]string)
	}
	rootCommand.Annotations["menu-id"] = fmt.Sprintf("%d", newMenuIndex)
	newMenuIndex++
	populateMenuLevel(rootCommand, HomeMenu)
}

// create and return a menu from a map[string]string
func MenuFromMap(cmd *cobra.Command, theMap map[string]string, helpMap map[string]string) *MenuForm {
	root := MenuForm{
		Name: cmd.Use,
		Cmd:  cmd,
	}

	for key := range theMap {
		help, present := helpMap[key]
		if !present {
			help = ""
		}
		menu := MenuForm{
			Name:  string(key),
			Key:   string(key[0]),
			Short: help,
			Cmd:   cmd,
		}
		root.Children = append(root.Children, &menu)
	}
	return &root
}
