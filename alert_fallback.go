//go:build !darwin && !linux

package main

import "fmt"

func alert(message string) {
	fmt.Println("\n*** Notification:", message, "***")
}