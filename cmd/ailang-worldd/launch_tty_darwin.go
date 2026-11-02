//go:build darwin

package main

import "syscall"

const ioctlGetTermios = syscall.TIOCGETA
