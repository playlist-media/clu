package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/playlist-media/clu/config"
)

var cmdLaunch = &Command{
	Run:      runLaunch,
	Usage:    "launch <name>",
	NeedsCfg: true,
	Category: "instances",
	Short:    "launch a named instance",
	Long: `
Launch an instance by name, based on the config file (clu.yaml)

Example:

	$ clu launch rinzler
`,
}

func runLaunch(cmd *Command, args []string) {
	if len(args) != 1 {
		cmd.PrintUsage()
		os.Exit(2)
	}
	name := args[0]
	instance, ok := cfg.Instances[name]
	if !ok {
		printFatal(fmt.Sprintf("Error: I don't know how to build %s - perhaps you should add it to config/clu.yaml?\n", name))
	}
	out := config.ProcessedCloudConfig(instance.CloudConfig)
	f, err := ioutil.TempFile("", "clu")
	defer f.Close()
	fmt.Fprint(f, out)
	f.Close()
	c := fmt.Sprintf("gcloud compute instances create %s ", name)
	c += "--image coreos-beta-444-3-0-v20141002 "
	c += "--image-project coreos-cloud "
	c += fmt.Sprintf("--machine-type %s ", instance.MachineType)
	c += fmt.Sprintf("--metadata-from-file user-data=%s ", f.Name())
	c += fmt.Sprintf("--project %s ", cfg.Global.ProjectID)
	c += fmt.Sprintf("--tags cluster %s ", instance.Kind)
	c += fmt.Sprintf("--zone %s ", cfg.Global.Zone)
	c += instance.MachineOpts
	_, _, err = runCommand(c)
	if err != nil {
		printFatal(err.Error())
	}
	err = os.Remove(f.Name())
	if err != nil {
		printFatal(err.Error())
	}
}

var cmdUpdateCloudConfig = &Command{
	Run:      runUpdateCloudConfig,
	Usage:    "update-cloud-config <name>",
	NeedsCfg: true,
	Category: "instances",
	Short:    "update the cloud-config for an instance",
	Long: `
Builds and updates and instance's cloud-config based on the config
file (clu.yaml)

Example:

	$ clu update rinzler
`,
}

func runUpdateCloudConfig(cmd *Command, args []string) {
	if len(args) != 1 {
		cmd.PrintUsage()
		os.Exit(2)
	}
	name := args[0]
	instance, ok := cfg.Instances[name]
	if !ok {
		printFatal(fmt.Sprintf("Error: I don't know how to update %s - perhaps you should add it to config/clu.yaml?\n", name))
	}

	out := config.ProcessedCloudConfig(instance.CloudConfig)
	f, err := ioutil.TempFile("", "clu")
	defer f.Close()
	fmt.Fprint(f, out)
	f.Close()
	c := fmt.Sprintf("gcloud compute instances add-metadata %s ", name)
	c += fmt.Sprintf("--metadata-from-file user-data=%s ", f.Name())
	c += fmt.Sprintf("--project %s ", cfg.Global.ProjectID)
	c += fmt.Sprintf("--zone %s ", cfg.Global.Zone)
	_, _, err = runCommand(c)
	if err != nil {
		printFatal(err.Error())
	}
	err = os.Remove(f.Name())
	if err != nil {
		printFatal(err.Error())
	}
}

var cmdRestart = &Command{
	Run:      runRestart,
	Usage:    "restart <name>",
	NeedsCfg: true,
	Category: "instances",
	Short:    "restart an instance",
	Long: `
Restarts an instance by name based on the config file (clu.yaml)

Example:

	$ clu restart rinzler
`,
}

func runRestart(cmd *Command, args []string) {
	if len(args) != 1 {
		cmd.PrintUsage()
		os.Exit(2)
	}
	name := args[0]
	_, ok := cfg.Instances[name]
	if !ok {
		printFatal(fmt.Sprintf("Error: I don't know how to restart %s - perhaps you should add it to config/clu.yaml?\n", name))
	}
	c := fmt.Sprintf("gcloud compute instances reset %s ", name)
	c += fmt.Sprintf("--project %s ", cfg.Global.ProjectID)
	c += fmt.Sprintf("--zone %s ", cfg.Global.Zone)
	_, _, err := runCommand(c)
	if err != nil {
		printFatal(err.Error())
	}
}
