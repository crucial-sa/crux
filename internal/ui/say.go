package ui

import (
	"fmt"
	"strings"
)

func Say(message ...string) {
	fmt.Printf("%s\n", strings.Join(message, ""))
}
