/*
 * Copyright (C) 2020-2021 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package collection

import "strings"

// Removes spaces leaving only the strings
func parseListWithCleanup(input string, sep string) (newS []string) {
	if len(input) == 0 {
		newS = []string{} // initialisation of empty arrays in function returns []string(nil) instead of []string{}
		return
	}
	split := strings.Split(input, sep)
	for _, s := range split {
		tempString := strings.TrimSpace(s)
		if tempString != "" {
			newS = append(newS, tempString)
		}
	}
	return
}

func ParseCommaSeparatedList(input string) []string {
	return parseListWithCleanup(input, ",")
}
