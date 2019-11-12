/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package strings

import (
	"fmt"
	"strings"
)

var stringTables [][]string
var reverseTables []map[string]int

var versions []string = []string{
	v0,
	v1,
}

var DefaultVersion = 1

func init() {
	for _, v := range versions {
		s := strings.Split(strings.TrimSpace(v), "\n")
		stringTables = append(stringTables, s)

		m := map[string]int{}
		for i := range s {
			m[s[i]] = i
		}
		reverseTables = append(reverseTables, m)
	}
}

func GetTable(i int) ([]string, error) {
	if i < 0 || i > len(stringTables) {
		return nil, fmt.Errorf("unable to lookup string table version %v", i)
	}
	return stringTables[i], nil
}

func GetReverseTable(i int) (map[string]int, error) {
	if i < 0 || i > len(reverseTables) {
		return nil, fmt.Errorf("unable to lookup reverse string table version %v", i)
	}
	return reverseTables[i], nil
}
