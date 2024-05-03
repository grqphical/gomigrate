package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/grqphical/gomigrate/internal/database"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   ListCommand,
}

func ListCommand(cmd *cobra.Command, args []string) {
	godotenv.Load("./.env")

	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		fmt.Println(
			"ERROR: DATABASE_URL env variable is not set. Make sure it is set or present in a '.env' file",
		)
		return
	}

	db, err := database.NewDatabaseConn(db_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		return
	}

	err = database.ListMigrations(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		return
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
