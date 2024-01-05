package main

import (
	"fmt"
	"fuaoc2023/day20/u"
	"strings"
)

const Part1Debug = 1

func main() {
	fmt.Println(Part1(u.Linewisefile_chan("input")))
	fmt.Println(Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {

	var connections = make(map[string][]string)
	var incoming_connections = make(map[string][]string)
	var flipflop_temp = make(map[string]bool)
	var conjunction_temp = make(map[string][]string)

	var name_to_type = make(map[string]byte)
	// var linenum ModuleID
	for line := range lines {

		sp := strings.Split(line, " -> ")
		module_type_string := sp[0][0]
		module_name := sp[0][1:]
		switch module_type_string {
		case 'b':
			module_name = "broadcaster"
		case '%':
			// flipflop
			flipflop_temp[module_name] = true
		case '&':
			conjunction_temp[module_name] = []string{}
		}
		name_to_type[module_name] = module_type_string

		// Connections
		for _, destname := range strings.Split(sp[1], ", ") {
			connections[module_name] = append(connections[module_name], destname)
			incoming_connections[destname] = append(incoming_connections[destname], module_name)
		}
		// linenum++
	}
	// Add button
	name_to_type["button"] = 'x'
	connections["button"] = append(connections["button"], "broadcaster")

	// Conjunction prepopulate
	for name := range conjunction_temp {
		conjunction_temp[name] = append(conjunction_temp[name], incoming_connections[name]...)
	}

	// Final construction
	var ffstate uint64
	var ffstate_next uint64 = 1
	var conjstate uint64
	var conjstate_next uint64 = 1

	var flipflops = make(map[string]*Flipflop)
	var conjunctions = make(map[string]*Conjunction)

	for name := range flipflop_temp {
		flipflops[name] = &Flipflop{
			statemask: ffstate_next,
		}
		ffstate_next <<= 1
	}

	for name, incoming_list := range conjunction_temp {
		var new_mask uint64
		new_connections_map := make(map[string]uint64)
		for _, name := range incoming_list {
			new_connections_map[name] = conjstate_next
			new_mask |= conjstate_next
			conjstate_next <<= 1
		}

		conjunctions[name] = &Conjunction{
			in_connections: new_connections_map,
			statemask:      new_mask,
		}
	}

	// Simulation
	// var history [][2]uint64
	// var state_to_history_index = make(map[[2]uint64]int)
	// var pulse_history [][2]int
	var high_pulse_count int
	var low_pulse_count int

	var pulse_queue PulseQueue
	var button_presses int
	var press_high_pulse_count int
	var press_low_pulse_count int
	// push_history := func() int {
	// 	state := [2]uint64{ffstate, conjstate}
	// 	history = append(history, state)
	// 	pulse_history = append(pulse_history, [2]int{press_high_pulse_count, press_low_pulse_count})

	// 	if history_index, ok := state_to_history_index[state]; ok {
	// 		return history_index
	// 	} else {
	// 		state_to_history_index[state] = button_presses
	// 	}
	// 	return -1
	// }
	// push_history()
	for button_presses = 1; button_presses <= 1000; button_presses++ {
		press_high_pulse_count = 0
		press_low_pulse_count = 0
		// Button Press
		pulse_queue.Push(Pulse{
			source:      "button",
			destination: "broadcaster",
			high:        false,
		})

		// Process pulses until empty
		for mpulse, not_empty := pulse_queue.Pop(); not_empty; mpulse, not_empty = pulse_queue.Pop() {
			if mpulse.high {
				press_high_pulse_count++
			} else {
				press_low_pulse_count++
			}

			dtype := name_to_type[mpulse.destination]
			new_pulse_template := Pulse{
				source:      mpulse.destination,
				destination: "",
				high:        false,
			}
			switch dtype {
			case 'b':
				new_pulse_template.high = mpulse.high
			case '%':
				// Flipflop
				if mpulse.high {
					// Ignored
					continue
				}
				ff := flipflops[mpulse.destination]
				ffstate ^= ff.statemask
				new_pulse_template.high = (ffstate & ff.statemask) > 0
			case '&':
				// Conjunction
				conj := conjunctions[mpulse.destination]
				input_mask, ok := conj.in_connections[mpulse.source]
				if !ok {
					panic("input not in conj's input map")
				}
				if mpulse.high {
					conjstate |= input_mask
				} else {
					conjstate &= ^input_mask
				}

				if (conjstate & conj.statemask) == conj.statemask {
					new_pulse_template.high = false
				} else {
					new_pulse_template.high = true
				}
			}

			// Push a pulse to each output
			for _, output := range connections[mpulse.destination] {
				new_pulse := new_pulse_template
				new_pulse.destination = output
				pulse_queue.Push(new_pulse)
			}
		}

		// if found_index := push_history(); found_index >= 0 {
		// 	panic("WOW")
		// }
		high_pulse_count += press_high_pulse_count
		low_pulse_count += press_low_pulse_count
	}

	return high_pulse_count * low_pulse_count
}

func Part2(lines <-chan string) int {
	var connections = make(map[string][]string)
	var incoming_connections = make(map[string][]string)
	var flipflop_temp = make(map[string]bool)
	var conjunction_temp = make(map[string][]string)

	var name_to_type = make(map[string]byte)
	// var linenum ModuleID
	for line := range lines {

		sp := strings.Split(line, " -> ")
		module_type_string := sp[0][0]
		module_name := sp[0][1:]
		switch module_type_string {
		case 'b':
			module_name = "broadcaster"
		case '%':
			// flipflop
			flipflop_temp[module_name] = true
		case '&':
			conjunction_temp[module_name] = []string{}
		}
		name_to_type[module_name] = module_type_string

		// Connections
		for _, destname := range strings.Split(sp[1], ", ") {
			connections[module_name] = append(connections[module_name], destname)
			incoming_connections[destname] = append(incoming_connections[destname], module_name)
		}
		// linenum++
	}
	// Add button
	name_to_type["button"] = 'x'
	connections["button"] = append(connections["button"], "broadcaster")

	// Conjunction prepopulate
	for name := range conjunction_temp {
		conjunction_temp[name] = append(conjunction_temp[name], incoming_connections[name]...)
	}

	// Final construction
	var ffstate uint64
	var ffstate_next uint64 = 1
	var conjstate uint64
	var conjstate_next uint64 = 1

	var flipflops = make(map[string]*Flipflop)
	var conjunctions = make(map[string]*Conjunction)

	for name := range flipflop_temp {
		flipflops[name] = &Flipflop{
			statemask: ffstate_next,
		}
		ffstate_next <<= 1
	}

	for name, incoming_list := range conjunction_temp {
		var new_mask uint64
		new_connections_map := make(map[string]uint64)
		for _, name := range incoming_list {
			new_connections_map[name] = conjstate_next
			new_mask |= conjstate_next
			conjstate_next <<= 1
		}

		conjunctions[name] = &Conjunction{
			in_connections: new_connections_map,
			statemask:      new_mask,
		}
	}

	// Simulation
	var pulse_queue PulseQueue
	var button_presses int

	important_senders := make(map[string]int)
	for i, name := range incoming_connections["jq"] {
		important_senders[name] = i
	}
	var low_pulse_cycle = make(map[string]int)

	button_loop:
	for button_presses = 1; ; button_presses++ {
		var rx_low_presses int
		var rx_high_presses int
		// Button Press
		pulse_queue.Push(Pulse{
			source:      "button",
			destination: "broadcaster",
			high:        false,
		})

		// Process pulses until empty
		for mpulse, not_empty := pulse_queue.Pop(); not_empty; mpulse, not_empty = pulse_queue.Pop() {
			if mpulse.destination == "rx" {
				if mpulse.high {
					rx_high_presses++
				} else {
					rx_low_presses++
				}
			}
			if mpulse.destination == "jq" && mpulse.high {
				low_pulse_cycle[mpulse.source] = button_presses
				if len(low_pulse_cycle) == len(important_senders) {
					break button_loop
				}
			}

			dtype := name_to_type[mpulse.destination]
			new_pulse_template := Pulse{
				source:      mpulse.destination,
				destination: "",
				high:        false,
			}
			switch dtype {
			case 'b':
				new_pulse_template.high = mpulse.high
			case '%':
				// Flipflop
				if mpulse.high {
					// Ignored
					continue
				}
				ff := flipflops[mpulse.destination]
				ffstate ^= ff.statemask
				new_pulse_template.high = (ffstate & ff.statemask) > 0
			case '&':
				// Conjunction
				conj := conjunctions[mpulse.destination]
				input_mask, ok := conj.in_connections[mpulse.source]
				if !ok {
					panic("input not in conj's input map")
				}
				if mpulse.high {
					conjstate |= input_mask
				} else {
					conjstate &= ^input_mask
				}

				if (conjstate & conj.statemask) == conj.statemask {
					new_pulse_template.high = false
				} else {
					new_pulse_template.high = true
				}
			}

			// Push a pulse to each output
			for _, output := range connections[mpulse.destination] {
				new_pulse := new_pulse_template
				new_pulse.destination = output
				pulse_queue.Push(new_pulse)
			}
		}
	}

	var cycles []int
	for _, val := range low_pulse_cycle {
		cycles = append(cycles, val)
	}
	return u.LCM(cycles...)
}
