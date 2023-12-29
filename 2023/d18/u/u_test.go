package u

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseNumsWithSpaces(t *testing.T) {
	// Test with spaces
	input := "6    4  1 8"
	sp := strings.Split(input, " ")

	out := ParseNums(sp)

	fmt.Println(out)
}