package service

import (
	"fmt"
	"strings"

	"github.com/nodeset-org/hyperdrive/shared/types"
	"github.com/urfave/cli/v2"

	"github.com/rocket-pool/smartnode/rocketpool-cli/utils"
	"github.com/rocket-pool/smartnode/rocketpool-cli/utils/terminal"
	"github.com/rocket-pool/smartnode/shared/config"
	cfgtypes "github.com/rocket-pool/smartnode/shared/types/config"
	"github.com/rocket-pool/smartnode/shared/utils/input"
)

var (
	ignoreSlashTimerFlag *cli.BoolFlag = &cli.BoolFlag{
		Name:  "ignore-slash-timer",
		Usage: "Bypass the safety timer that forces a delay when switching to a new ETH2 client",
	}
	tailFlag *cli.StringFlag = &cli.StringFlag{
		Name:    "tail",
		Aliases: []string{"t"},
		Usage:   "The number of lines to show from the end of the logs (number or \"all\")",
		Value:   "100",
	}
)

// Creates CLI argument flags from the parameters of the configuration struct
func createFlagsFromConfigParams(sectionName string, params []types.IParameter, configFlags []cli.Flag, network cfgtypes.Network) []cli.Flag {
	for _, param := range params {
		common := param.GetCommon()
		var paramName string
		if sectionName == "" {
			paramName = common.ID
		} else {
			paramName = fmt.Sprintf("%s-%s", sectionName, common.ID)
		}

		usage := common.Description
		options := param.GetOptions()
		if len(options) > 0 {
			optionStrings := []string{}
			for _, option := range options {
				optionStrings = append(optionStrings, fmt.Sprint(option.GetValueAsString()))
			}
			usage = fmt.Sprintf("%s\nOptions: %s\n", common.Description, strings.Join(optionStrings, ", "))
		}

		defaultVal := param.GetDefaultAsString(network)
		configFlags = append(configFlags, &cli.StringFlag{
			Name:  paramName,
			Usage: usage,
			Value: defaultVal,
		})
	}

	return configFlags
}

