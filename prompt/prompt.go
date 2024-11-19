package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func StringInputPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")

		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func IntegerInputPrompt(label string) (int, error) {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")

		input, _ := r.ReadString('\n')
		input = strings.TrimSpace(input)

		// attempt to convert the input string to an integer
		integerValue, err := strconv.Atoi(input)
		if err != nil {
			return 0, errors.New("Please enter a valid number")
		}

		return integerValue, nil
	}
}
