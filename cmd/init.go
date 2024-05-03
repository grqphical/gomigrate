package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/grqphical/gomigrate/internal/database"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initalizes a database for migrations",
	Long:  ``,
	Run:   initCommand,
}

func initCommand(cmd *cobra.Command, args []string) {
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

	err = database.CreateMigrationTable(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		return
	}

	err = os.MkdirAll("migrations", 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		return
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
}
