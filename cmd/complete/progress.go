package complete

import (
	"dsacli/db"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// GetProgressCommand returns a command to check progression gate status
func GetProgressCommand(database db.Database) *cobra.Command {
	return &cobra.Command{
		Use:   "progress",
		Short: "Check progression gate status for difficulty tiers",
		Long:  `Check the mastery percentage for each difficulty tier and see which tiers are unlocked.`,
		Run:   progressCmd(database),
	}
}

func progressCmd(database db.Database) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := executeProgress(database); err != nil {
			color.Red("Error: %v", err)
			return
		}
	}
}

func executeProgress(database db.Database) error {
	difficulties := []string{"easy", "medium", "hard"}

	color.Cyan("🎯 Progression Gate Status\n")

	for _, difficulty := range difficulties {
		questions, err := database.GetQuestionsByDifficulty(difficulty)
		if err != nil {
			return fmt.Errorf("failed to get %s questions: %w", difficulty, err)
		}

		if len(questions) == 0 {
			color.Yellow("No %s questions found", difficulty)
			continue
		}

		masteredCount := 0
		for _, q := range questions {
			if q.Mastered {
				masteredCount++
			}
		}

		masteryPercentage := float64(masteredCount) / float64(len(questions)) * 100
		isUnlocked := masteryPercentage > 50.0

		statusIcon := "🔒"
		statusColor := color.Red
		if isUnlocked {
			statusIcon = "🔓"
			statusColor = color.Green
		}

		statusColor("%s %s: %d/%d mastered (%.1f%%)",
			statusIcon, difficulty, masteredCount, len(questions), masteryPercentage)

		if isUnlocked {
			color.Green("   ✅ Unlocked - You can progress to the next tier!")
		} else {
			needed := int(float64(len(questions))*0.51) - masteredCount
			if needed <= 0 {
				needed = 1
			}
			color.Yellow("   ⏳ Need %d more mastered questions to unlock", needed)
		}
		fmt.Println()
	}

	return nil
}
