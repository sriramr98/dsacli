package seed

import (
	"dsacli/db"
	"dsacli/types"
	"encoding/json"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func GetCommand(db db.Database) *cobra.Command {
	return &cobra.Command{
		Use:   "seed [json_path]",
		Short: "Add problems to database",
		Long:  "Use this command to add problems to the database",
		Run:   runSeed(db),
		Args:  cobra.ExactArgs(1),
	}
}
func runSeed(db db.Database) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		executeSeed(db, args)
	}
}

func executeSeed(db db.Database, args []string) {
	problemFilePath := args[0]
	questions, err := readQuestions(problemFilePath)
	if err != nil {
		color.Red("Error reading questions from file: %s", err)
		return
	}

	color.Yellow("Inserting %d questions into database", len(questions))

	if err := db.InsertQuestions(questions); err != nil {
		color.Red("Error inserting questions into database: %s", err)
		return
	}

	color.Green("Successfully seeded %d questions into the database", len(questions))
}

func readQuestions(path string) ([]types.Question, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			color.Red("Unable to find file in path %s", path)
			return nil, err
		}
	}

	var questions []types.Question
	if err := json.Unmarshal(file, &questions); err != nil {
		color.Red("Unable to read questions from file with error: %s", err)
		return nil, err
	}

	return questions, nil
}
