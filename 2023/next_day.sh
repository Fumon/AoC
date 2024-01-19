#!/bin/bash

# Get the current directory name
current_dir=$(basename "$PWD")

# Extract the number from the directory name
if [[ $current_dir =~ d([0-9]{2}) ]]; then
    number=${BASH_REMATCH[1]}
else
    echo "Error: Current directory name does not match expected format 'd##'."
    exit 1
fi

# Increment the number
next_number=$((10#$number + 1)) # 10# ensures the number is treated as decimal

# Format the new directory name
next_dir=$(printf "../d%02d" $next_number)

# Create the new directory
mkdir "$next_dir"

# Copy the 'u' folder
cp -r ./u "$next_dir"

# Change to the new directory
cd "$next_dir"

# Initialize the Go module
go mod init "fuaoc2023/day$next_number"

# Create blank files
touch input testinput

# Create main.go with content
cat <<EOF > main.go
package main

import (
	"fmt"
	"fuaoc2023/day$next_number/u"
)

func main() {
	fmt.Println(Part1(u.Linewisefile_chan("input")))
	fmt.Println(Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {
	return 0
}

func Part2(lines <-chan string) int {
	return 0
}
EOF

# Create main_test.go with content
cat <<EOF > main_test.go
package main

import (
	"fuaoc2023/day$next_number/u"
	"testing"
)

func TestPart1(t *testing.T) {
	expected := 0
	result := Part1(u.Linewisefile_chan("testinput"))
	if result != expected {
		t.Error("Part1 failed on testinput. Expected ", expected, " got ", result)
	}
}

func TestPart2(t *testing.T) {
	expected := 0
	result := Part2(u.Linewisefile_chan("testinput"))
	if result != expected {
		t.Error("Part1 failed on testinput. Expected ", expected, " got ", result)
	}
}
EOF

# Open Visual Studio Code with the new directory
code .