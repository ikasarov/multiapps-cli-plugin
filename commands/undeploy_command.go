package commands

import (
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"

	baseclient "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
)

//UndeployCommand is a command for undeploying MTAs
type UndeployCommand struct {
	BaseCommand
	processTypeProvider ProcessTypeProvider
}

func NewUndeployCommand() *UndeployCommand {
	return &UndeployCommand{BaseCommand: BaseCommand{}, processTypeProvider: &undeployCommandProcessTypeProvider{}}
}

// GetPluginCommand returns the plugin command details
func (c *UndeployCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "undeploy",
		HelpText: "Undeploy a multi-target app",
		UsageDetails: plugin.Usage{
			Usage: `Undeploy a multi-target app
   cf undeploy MTA_ID [-u URL] [-f] [--retries RETRIES] [--namespace NAMESPACE] [--delete-services] [--delete-service-brokers] [--no-restart-subscribed-apps] [--do-not-fail-on-missing-permissions] [--abort-on-error]

   Perform action on an active undeploy operation
   cf undeploy -i OPERATION_ID -a ACTION [-u URL]`,
			Options: map[string]string{
				deployServiceURLOpt:                    "Deploy service URL, by default 'deploy-service.<system-domain>'",
				operationIDOpt:                         "Active undeploy operation id",
				actionOpt:                              "Action to perform on the active undeploy operation (abort, retry, monitor)",
				forceOpt:                               "Force undeploy without confirmation",
				util.GetShortOption(deleteServicesOpt): "Delete services",
				util.GetShortOption(deleteServiceBrokersOpt):       "Delete service brokers",
				util.GetShortOption(noRestartSubscribedAppsOpt):    "Do not restart subscribed apps, updated during the undeployment",
				util.GetShortOption(noFailOnMissingPermissionsOpt): "Do not fail on missing permissions for admin operations",
				util.GetShortOption(abortOnErrorOpt):               "Auto-abort the process on any errors",
				util.GetShortOption(retriesOpt):                    "Retry the operation N times in case a non-content error occurs (default 3)",
				// TODO: find out what's going on here??
				util.GetShortOption(namespaceOpt): "Specify the (optional) namespace the target mta is in",
			},
		},
	}
}

