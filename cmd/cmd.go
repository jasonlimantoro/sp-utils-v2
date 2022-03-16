package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"git.garena.com/jason.limantoro/shopee-utils-v2/cmd/createmergerequest"
	listmergerequestcmd "git.garena.com/jason.limantoro/shopee-utils-v2/cmd/listmergerequest"
	reviewmergerequestcmd "git.garena.com/jason.limantoro/shopee-utils-v2/cmd/reviewmergerequest"
	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/registry"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/reviewmergerequest"
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
								Description:  "repository to create the MR in",
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
								Name:         "state",
								Description:  "state of merge requests (e.g., opened, closed, merged, locked)",
								Shorthand:    "s",
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
								Description:  "repository to create the MR in",
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
								DefaultValue: reviewmergerequest.DefaultCodeReviewMessageTemplate,
								Required:     false,
								Persistent:   false,
							},
						},
						Runner: reviewmergerequestcmd.NewRunner(diRegistry.ReviewMergeRequestModule),
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
				if flag.Persistent {
					flagsMap[flag.Name] = cmd.PersistentFlags().Lookup(flag.Name).Value.String()
				} else {
					flagsMap[flag.Name] = cmd.Flags().Lookup(flag.Name).Value.String()
				}
			}

			if command.Runner == nil {
				return cmd.Help()
			}

			return command.Runner.Run(context.TODO(), flagsMap)
		},
	}

	for _, flag := range command.Flags {
		if flag.Persistent {
			cobraCmd.PersistentFlags().StringP(
				flag.Name,
				flag.Shorthand,
				flag.DefaultValue,
				flag.Description,
			)
		} else {
			cobraCmd.Flags().StringP(
				flag.Name,
				flag.Shorthand,
				flag.DefaultValue,
				flag.Description,
			)
		}

		if flag.Required {
			_ = cobraCmd.MarkFlagRequired(flag.Name)
		}
	}

	for _, subCommand := range command.SubCommands {
		subCommandCobra := initCobra(subCommand)
		cobraCmd.AddCommand(subCommandCobra)
	}

	return cobraCmd
}
