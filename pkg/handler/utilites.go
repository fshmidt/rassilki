package handler

import (
	"github.com/fshmidt/rassilki"
)

func findUpdateLevel(original rassilki.Rassilka, input rassilki.UpdateRassilka) int {
	if (input.Message != nil && *input.Message != original.Message) || (input.Filter != nil && original.Filter == nil) {
		*input.Recreated = true
		return recreate_messages_table
	}
	if input.Filter != nil && original.Filter != nil {

		if len(*input.Filter) < len(original.Filter) {
			return deleting_tags
		}

		originalSlice := []string(original.Filter)

		NotContain, _ := notContain(*input.Filter, originalSlice)
		if NotContain {
			*input.Recreated = true
			return deleting_tags
		}

		if len(*input.Filter) > len(original.Filter) {
			*input.Supplemented = true
			return add_clients_to_messages_table
		}
	}

	if input.StartTime != nil {
		if *input.StartTime != original.StartTime {
			return time_changes
		}
	}
	if input.EndTime != nil {
		if *input.EndTime != original.EndTime {
			return time_changes
		}
	}
	return no_changes
}

func notContain(input, original []string) (bool, []string) {
	inputSet := make(map[string]struct{})
	var differance []string

	for _, val := range input {
		inputSet[val] = struct{}{}
	}

	for _, val := range original {
		if _, exists := inputSet[val]; !exists {
			return true, nil
		} else {
			delete(inputSet, val)
		}
	}
	for k, _ := range inputSet {
		differance = append(differance, k)
	}
	return false, differance
}