// Execute executes the command
func (c *UndeployCommand) Execute(args []string) ExecutionStatus {
	log.Tracef("Executing command '"+c.name+"': args: '%v'\n", args)

	var host string
	var operationID string
	var namespace string
	var actionID string
	var force bool
	var deleteServices bool
	var noRestartSubscribedApps bool
	var deleteServiceBrokers bool
	var noFailOnMissingPermissions bool
	var abortOnError bool
	var retries uint
	flags, err := c.CreateFlags(&host, args)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}
	flags.BoolVar(&force, forceOpt, false, "")
	flags.StringVar(&operationID, operationIDOpt, "", "")
	flags.StringVar(&namespace, namespaceOpt, "", "")
	flags.StringVar(&actionID, actionOpt, "", "")
	flags.BoolVar(&deleteServices, deleteServicesOpt, false, "")
	flags.BoolVar(&noRestartSubscribedApps, noRestartSubscribedAppsOpt, false, "")
	flags.BoolVar(&deleteServiceBrokers, deleteServiceBrokersOpt, false, "")
	flags.BoolVar(&noFailOnMissingPermissions, noFailOnMissingPermissionsOpt, false, "")
	flags.BoolVar(&abortOnError, abortOnErrorOpt, false, "")
	flags.UintVar(&retries, retriesOpt, 3, "")

	parser := NewCommandFlagsParser(flags, NewProcessActionExecutorCommandArgumentsParser([]string{"MTA_ID"}), NewDefaultCommandFlagsValidator(nil))
	err = parser.Parse(args)
	if err != nil {
		c.Usage(err.Error())
		return Failure
	}

	context, err := c.GetContext()
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	if operationID != "" || actionID != "" {
		return c.ExecuteAction(operationID, actionID, retries, host)
	}

	mtaID := args[0]
	if !force && !ui.Confirm("Really undeploy multi-target app %s? (y/n)", terminal.EntityNameColor(mtaID)) {
		ui.Warn("Undeploy cancelled")
		return Failure
	}

	// Print initial message
	ui.Say("Undeploying multi-target app %s in org %s / space %s as %s...",
		terminal.EntityNameColor(mtaID), terminal.EntityNameColor(context.Org),
		terminal.EntityNameColor(context.Space), terminal.EntityNameColor(context.Username))

	// Create rest client
	mtaClient, err := c.NewMtaClient(host)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}

	// Check if a deployed MTA with the specified ID exists
	_, err = mtaClient.GetMta(mtaID)
	if err != nil {
		ce, ok := err.(*baseclient.ClientError)
		if ok && ce.Code == 404 && strings.Contains(fmt.Sprint(ce.Description), mtaID) {
			ui.Failed("Multi-target app %s not found", terminal.EntityNameColor(mtaID))
			return Failure
		}
		ui.Failed("Could not get multi-target app %s: %s", terminal.EntityNameColor(mtaID), baseclient.NewClientError(err))
		return Failure

	}

	// Check for an ongoing operation for this MTA ID and abort it
	wasAborted, err := c.CheckOngoingOperation(mtaID, namespace, host, force)
	if err != nil {
		ui.Failed(err.Error())
		return Failure
	}
	if !wasAborted {
		return Failure
	}

	processBuilder := util.NewProcessBuilder()
	processBuilder.ProcessType(c.processTypeProvider.GetProcessType())
	processBuilder.Parameter("mtaId", mtaID)
	processBuilder.Parameter("noRestartSubscribedApps", strconv.FormatBool(noRestartSubscribedApps))
	processBuilder.Parameter("deleteServices", strconv.FormatBool(deleteServices))
	processBuilder.Parameter("deleteServiceBrokers", strconv.FormatBool(deleteServiceBrokers))
	processBuilder.Parameter("noFailOnMissingPermissions", strconv.FormatBool(noFailOnMissingPermissions))
	processBuilder.Parameter("abortOnError", strconv.FormatBool(abortOnError))
	processBuilder.Parameter("namespace", namespace)
	operation := processBuilder.Build()

	// Create the new process
	responseHeader, err := mtaClient.StartMtaOperation(*operation)
	if err != nil {
		ui.Failed("Could not create undeploy process: %s", err)
		return Failure
	}

	// Monitor process execution
	return NewExecutionMonitorFromLocationHeader(c.name, responseHeader.Location.String(), retries, []*models.Message{}, mtaClient).Monitor()
}

type undeployCommandProcessTypeProvider struct{}

func (d undeployCommandProcessTypeProvider) GetProcessType() string {
	return "UNDEPLOY"
}

type ProcessActionExecutorCommandArgumentsParser struct {
	positionalArgNames []string
}

func NewProcessActionExecutorCommandArgumentsParser(positionalArgNames []string) ProcessActionExecutorCommandArgumentsParser {
	return ProcessActionExecutorCommandArgumentsParser{positionalArgNames: positionalArgNames}
}

func (p ProcessActionExecutorCommandArgumentsParser) ParseFlags(flags *flag.FlagSet, args []string) error {
	operationExecutorOptions := make(map[string]string)
	for _, arg := range args {
		optionFlag := flags.Lookup(strings.Replace(arg, "-", "", 1))
		if optionFlag != nil && (operationIDOpt == optionFlag.Name || actionOpt == optionFlag.Name) {
			operationExecutorOptions[optionFlag.Name] = arg
		}
	}

	if len(operationExecutorOptions) > 2 {
		return fmt.Errorf("Options %s and %s should be specified only once", operationIDOpt, actionOpt)
	}

	if len(operationExecutorOptions) > 0 && len(operationExecutorOptions) < len([]string{operationIDOpt, actionOpt}) {
		var keys []string
		for _, key := range []string{operationIDOpt, actionOpt} {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return fmt.Errorf("All the %s options should be specified together", strings.Join(keys, " "))
	}

	return NewDefaultCommandFlagsParser(p.determinePositionalArguments(operationExecutorOptions)).ParseFlags(flags, args)
}

func (p ProcessActionExecutorCommandArgumentsParser) determinePositionalArguments(operationExecutorOptions map[string]string) []string {
	if len(operationExecutorOptions) == 2 {
		return []string{}
	}

	return p.positionalArgNames
}
