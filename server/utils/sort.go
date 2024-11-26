package utils

import (
	"oneStepGps/models"
	"sort"
)

func SortDevicesByName(devices []models.Device, ascending bool) {
	sort.Slice(devices, func(i, j int) bool {
		if ascending {
			return devices[i].Name < devices[j].Name
		}
		return devices[i].Name > devices[j].Name
	})
}