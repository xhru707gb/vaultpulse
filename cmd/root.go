package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpulse/internal/alert"
	"github.com/vaultpulse/internal/config"
	"github.com/vaultpulse/internal/expiry"
	"github.com/vaultpulse/internal/vault"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "vaultpulse",
	Short: "Monitor HashiCorp Vault secret expiration and rotation schedules",
	Long: `vaultpulse is a lightweight CLI tool that checks Vault secret
expiration dates and rotation schedules, then fires alerting hooks
when secrets are approaching expiry or have already expired.`,
	RunE: runCheck,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file (default: vaultpulse.yaml)")
	rootCmd.Flags().BoolP("table", "t", false, "render output as a formatted table")
}

func runCheck(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	checker := expiry.NewChecker(client, cfg)
	statuses, err := checker.Check()
	if err != nil {
		return fmt.Errorf("checking secrets: %w", err)
	}

	showTable, _ := cmd.Flags().GetBool("table")
	if showTable {
		fmt.Print(expiry.FormatTable(statuses))
	}

	notifier := alert.NewNotifier(cfg)
	if err := notifier.Notify(statuses); err != nil {
		return fmt.Errorf("sending alerts: %w", err)
	}

	return nil
}
