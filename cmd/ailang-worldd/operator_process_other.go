//go:build !darwin && !linux

package main

import "errors"

func listWorldProcesses() ([]worldProcess, error) {
	return nil, errors.New("operator upgrade requires the macOS process inventory")
}

func operatorProcessIdentity(int) (string, bool) { return "", false }
