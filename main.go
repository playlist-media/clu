package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	flag "github.com/bgentry/pflag"
	"github.com/bugsnag/bugsnag-go"
	"github.com/mgutz/ansi"
	"github.com/playlist-media/clu/config"
	"github.com/playlist-media/clu/term"
)

var (
	stdin = bufio.NewReader(os.Stdin)
)

type Command struct {
	Run      func(cmd *Command, args []string)
	Flag     flag.FlagSet
	NeedsCfg bool

	Usage    string // first word is the command name
	Category string // i.e. "App", "Instance", etc.
	Short    string // `clu help`
	Long     string // `clu help cmd`
}

func (c *Command) PrintUsage() {
	if c.Runnable() {
		fmt.Fprintf(os.Stderr, "Usage: clu %s\n", c.FullUsage())
	}
	fmt.Fprintf(os.Stderr, "Use 'clu help %s' for more information.\n", c.Name())
}

func (c *Command) PrintLongUsage() {
	if c.Runnable() {
		fmt.Printf("Usage: clu %s\n", c.FullUsage())
	}
	fmt.Println(strings.Trim(c.Long, "\n"))
}

func (c *Command) FullUsage() string {
	return c.Usage
}

func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

const extra = " (extra)"

func (c *Command) List() bool {
	return c.Short != "" && !strings.HasSuffix(c.Short, extra)
}

func (c *Command) ListAsExtra() bool {
	return c.Short != "" && strings.HasSuffix(c.Short, extra)
}

func (c *Command) ShortExtra() string {
	return c.Short[:len(c.Short)-len(extra)]
}

var commands = []*Command{
	cmdLaunch,
	cmdUpdateCloudConfig,
	cmdRestart,
	cmdSubmit,
	cmdVersion,
	cmdHelp,

	helpCommands,
	helpEnviron,
	helpMore,
	helpAbout,

	// unlisted
	cmdUpdate,
}

var (
	cfg       *config.Config
	cluAgent  = "hk/" + Version + " (" + runtime.GOOS + "; " + runtime.GOARCH + ")"
	userAgent = cluAgent
)

func main() {
	log.SetFlags(0)

	// ensure no global args, ensure command specified
	args := os.Args[1:]
	if len(args) < 1 || strings.IndexRune(args[0], '-') == 0 {
		printUsageTo(os.Stderr)
		os.Exit(2)
	}

	// run update command early
	if args[0] == cmdUpdate.Name() {
		cmdUpdate.Run(cmdUpdate, args)
		return
	} else if updater != nil {
		defer updater.backgroundRun() // will not run if os.Exit is called
	}

	if !term.IsANSI(os.Stdout) {
		ansi.DisableColors(true)
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			defer recoverPanic()

			cmd.Flag.SetDisableDuplicates(true) // disallow duplicate flag options
			if !gitConfigBool("clu.strict-flag-ordering") {
				cmd.Flag.SetInterspersed(true) // allow flags & non-flag args to mix
			}
			cmd.Flag.Usage = func() {
				cmd.PrintUsage()
			}
			if err := cmd.Flag.Parse(args[1:]); err == flag.ErrHelp {
				cmdHelp.Run(cmdHelp, args[:1])
				return
			} else if err != nil {
				printError(err.Error())
				os.Exit(2)
			}
			if cmd.NeedsCfg {
				getConfig()
			}
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
	if g := suggest(args[0]); len(g) > 0 {
		fmt.Fprintf(os.Stderr, "Possible alternatives: %v\n", strings.Join(g, " "))
	}
	fmt.Fprintf(os.Stderr, "Run 'hk help' for usage.\n")
	os.Exit(2)
}

var bugsnagClient = bugsnag.New(bugsnag.Configuration{
	APIKey:          "dd1201fbded336dd2103d1ba60dbacdd",
	AppVersion:      "1.2.3",
	Hostname:        "clu",
	ProjectPackages: []string{"main", "github.com/playlist-media/clu/*"},
})

func recoverPanic() {
	if Version != "dev" {
		if rec := recover(); rec != nil {
			message := ""
			switch rec := rec.(type) {
			case error:
				message = rec.Error()
			default:
				message = fmt.Sprintf("%v", rec)
			}
			if err := bugsnagClient.Notify(errors.New(message)); err != nil {
				printError("reporting crash failed: %s", err.Error())
				panic(rec)
			}
			printFatal("clu encountered and reported an internal client error")
		}
	}
}