// Register commands
func RegisterCommands(app *cli.App, name string, aliases []string) {

	configFlags := []cli.Flag{}
	cfgTemplate := config.NewRocketPoolConfig("", false)
	network := cfgTemplate.Smartnode.Network.Value.(cfgtypes.Network)

	// Root params
	configFlags = createFlagsFromConfigParams("", cfgTemplate.GetParameters(), configFlags, network)

	// Subconfigs
	for sectionName, subconfig := range cfgTemplate.GetSubconfigs() {
		configFlags = createFlagsFromConfigParams(sectionName, subconfig.GetParameters(), configFlags, network)
	}

	app.Commands = append(app.Commands, &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage Hyperdrive service",
		Flags: []cli.Flag{
			utils.ComposeFileFlag,
		},
		Subcommands: []*cli.Command{
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "Install the Hyperdrive service",
				Flags: []cli.Flag{
					utils.YesFlag,
					installVerboseFlag,
					installNoDepsFlag,
					installPathFlag,
					installVersionFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return installService(c)
				},
			},

			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Configure the Rocket Pool service",
				Flags:   configFlags,
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return configureService(c)
				},
			},

			{
				Name:    "sync",
				Aliases: []string{"y"},
				Usage:   "Get the sync progress of the Execution and Consensus clients",
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return getSyncProgress(c)
				},
			},

			{
				Name:    "status",
				Aliases: []string{"u"},
				Usage:   "View the Rocket Pool service status",
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return serviceStatus(c)
				},
			},

			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "Start the Rocket Pool service",
				Flags: []cli.Flag{
					ignoreSlashTimerFlag,
					utils.YesFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return startService(c, false)
				},
			},

			{
				Name:      "stop",
				Aliases:   []string{"pause", "p"},
				Usage:     "Pause the Rocket Pool service",
				UsageText: "rocketpool service pause [options]",
				Flags: []cli.Flag{
					utils.YesFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return stopService(c)
				},
			},

			{
				Name:      "logs",
				Aliases:   []string{"l"},
				Usage:     "View the Rocket Pool service logs",
				ArgsUsage: "[service names]",
				Flags: []cli.Flag{
					tailFlag,
				},
				Action: func(c *cli.Context) error {
					// Run command
					return serviceLogs(c, c.Args().Slice()...)
				},
			},

			{
				Name:    "stats",
				Aliases: []string{"a"},
				Usage:   "View the Rocket Pool service stats",
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return serviceStats(c)
				},
			},

			{
				Name:  "compose",
				Usage: "View the Rocket Pool service docker compose config",
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return serviceCompose(c)
				},
			},

			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "View the Rocket Pool service version information",
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return serviceVersion(c)
				},
			},

			{
				Name:    "prune-ec",
				Aliases: []string{"prune-eth1", "n"},
				Usage:   "Shuts down the main ETH1 client and prunes its database, freeing up disk space, then restarts it when it's done.",
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return pruneExecutionClient(c)
				},
			},

			{
				Name:    "install-update-tracker",
				Aliases: []string{"d"},
				Usage:   "Install the update tracker that provides the available system update count to the metrics dashboard",
				Flags: []cli.Flag{
					utils.YesFlag,
					installUpdateTrackerVerboseFlag,
					installUpdateTrackerVersionFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return installUpdateTracker(c)
				},
			},

			{
				Name:    "check-cpu-features",
				Aliases: []string{"ccf"},
				Usage:   "Checks if your CPU supports all of the features required by the \"modern\" version of certain client images. If not, it prints what features are missing.",
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return checkCpuFeatures()
				},
			},

			{
				Name:  "get-config-yaml",
				Usage: "Generate YAML that shows the current configuration schema, including all of the parameters and their descriptions",
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return getConfigYaml(c)
				},
			},

			{
				Name:      "export-ec-data",
				Aliases:   []string{"export-eth1-data"},
				Usage:     "Exports the execution client (eth1) chain data to an external folder. Use this if you want to back up your chain data before switching execution clients.",
				ArgsUsage: "target-folder",
				Flags: []cli.Flag{
					exportEcDataForceFlag,
					exportEcDataDirtyFlag,
					utils.YesFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 1); err != nil {
						return err
					}
					targetDir := c.Args().Get(0)

					// Run command
					return exportEcData(c, targetDir)
				},
			},

			{
				Name:      "import-ec-data",
				Aliases:   []string{"import-eth1-data"},
				Usage:     "Imports execution client (eth1) chain data from an external folder. Use this if you want to restore the data from an execution client that you previously backed up.",
				ArgsUsage: "source-folder",
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 1); err != nil {
						return err
					}
					sourceDir := c.Args().Get(0)

					// Run command
					return importEcData(c, sourceDir)
				},
			},

			{
				Name:    "resync-ec",
				Aliases: []string{"resync-eth1"},
				Usage:   fmt.Sprintf("%sDeletes the main Execution client's chain data and resyncs it from scratch. Only use this as a last resort!%s", terminal.ColorRed, terminal.ColorReset),
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return resyncExecutionClient(c)
				},
			},

			{
				Name:    "resync-cc",
				Aliases: []string{"resync-eth2"},
				Usage:   fmt.Sprintf("%sDeletes the Consensus client's chain data and resyncs it from scratch. Only use this as a last resort!%s", terminal.ColorRed, terminal.ColorReset),
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return resyncConsensusClient(c)
				},
			},

			{
				Name:    "terminate",
				Aliases: []string{"t"},
				Usage:   fmt.Sprintf("%sDeletes all of the Rocket Pool Docker containers and volumes, including your ETH1 and ETH2 chain data and your Prometheus database (if metrics are enabled). Also removes your entire `.rocketpool` configuration folder, including your wallet, password, and validator keys. Only use this if you are cleaning up the Smartnode and want to start over!%s", terminal.ColorRed, terminal.ColorReset),
				Flags: []cli.Flag{
					utils.YesFlag,
				},
				Action: func(c *cli.Context) error {
					// Validate args
					if err := input.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run command
					return terminateService(c)
				},
			},
		},
	})
}
