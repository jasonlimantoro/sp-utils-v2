package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"git.garena.com/jason.limantoro/shopee-utils-v2/cmd/createcard"
	"git.garena.com/jason.limantoro/shopee-utils-v2/cmd/createdraft"
	"git.garena.com/jason.limantoro/shopee-utils-v2/cmd/createlist"
	"git.garena.com/jason.limantoro/shopee-utils-v2/cmd/createmergerequest"
	"git.garena.com/jason.limantoro/shopee-utils-v2/cmd/getweeklyupdates"
	"git.garena.com/jason.limantoro/shopee-utils-v2/cmd/listmergerequest"
	"git.garena.com/jason.limantoro/shopee-utils-v2/cmd/reviewmergerequest"
	"git.garena.com/jason.limantoro/shopee-utils-v2/cmd/syncrepo"
	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/registry"
)

type Command struct {
	Name        string
	Description string
	SubCommands []Command
	Flags       []Flag
	Runner      Runner
}

type Flag struct {
	Name         string
	Description  string
	Shorthand    string
	DefaultValue string
	Required     bool
	Persistent   bool
	Env          string
}

type Runner interface {
	Run(ctx context.Context, flags map[string]string) error
}

func initCommand(diRegistry *registry.Registry) Command {
	var command = Command{
		Name:        "shopee-utils-v2",
		Description: "Utility Command line V2",
		SubCommands: []Command{
			{
				Name:        "merge-request",
				Description: "Merge Request Commands",
				SubCommands: []Command{
					{
						Name:        "create",
						Description: "Create MR",
						Flags: []Flag{
							{
								Name:         "repository",
								Description:  "repository to create the MR in",
								Shorthand:    "r",
								DefaultValue: "",
								Required:     true,
								Persistent:   false,
							},
							{
								Name:         "source-branch",
								Description:  "your feature branch",
								Shorthand:    "s",
								DefaultValue: "",
								Required:     true,
								Persistent:   false,
							},
							{
								Name:         "target-branch",
								Description:  "target branches, like uat or master",
								Shorthand:    "t",
								DefaultValue: "uat,master",
								Required:     false,
								Persistent:   false,
							},
							{
								Name:         "description",
								Description:  "description of MR",
								Shorthand:    "d",
								DefaultValue: "",
								Required:     true,
								Persistent:   false,
							},
							{
								Name:         "jira",
								Description:  "relevant jira tickets (e.g. SPOT-1234,SPOT-3245)",
								Shorthand:    "j",
								DefaultValue: "",
								Required:     true,
								Persistent:   false,
							},
						},
						Runner: createmergerequestcmd.NewRunner(diRegistry.CreateMergeRequestModule),
					},
					{
						Name:        "list",
						Description: "List MR",
						Flags: []Flag{
							{
								Name:         "repository",
								Description:  "repositories the MR(s) belong to",
								Shorthand:    "r",
								DefaultValue: "",
								Required:     true,
								Persistent:   false,
							},
							{
								Name:         "jira",
								Description:  "relevant jira tickets (e.g. SPOT-1234,SPOT-3245)",
								Shorthand:    "j",
								DefaultValue: "",
								Required:     false,
								Persistent:   false,
							},
							{
								Name:         "state",
								Description:  "state of merge requests (e.g., opened, closed, merged, locked)",
								Shorthand:    "s",
								DefaultValue: "opened",
								Required:     false,
								Persistent:   false,
							},
							{
								Name:         "search",
								Description:  "search merge request",
								Shorthand:    "",
								DefaultValue: "",
								Required:     false,
								Persistent:   false,
							},
						},
						Runner: listmergerequestcmd.NewRunner(diRegistry.ListMergeRequestModule),
					},
					{
						Name:        "review",
						Description: "Construct code review message",
						Flags: []Flag{
							{
								Name:         "repository",
								Description:  "repository the MR belongs to",
								Shorthand:    "r",
								DefaultValue: "",
								Required:     true,
								Persistent:   false,
							},
							{
								Name:         "jira",
								Description:  "relevant jira tickets (e.g. SPOT-1234,SPOT-3245)",
								Shorthand:    "j",
								DefaultValue: "",
								Required:     true,
								Persistent:   false,
							},
							{
								Name:         "template",
								Description:  "code review message template file path",
								Shorthand:    "t",
								DefaultValue: "",
								Required:     false,
								Persistent:   false,
							},
						},
						Runner: reviewmergerequestcmd.NewRunner(diRegistry.ReviewMergeRequestModule),
					},
				},
			},
			{
				Name:        "daily-updates",
				Description: "Daily Updates Commands",
				SubCommands: []Command{
					{
						Name:        "create-card",
						Description: "create a Trello card",
						SubCommands: nil,
						Flags: []Flag{
							{
								Name:        "title",
								Description: "title of the card",
								Shorthand:   "t",
								Required:    true,
							},
							{
								Name:         "list-name",
								Description:  "list name the card will be created in",
								Shorthand:    "l",
								DefaultValue: time.Now().Format("02-Jan-2006"),
								Required:     false,
							},
							{
								Name:        "jira-link",
								Description: "related jira ticket link (e.g. https://jira.shopee.io/browse/SPOT-1234)",
								Shorthand:   "j",
								Required:    false,
							},
							{
								Name:        "epic-link",
								Description: "epic ticket link (e.g. https://jira.shopee.io/browse/SPOT-1234)",
								Shorthand:   "e",
								Required:    false,
							},
							{
								Name:        "td-link",
								Description: "TD link",
								Required:    false,
							},
							{
								Name:        "prd-link",
								Description: "PRD link",
								Required:    false,
							},
						},
						Runner: createcardcmd.NewRunner(diRegistry.CreateCardModule),
					},
					{
						Name:        "create-list",
						Description: "create a Trello list for next working day",
						Flags: []Flag{
							{
								Name:        "operation-type",
								Description: "0: today, 1: next working day",
								Required:    true,
							},
						},
						Runner: createlistcmd.NewRunner(diRegistry.CreateListModule),
					},
					{
						Name:        "get",
						Description: "Get weekly updates",
						SubCommands: nil,
						Flags: []Flag{
							{
								Name:        "delta-week",
								Description: "Delta week from now (0: current week, 1: last week)",
								Shorthand:   "",
								Required:    true,
								Persistent:  false,
							},
							{
								Name:        "template",
								Description: "template file path",
								Shorthand:   "t",
								Required:    false,
								Persistent:  false,
							},
							{
								Name:        "out",
								Description: "output file",
								Shorthand:   "o",
								Required:    false,
								Persistent:  false,
							},
						},
						Runner: getweeklyupdatescmd.NewRunner(diRegistry.GetWeeklyUpdates),
					},
					{
						Name:        "create-draft",
						Description: "Create Gmail Draft",
						Flags: []Flag{
							{
								Name:         "input-file",
								Description:  "Input markdown file",
								Shorthand:    "i",
								DefaultValue: "",
								Required:     true,
							},
							{
								Name:         "delta-week",
								Description:  "Delta week from now (0: current week, 1: last week)",
								Shorthand:    "",
								DefaultValue: "",
								Required:     true,
							},
						},
						Runner: createdraftcmd.NewRunner(diRegistry.CreateDraftModule),
					},
				},
			},
			{
				Name:        "repo",
				Description: "Repo commands",
				SubCommands: []Command{
					{
						Name:        "sync",
						Description: "sync specified branches against remote",
						SubCommands: nil,
						Flags: []Flag{
							{
								Name:        "repo-file",
								Description: "csv file containing repository path and target branches",
								Shorthand:   "",
								Required:    false,
								Persistent:  false,
								Env:         "SYNC_REPO_FILE",
							},
						},
						Runner: syncrepocmd.NewRunner(diRegistry.SyncRepoModule),
					},
				},
			},
		},
	}

	return command
}

