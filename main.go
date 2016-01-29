package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/cli/plugin/models"
)

/*
*	This is the struct implementing the interface defined by the core CLI. It can
*	be found at  "github.com/cloudfoundry/cli/plugin/plugin.go"
*
 */
type DoctorPlugin struct {
	ui terminal.UI
}

/*
*	This function must be implemented by any plugin because it is part of the
*	plugin interface defined by the core CLI.
*
*	Run(....) is the entry point when the core CLI is invoking a command defined
*	by a plugin. The first parameter, plugin.CliConnection, is a struct that can
*	be used to invoke cli commands. The second paramter, args, is a slice of
*	strings. args[0] will be the name of the command, and will be followed by
*	any additional arguments a cli user typed in.
*
*	Any error handling should be handled with the plugin itself (this means printing
*	user facing errors). The CLI will exit 0 if the plugin exits 0 and will exit
*	1 should the plugin exits nonzero.
 */
func (c *DoctorPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	fmt.Printf("\n\n")
	var triageApps []string
	var triageRoutes []string
	var triageServices []string
	var totalNumberOfApps int
	var totalNumberOfRunningApps int

	c.ui = terminal.NewUI(os.Stdin, terminal.NewTeePrinter())
	c.ui.Say(terminal.WarningColor("doctor: time to triage cloudfoundry"))
	fmt.Printf("\n")
	c.CFMainChecks(cliConnection)

	listOfRunningApps := c.AppsStateRunning(cliConnection)
	listOfStoppedApps := c.AppsStateStopped(cliConnection)
	triageApps = c.CheckUpApps(cliConnection, triageApps, listOfRunningApps, listOfStoppedApps)
	triageRoutes = c.CheckUpRoutes(cliConnection, triageRoutes)
	triageServices = c.CheckUpServices(cliConnection, triageServices)

	if len(triageApps) == 0 && len(triageRoutes) == 0 && len(triageServices) == 0 {
		c.ui.Say(terminal.SuccessColor("doctor: Everything looks OK!"))
		return
	}

	// doctor run results
	if len(triageApps) > 0 {
		c.ui.Say(terminal.WarningColor("Detected triage points for apps: "))
		for _, v := range triageApps {

			c.ui.Say(terminal.LogStderrColor(strings.Split(v, "___")[0]+" <---> ") + terminal.LogStderrColor(strings.Split(v, "___")[1]))
		}
	}
	c.ui.Say(" ")

	if len(triageRoutes) > 0 {
		c.ui.Say(terminal.WarningColor("Following routes do not have any app bound to them:"))

		for _, y := range triageRoutes {
			c.ui.Say(terminal.LogStderrColor(y))
		}
	}
	fmt.Printf("\n")
	if len(triageServices) > 0 {
		c.ui.Say(terminal.WarningColor("Following services do not have any app bound to them:"))

		for _, y := range triageServices {
			c.ui.Say(terminal.LogStderrColor(y))
		}
	}
	fmt.Printf("\n")
	results, err := cliConnection.GetApps()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	for _, app := range results {
		if app.State == "started" {
			totalNumberOfApps++
			totalNumberOfRunningApps++
		} else {
			totalNumberOfApps++
		}
	}

	c.ui.Say(terminal.WarningColor("Total Number of Apps: " + strconv.Itoa(totalNumberOfApps)))
	c.ui.Say(terminal.WarningColor("Total Number of Running Apps: " + strconv.Itoa(totalNumberOfRunningApps)))
}

// CheckUpRoutes performs checkup on currently defined routes in cloudfoundry
func (c *DoctorPlugin) CheckUpRoutes(cliConnection plugin.CliConnection, triageRoutes []string) []string {
	results, err := cliConnection.CliCommandWithoutTerminalOutput("routes")
	if err != nil {
		c.ui.Failed(err.Error())
	}

	for _, line := range results {
		// regex to match cf routes output and see if there are unbound routes
		match, _ := regexp.MatchString("^[a-zA-Z]*\\s*\\S*\\s*\\S*\\s*", line)
		if match {
			parts := strings.Fields(line)

			if len(parts) == 3 {
				triageRoutes = append(triageRoutes, "Host: "+parts[1]+" <--->  Domain: "+parts[2])
			}
		}

	}
	return triageRoutes
}

// CheckUpServices performs checkup on currently defined services in cloudfoundry
func (c *DoctorPlugin) CheckUpServices(cliConnection plugin.CliConnection, triageServices []string) []string {
	results, err := cliConnection.CliCommandWithoutTerminalOutput("services")
	if err != nil {
		c.ui.Failed(err.Error())
	}

	for _, line := range results {
		// regex to match cf services output and see if there are unbound services

		match, _ := regexp.MatchString("^\\S*\\s*\\S*\\s*\\D*\\s*create succeeded", line)
		if match {
			parts := strings.Fields(line)

			if len(parts) == 5 || len(parts) == 7 {
				triageServices = append(triageServices, "Service: "+parts[1]+" <--->  Name: "+parts[0])
			}

		}

	}
	return triageServices
}

