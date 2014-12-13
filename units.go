package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var cmdSubmit = &Command{
	Run:      runSubmit,
	Usage:    "submit <unit>",
	NeedsCfg: true,
	Category: "units",
	Short:    "submit unit to cluster",
	Long: `
Submits a unit defined in the config file (clu.yaml) to the cluster.

Examples:

	$ clu submit logspout-service
`,
}

func runSubmit(cmd *Command, args []string) {
	if len(args) != 1 {
		cmd.PrintUsage()
		os.Exit(2)
	}

	name := args[0]

	unit, ok := cfg.Units[name]
	if !ok {
		printFatal(fmt.Sprintf("could not find unit by key %s", name))
	}

	d, err := ioutil.TempDir("", "clu")
	if err != nil {
		printFatal(err.Error())
	}

	err = ioutil.WriteFile(path.Join(d, unit.Name), []byte(unit.Content), 0644)
	if err != nil {
		printFatal(err.Error())
	}

	c := fmt.Sprintf("fleetctl submit %s", path.Join(d, unit.Name))
	fmt.Println("+ " + c)
	_, e, err := runCommand(c)
	if err != nil {
		printFatal(e)
	}
	err = os.RemoveAll(d)
	if err != nil {
		printFatal(err.Error())
	}
}
