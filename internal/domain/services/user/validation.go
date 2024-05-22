package user

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (u UserService) validatePhoneMessage(phone string) (valid bool) {
	pattern := `^\+?[1-9]\d{1,14}$`

	match, err := regexp.MatchString(pattern, phone)
	if err != nil {
		return false
	}
	return match
}
func (u UserService) validateNameMessage(name string) (valid bool) {
	nameItems := strings.Split(name, " ")

	validNameItemsLength := len(nameItems) == 2
	if validNameItemsLength {
		return true
	}
	return false
}

func (u UserService) validateBirthDateMessage(birthDate string) (valid bool) {
	currentTime := time.Now()
	dateItems := strings.Split(birthDate, ".")
	validDateToday := false

	validLength := len(birthDate) == 10

	validDateItemsLength := len(dateItems) == 3

	if validDateItemsLength {
		intDay, err := strconv.Atoi(dateItems[0])
		if err != nil {
			return false
		}
		intMonth, err := strconv.Atoi(dateItems[1])
		if err != nil {
			return false
		}
		intYear, err := strconv.Atoi(dateItems[2])
		if err != nil {
			return false
		}

		unvalidatedDate := time.Date(intYear, time.Month(intMonth), intDay, 0, 0, 0, 0, time.UTC)
		validDateToday = unvalidatedDate.Before(currentTime)
	}
	if validLength && validDateItemsLength && validDateToday {
		return true
	}

	return false
}
