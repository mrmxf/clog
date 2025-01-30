//  Copyright ©2019-2024  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/

// Package crayon uses [https://github.com/fatih/color] to provide role base
// colors for logging and highlighting on the console:
//
//   - Builtin
//   - Command
//   - Debug
//   - Dim
//   - Error
//   - File
//   - Heading
//   - Info
//   - Success
//   - Text
//   - Url
//   - Warning
//   - Xit
//
// Roles are defined in [CrayonColors] with a typical usage initialised with
// [Color] and then assigning a shorthand for the few colors you want to use:
//
//	s:= ttycrayon.Color().Success
//	i:= ttycrayon.Color().Info
//	e:= ttycrayon.Color().Error
//	fmt.Printf("%s %s and %s", i("exit with"), s("Success"), e("error"))
//
// Color scheme can be exported to bash/zsh with [GetBashString] and you can
// visualise with [SampleColors].
package crayon

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const escape = "\x1b"

// CrayonColors is a struct of functions for marking up TTY text.
//
// The target application is making CLI more legible and the roles are based
// on the sort of things that clog does. If you've got a different application
// then fork this repo and define your own
type CrayonColors struct {
	Builtin func(a ...interface{}) string //a builtin like Core
	Command func(a ...interface{}) string //CLI command like godoc
	Debug     func(a ...interface{}) string //debug or de-emphasise something
	Dim     func(a ...interface{}) string //dim or de-emphasise something
	Error   func(a ...interface{}) string //error
	File    func(a ...interface{}) string //file or folder names
	Heading func(a ...interface{}) string //headings
	Info    func(a ...interface{}) string //information messages (not body text)
	Success func(a ...interface{}) string //success
	Text    func(a ...interface{}) string //plain text
	Url     func(a ...interface{}) string //URL  / Uri / links
	Warning func(a ...interface{}) string //Warning
	Xit     func(a ...interface{}) string //stop coloring (used for bash export)

	B func(a ...interface{}) string // shorthand for: Builtin
	C func(a ...interface{}) string // shorthand for: Command
	D func(a ...interface{}) string // shorthand for: Dim & Debug
	E func(a ...interface{}) string // shorthand for: Error
	F func(a ...interface{}) string // shorthand for: File
	H func(a ...interface{}) string // shorthand for: Heading
	I func(a ...interface{}) string // shorthand for: Info
	S func(a ...interface{}) string // shorthand for: Success
	T func(a ...interface{}) string // shorthand for: Text
	U func(a ...interface{}) string // shorthand for: Url
	W func(a ...interface{}) string // shorthand for: Warning
	X func(a ...interface{}) string // shorthand for: Xit

	Amd func(a ...interface{}) string // AMD highlighter
	Arm func(a ...interface{}) string // Arm highlighter
	Lnx func(a ...interface{}) string // Linux highlighter
	Mac func(a ...interface{}) string // Mac highlighter
	Wsm func(a ...interface{}) string // Wasm highlighter
	Win func(a ...interface{}) string // Win highlighter

	//The bg variants all have solid backgrounds
	Bbg func(a ...interface{}) string // background emphasis: Builtin
	Cbg func(a ...interface{}) string // background emphasis: Command
	Dbg func(a ...interface{}) string // background emphasis: Dim / Debug
	Ebg func(a ...interface{}) string // background emphasis: Error
	Fbg func(a ...interface{}) string // background emphasis: File
	Hbg func(a ...interface{}) string // background emphasis: Heading
	Ibg func(a ...interface{}) string // background emphasis: Info
	Sbg func(a ...interface{}) string // background emphasis: Success
	Tbg func(a ...interface{}) string // background emphasis: Text
	Ubg func(a ...interface{}) string // background emphasis: Url
	Wbg func(a ...interface{}) string // background emphasis: Warning
	Xbg func(a ...interface{}) string // background emphasis: Xit
}

var crayonSprint CrayonColors

// ansiExit returns a Sprint() function that prepends the noColor escape seq.
func ansiExit() func(a ...interface{}) string {
	return func(a ...interface{}) string {
		return escape + "[0m" + fmt.Sprint(a...)
	}
}

