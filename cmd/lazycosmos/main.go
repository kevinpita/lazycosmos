// Command lazycosmos starts the LazyCosmos CLI.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultGreeting = "Hello from LazyCosmos"
	defaultProfile  = "local"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	rootCmd, err := newRootCommand()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() (*cobra.Command, error) {
	var configFile string

	rootCmd := &cobra.Command{
		Use:           "lazycosmos",
		Short:         "Keyboard-first TUI for Cosmos SDK app development workflows",
		Version:       buildInfo(),
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return initConfig(configFile)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return runTUI()
		},
	}
	rootCmd.SetVersionTemplate("{{.Name}} {{.Version}}\n")

	flags := rootCmd.PersistentFlags()
	flags.StringVar(&configFile, "config", "", "config file")
	flags.String("profile", defaultProfile, "profile name")
	flags.String("greeting", defaultGreeting, "startup greeting")

	viper.SetDefault("profile", defaultProfile)
	viper.SetDefault("greeting", defaultGreeting)
	viper.SetEnvPrefix("lazycosmos")
	viper.AutomaticEnv()

	if err := viper.BindPFlag("profile", flags.Lookup("profile")); err != nil {
		return nil, fmt.Errorf("bind profile flag: %w", err)
	}
	if err := viper.BindPFlag("greeting", flags.Lookup("greeting")); err != nil {
		return nil, fmt.Errorf("bind greeting flag: %w", err)
	}

	rootCmd.AddCommand(newVersionCommand())

	return rootCmd, nil
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return printVersion(cmd.OutOrStdout())
		},
	}
}

func initConfig(configFile string) error {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("find home directory: %w", err)
		}

		viper.AddConfigPath(homeDir)
		viper.SetConfigName(".lazycosmos")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if configFile == "" && errors.As(err, &notFound) {
			return nil
		}

		return fmt.Errorf("read config: %w", err)
	}

	return nil
}

func runTUI() error {
	program := tea.NewProgram(appModel{
		greeting: viper.GetString("greeting"),
		profile:  viper.GetString("profile"),
	})

	if _, err := program.Run(); err != nil {
		return fmt.Errorf("run TUI: %w", err)
	}

	return nil
}

func buildInfo() string {
	return fmt.Sprintf("%s (commit %s, built %s)", version, commit, date)
}

func printVersion(writer io.Writer) error {
	if _, err := fmt.Fprintf(writer, "lazycosmos %s\n", buildInfo()); err != nil {
		return fmt.Errorf("print version: %w", err)
	}

	return nil
}

type appModel struct {
	greeting string
	profile  string
}

func (model appModel) Init() tea.Cmd {
	return nil
}

func (model appModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	if message, ok := message.(tea.KeyMsg); ok {
		switch message.String() {
		case "ctrl+c", "esc", "q":
			return model, tea.Quit
		}
	}

	return model, nil
}

func (model appModel) View() string {
	return fmt.Sprintf(
		"%s\n\nProfile: %s\n\nPress q to quit.\n",
		model.greeting,
		model.profile,
	)
}
