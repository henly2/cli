package cli_helper

import (
	"github.com/henly2/cli"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"github.com/fatih/color"
)

type RunArgs func(args []string) int
type Function struct {
	// Help should return long-form help text that includes the command-line
	// usage, a brief few sentences explaining the function of the command,
	// and the complete list of flags the command accepts.
	Help_ string

	// Run should run the actual command with the given CLI instance and
	// command-line arguments. It should return the exit status when it is
	// finished.
	//
	// There are a handful of special exit codes this can return documented
	// above that change behavior.
	RunArgs_ RunArgs

	// Synopsis should return a one-line, short synopsis of the command.
	// This should be less than 50 characters ideally.
	Synopsis_ string
}

func (f *Function) Help() string {
	return f.Help_
}

func (f *Function) Run(args []string) int {
	if f.Synopsis_ != "" {
		needArgs := strings.Split(f.Synopsis_, " ")
		if len(needArgs) != len(args) {
			return cli.RunResultHelp
		}
	}
	return f.RunArgs_(args)
}

func (f *Function) Synopsis() string {
	return f.Synopsis_
}

var (
	app_ string
	version_ string

	//bui *cli.BasicUi
	//cui *cli.ColoredUi
	//pui *cli.PrefixedUi
	cpui *cli.ColoredPrefixedUi
	commandFactory map[string]cli.CommandFactory

	exit chan os.Signal
)

func Init(app, version string)  {
	app_ = app
	version_ = version

	exit = make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	//bui = &cli.BasicUi{
	//	Reader:      os.Stdin,
	//	Writer:      os.Stdout,
	//	ErrorWriter: os.Stdout,
	//}
	//cui = &cli.ColoredUi{
	//	OutputColor: cli.UiColorNone,
	//	InfoColor:cli.UiColor{int(color.FgHiCyan), true},
	//	ErrorColor:cli.UiColorRed,
	//	WarnColor:cli.UiColorYellow,
	//	Ui: bui,
	//}
	//pui = &cli.PrefixedUi{
	//	AskPrefix:"[none]",
	//	AskSecretPrefix :"[none]",
	//	OutputPrefix    :"",
	//	InfoPrefix      :"[info]",
	//	ErrorPrefix     :"[err]",
	//	WarnPrefix      :"[war]",
	//	Ui: cui,
	//}
	cpui = cli.NewColoredPrefixedUi()
	cpui.AskPrefix 			= "[none]"
	cpui.AskSecretPrefix 	= "[none]"

	cpui.InfoPrefix 		= "[info]"
	cpui.OutputPrefix 		= ""
	cpui.ErrorPrefix 		= "[error]"
	cpui.WarnPrefix 		= "[warning]"

	cpui.AskPrefixColorAttributes[0] = append(cpui.AskPrefixColorAttributes[0], color.BgBlue)
	cpui.AskPrefixColorAttributes[0] = append(cpui.AskPrefixColorAttributes[0], color.FgHiWhite)

	cpui.AskSecretPrefixColorAttributes[0] = append(cpui.AskSecretPrefixColorAttributes[0], color.BgBlue)
	cpui.AskSecretPrefixColorAttributes[0] = append(cpui.AskSecretPrefixColorAttributes[0], color.FgHiWhite)

	cpui.OutputPrefixColorAttributes[1] = append(cpui.OutputPrefixColorAttributes[1], color.Bold)
	cpui.OutputPrefixColorAttributes[1] = append(cpui.OutputPrefixColorAttributes[1], color.FgHiWhite)

	cpui.ErrorPrefixColorAttributes[0] = append(cpui.ErrorPrefixColorAttributes[0], color.BgRed)
	cpui.ErrorPrefixColorAttributes[0] = append(cpui.ErrorPrefixColorAttributes[0], color.FgHiWhite)

	cpui.WarnPrefixColorAttributes[0] = append(cpui.WarnPrefixColorAttributes[0], color.BgYellow)
	cpui.WarnPrefixColorAttributes[0] = append(cpui.WarnPrefixColorAttributes[0], color.FgHiWhite)

	commandFactory = make(map[string]cli.CommandFactory)
}

func AddCommand(name string, help string, synopsis string, run RunArgs) (error) {
	if _, ok := commandFactory [name]; ok{
		err := fmt.Errorf("Exist name '%s'", name)
		cpui.Error(err.Error())
		return err
	}
	c := &Function{
		Help_:help,
		Synopsis_:synopsis,
		RunArgs_:run,
	}

	commandFactory[name] = func()(cli.Command, error) {
		return c, nil
	}

	return nil
}

func GetInput(name string) string {
	input, _ := cpui.Ask("Input "+name+">")
	select{
	case <-exit :
		os.Exit(0)
	default:
	}
	return input
}

func GetInputSecret(name string) string {
	input, _ := cpui.AskSecret("Input "+name+">")
	select{
	case <-exit :
		os.Exit(0)
	default:
	}
	return input
}

func SwitchPrefix(prefix string) {
	cpui.AskPrefix = prefix
	cpui.AskSecretPrefix = prefix
}

//func GetBui() *cli.BasicUi {
//	return bui
//}
//
//func GetCui() *cli.ColoredUi {
//	return cui
//}
//
//func GetPui() *cli.PrefixedUi {
//	return pui
//}

func GetCPui() *cli.ColoredPrefixedUi {
	return cpui
}

func RunOnce(argv []string) (int, error) {
	c1 := cli.NewCLI(app_, version_)
	c1.Args = argv
	c1.Commands = commandFactory
	return c1.Run()
}

func Run()  {
	RunOnce([]string{""})

	for {
		input := GetInput("command")
		argv := strings.Split(input, " ")


		exitStatus, err := RunOnce(argv)
		if err != nil {
			cpui.Error(err.Error())
		}

		if exitStatus < 0 {
			os.Exit(0)
		}
	}
}