func Execute() {
	diRegistry := registry.InitRegistry()
	command := initCommand(diRegistry)
	rootCobraCmd := initCobra(command)

	if err := rootCobraCmd.Execute(); err != nil {
		panic(err)
	}
}

func initCobra(command Command) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   command.Name,
		Short: command.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			flagsMap := make(map[string]string)

			for _, flag := range command.Flags {
				var flagValue string
				if flag.Persistent {
					flagValue = cmd.PersistentFlags().Lookup(flag.Name).Value.String()
				} else {
					flagValue = cmd.Flags().Lookup(flag.Name).Value.String()
				}

				if flagValue != "" {
					flagsMap[flag.Name] = flagValue
				} else {
					flagsMap[flag.Name] = viper.GetString(flag.Name)
				}
			}

			if command.Runner == nil {
				return cmd.Help()
			}

			return command.Runner.Run(context.TODO(), flagsMap)
		},
	}

	for _, flag := range command.Flags {
		defaultValue := flag.DefaultValue
		if flag.Env != "" {
			viper.BindEnv(flag.Name, flag.Env)
			defaultValue = viper.GetString(flag.Name)
		}

		if flag.Persistent {
			cobraCmd.PersistentFlags().StringP(
				flag.Name,
				flag.Shorthand,
				defaultValue,
				flag.Description,
			)
		} else {
			cobraCmd.Flags().StringP(
				flag.Name,
				flag.Shorthand,
				defaultValue,
				flag.Description,
			)
		}

		if flag.Required {
			if flag.Persistent {
				_ = cobraCmd.MarkPersistentFlagRequired(flag.Name)
			} else {
				_ = cobraCmd.MarkFlagRequired(flag.Name)
			}
		}
	}

	for _, subCommand := range command.SubCommands {
		subCommandCobra := initCobra(subCommand)
		cobraCmd.AddCommand(subCommandCobra)
	}

	return cobraCmd
}