// CheckUpApps performs checkup on applications and adds the result to triage map
func (c *DoctorPlugin) CheckUpApps(cliConnection plugin.CliConnection, triage []string, listOfRunningApps []plugin_models.GetAppsModel, listOfStoppedApps []plugin_models.GetAppsModel) []string {
	const alarmCPU float64 = 85.0

	for _, i := range listOfRunningApps {
		app, err := cliConnection.GetApp(i.Name)
		if err != nil {
			c.ui.Failed(err.Error())
		}

		if len(app.StagingFailedReason) > 0 {
			triage = append(triage, i.Name+" ___ StagingFailedReason: "+app.StagingFailedReason)
		}

		insts := app.Instances

		for _, ins := range insts {
			if ins.CpuUsage > alarmCPU {
				triage = append(triage, i.Name+" ___ CPU usage over %85 percent!")
			}

			if float64(ins.DiskUsage) > float64(ins.DiskQuota)*0.80 {
				triage = append(triage, i.Name+" ___ DiskUsage over %80 percent of DiskQuota")
			}

			if float64(ins.MemUsage) > float64(ins.MemQuota)*0.80 {
				triage = append(triage, i.Name+" ___ MemUsage over %80 percent of MemQuota")
			}

			if float64(ins.MemUsage) < float64(ins.MemQuota)*0.15 {
				triage = append(triage, i.Name+" ___ MemUsage lower than %15 percent of MemQuota, scaledown is an option.")
			}

			if len(insts) > 1 && float64(ins.MemUsage) < float64(ins.MemQuota)*0.15 && ins.CpuUsage < 10.0 {
				triage = append(triage, i.Name+" ___ app has more than one instance running with very low resource consumption. candidate for scaling down.")
			}

		}

		routes := app.Routes

		if len(routes) == 0 {
			triage = append(triage, i.Name+" ___ You have a running application that does not have a route!")
		}
	}

	for _, y := range listOfStoppedApps {
		app, err := cliConnection.GetApp(y.Name)
		if err != nil {
			c.ui.Failed(err.Error())
		}

		if len(app.StagingFailedReason) > 0 {
			triage = append(triage, y.Name+" ___ StagingFailedReason: "+app.StagingFailedReason)
		}
	}

	return triage

}

// AppsStateRunning will return a list of app whose state is running
func (c *DoctorPlugin) AppsStateRunning(cliConnection plugin.CliConnection) []plugin_models.GetAppsModel {
	var res []plugin_models.GetAppsModel
	appsListing, err := cliConnection.GetApps()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	for _, app := range appsListing {
		if app.State == "started" {
			res = append(res, app)
		}
	}
	return res
}

// AppsStateStopped will return a list of app whose state is running
func (c *DoctorPlugin) AppsStateStopped(cliConnection plugin.CliConnection) []plugin_models.GetAppsModel {
	var res []plugin_models.GetAppsModel
	appsListing, err := cliConnection.GetApps()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	for _, app := range appsListing {
		if app.State == "stopped" {
			res = append(res, app)
		}
	}
	return res
}

// CFMainChecks is responsible if the environment is okay for running doctor
func (c *DoctorPlugin) CFMainChecks(cliConnection plugin.CliConnection) {
	cliLogged, err := cliConnection.IsLoggedIn()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	cliHasOrg, err := cliConnection.HasOrganization()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	cliHasSpace, err := cliConnection.HasSpace()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	if cliLogged == false {
		panic("doctor cannot work without being logged in to CF")
	}

	if cliHasOrg == false || cliHasSpace == false {
		c.ui.Warn("WARN: It seems that your cloudfoundry has no space or org...")
	}
}

/*
*	This function must be implemented as part of the plugin interface
*	defined by the core CLI.
*
*	GetMetadata() returns a PluginMetadata struct. The first field, Name,
*	determines the name of the plugin which should generally be without spaces.
*	If there are spaces in the name a user will need to properly quote the name
*	during uninstall otherwise the name will be treated as seperate arguments.
*	The second value is a slice of Command structs. Our slice only contains one
*	Command Struct, but could contain any number of them. The first field Name
*	defines the command `cf basic-plugin-command` once installed into the CLI. The
*	second field, HelpText, is used by the core CLI to display help information
*	to the user in the core commands `cf help`, `cf`, or `cf -h`.
 */
func (c *DoctorPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "DoctorPlugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 2,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "doctor",
				HelpText: "doctor is responsible for scanning and reporting about anomalies present in your cloudfoundry.",
				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: "cf doctor\n",
				},
			},
		},
	}
}

/*
* Unlike most Go programs, the `Main()` function will not be used to run all of the
* commands provided in your plugin. Main will be used to initialize the plugin
* process, as well as any dependencies you might require for your
* plugin.
 */
func main() {
	plugin.Start(new(DoctorPlugin))
}

func (c *DoctorPlugin) showUsage(args []string) {
	for _, cmd := range c.GetMetadata().Commands {
		if cmd.Name == args[0] {
			fmt.Println("Invalid Usage: \n", cmd.UsageDetails.Usage)
		}
	}
}
