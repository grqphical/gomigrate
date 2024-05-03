package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

const MigrationDir = "migrations"

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new migration up and migration down SQL file in ./migrations",
	Long:  ``,
	Run:   createCommand,
}

func createCommand(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: migration name needs to be provided")
		return
	}

	name := args[0]

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s.up.sql", timestamp, name)
	path := filepath.Join(MigrationDir, filename)

	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"ERROR: could not create up migration. Make sure you have run gomigrate init first.",
		)
		return
	}
	defer file.Close()

	path = fmt.Sprintf("%d_%s.down.sql", timestamp, name)
	file, err = os.Create(path)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"ERROR: could not create down migration. Make sure you have run gomigrate init first.",
		)
		return
	}
	defer file.Close()
	fmt.Printf("Created migration: %s\n", name)
}

func init() {
	rootCmd.AddCommand(createCmd)
}
