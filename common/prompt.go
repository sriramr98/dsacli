package common

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func PromptInt(question string) (int, error) {
	fmt.Print(question + " ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	input = strings.TrimSpace(input)
	return strconv.Atoi(input)
}