// return a structure for coloring the ansi output - light Mode.
func Color() *CrayonColors {
	//enable color all the time
	color.NoColor = false

	builtinPlain := color.New(color.FgCyan).Add(color.Bold)
	builtinBlock := color.New(color.BgCyan).Add(color.FgHiYellow).Add(color.Bold)

	commandPlain := color.New(color.FgBlue)
	commandBlock := color.New(color.BgBlue).Add(color.FgYellow)

	dimPlain := color.New(color.FgWhite)
	dimBlock := color.New(color.BgWhite).Add(color.FgBlack)

	errorPlain := color.New(color.FgHiRed)
	errorBlock := color.New(color.BgHiRed).Add(color.FgWhite)

	filePlain := color.New(color.FgYellow)
	fileBlock := color.New(color.BgYellow).Add(color.FgBlack)

	headingPlain := color.New(color.FgCyan).Add(color.Bold)
	headingBlock := color.New(color.BgGreen).Add(color.FgBlack).Add(color.Bold)

	infoPlain := color.New(color.FgHiYellow)
	infoBlock := color.New(color.BgHiYellow).Add(color.FgBlue)

	successPlain := color.New(color.FgGreen)
	successBlock := color.New(color.BgGreen).Add(color.FgHiYellow)

	textPlain := color.New(color.FgBlack)
	textBlock := color.New(color.BgBlack).Add(color.FgHiWhite)

	urlPlain := color.New(color.FgCyan)
	urlBlock := color.New(color.FgCyan).Add(color.BgCyan)

	warningPlain := color.New(color.FgMagenta)
	warningBlock := color.New(color.FgMagenta).Add(color.BgMagenta)

	crayonSprint.Builtin = builtinPlain.SprintFunc()
	crayonSprint.Command = commandPlain.SprintFunc()
	crayonSprint.Dim = dimPlain.SprintFunc()
	crayonSprint.Error = errorPlain.SprintFunc()
	crayonSprint.File = filePlain.SprintFunc()
	crayonSprint.Heading = headingPlain.SprintFunc()
	crayonSprint.Info = infoPlain.SprintFunc()
	crayonSprint.Success = successPlain.SprintFunc()
	crayonSprint.Text = textPlain.SprintFunc()
	crayonSprint.Url = urlPlain.SprintFunc()
	crayonSprint.Warning = warningPlain.SprintFunc()
	crayonSprint.Xit = ansiExit()

	crayonSprint.B = crayonSprint.Builtin
	crayonSprint.C = crayonSprint.Command
	crayonSprint.D = crayonSprint.Dim
	crayonSprint.E = crayonSprint.Error
	crayonSprint.F = crayonSprint.File
	crayonSprint.H = crayonSprint.Heading
	crayonSprint.I = crayonSprint.Info
	crayonSprint.S = crayonSprint.Success
	crayonSprint.T = crayonSprint.Text
	crayonSprint.U = crayonSprint.Url
	crayonSprint.W = crayonSprint.Warning
	crayonSprint.X = ansiExit()

	crayonSprint.Amd = crayonSprint.Heading
	crayonSprint.Arm = crayonSprint.Success
	crayonSprint.Lnx = crayonSprint.Command
	crayonSprint.Mac = crayonSprint.Warning
	crayonSprint.Win = crayonSprint.Error

	crayonSprint.Bbg = builtinBlock.SprintFunc()
	crayonSprint.Cbg = commandBlock.SprintFunc()
	crayonSprint.Dbg = dimBlock.SprintFunc()
	crayonSprint.Ebg = errorBlock.SprintFunc()
	crayonSprint.Fbg = fileBlock.SprintFunc()
	crayonSprint.Hbg = headingBlock.SprintFunc()
	crayonSprint.Ibg = infoBlock.SprintFunc()
	crayonSprint.Sbg = successBlock.SprintFunc()
	crayonSprint.Tbg = textBlock.SprintFunc()
	crayonSprint.Ubg = urlBlock.SprintFunc()
	crayonSprint.Wbg = warningBlock.SprintFunc()
	crayonSprint.Xbg = ansiExit()

	return &crayonSprint
}

func SampleColors() string {
	c := Color()

	msg := ""
	msg = msg + "Builtin  " + c.B("Builtin") + "   " + c.Bbg("Builtin") + "\n"
	msg = msg + "Command  " + c.C("Command") + "   " + c.Cbg("Command") + "\n"
	msg = msg + "Dim      " + c.D("Dim    ") + "   " + c.Dbg("Dim    ") + "\n"
	msg = msg + "Error    " + c.E("Error  ") + "   " + c.Ebg("Error  ") + "\n"
	msg = msg + "File     " + c.F("File   ") + "   " + c.Fbg("File   ") + "\n"
	msg = msg + "Heading  " + c.H("Heading") + "   " + c.Hbg("Heading") + "\n"
	msg = msg + "Info     " + c.I("Info   ") + "   " + c.Ibg("Info   ") + "\n"
	msg = msg + "Success  " + c.S("Success") + "   " + c.Sbg("Success") + "\n"
	msg = msg + "Text     " + c.T("Text   ") + "   " + c.Tbg("Text   ") + "\n"
	msg = msg + "Url      " + c.U("Url    ") + "   " + c.Ubg("Url    ") + "\n"
	msg = msg + "Warning  " + c.W("Warning") + "   " + c.Wbg("Warning") + "\n"
	return msg
}

func toBashStr(bashVars []string, outputs []string) string {
	// start with the common escape root
	bashStr := ""
	bashEscape := "\\e"
	for i := range bashVars {
		slices := strings.Split(outputs[i], "XXX")
		bashCode := strings.ReplaceAll(slices[0], escape, bashEscape)
		bashStr = fmt.Sprintf("%s%s=\"%s\";", bashStr, bashVars[i], bashCode)
	}
	return bashStr
}

func GetBashString(darkMode bool) string {
	c := Color()
	if darkMode {
		c = Color()
	}
	x := "XXX"
	bashVars := []string{"cC", "cE", "cI", "cF", "cH", "cS", "cT", "cU", "cW", "cX", "cAmd", "cArm", "cLnx", "cMac", "cWin"}
	outputs := []string{c.C(x), c.E(x), c.I(x), c.F(x), c.H(x), c.S(x), c.T(x), c.U(x), c.W(x), c.X(x), c.Amd(x), c.Arm(x), c.Lnx(x), c.Mac(x), c.Win(x)}
	return toBashStr(bashVars, outputs)
}
