package common

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"strconv"
	"strings"
)

func OneToFiveRatingValidator(input string) error {
	value, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || value < 1 || value > 5 {
		return fmt.Errorf("please enter a number between 1 and 5")
	}
	return nil
}

func NumberValidator(input string) error {
	num, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return fmt.Errorf("please enter a valid number >= -1")
	}
	if num < -1 {
		return fmt.Errorf("please enter a number >= -1")
	}
	return nil
}

type IntValidator func(string) error

func PromptInt(question string, validator IntValidator) (int, error) {
	prompt := promptui.Prompt{
		Label: question,
		Validate: func(s string) error {
			return validator(s)
		},
	}
	result, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			return 0, fmt.Errorf("interrupted by user")
		}
		return 0, fmt.Errorf("prompt failed: %w", err)
	}
	value, err := strconv.Atoi(strings.TrimSpace(result))
	return value, err
}

func PromptSelect(question string, items []string) (int, error) {
	prompt := promptui.Select{
		Label: question,
		Items: items,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			return -1, fmt.Errorf("interrupted by user")
		}
		return -1, fmt.Errorf("prompt failed: %w", err)
	}
	return idx, nil
}
