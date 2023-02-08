package utils

import (
	"strings"

	"github.com/shirou/gopsutil/disk"
)

/*
	Drive Letters funcs are getting mounted/available mounting points (letters)
	The implementation is taken from here: https://stackoverflow.com/a/57569918/10352443

	Innokentiy Sokolov
	https://github.com/keshon

	2022-03-24
*/

// Array of mounted points (letters)
func MountedDriveLetters() []string {
	var letters []string
	partitions, _ := disk.Partitions(false)
	for _, partition := range partitions {
		letters = append(letters, strings.ReplaceAll(partition.Mountpoint, ":", ""))
	}

	return letters
}

// Array of available mounting points (letters)
func AvalableDriveLetters() []string {
	allLetters := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	return StringArrayDifference(allLetters, MountedDriveLetters())
}
