package main

func displayCommitter(value string) string {
	if value == "" {
		return "(repository git config)"
	}
	return value
}
