package lib

import (
	"os"
	"strconv"
)

// Napis funkci, ktera ma na vstupu jmeno promenne a defaultni hodnotu a vraci hodnotu promenne v spravnem typu
// pokud promenna neexistuje vraci defaultni hodnotu
// pokud promenna existuje vraci ji ve spravnem typu
func GetEnv[T interface{}](name string, def T) T {
	ret, err := os.LookupEnv(name)
	if err != nil {
		// pokud nemam promennou vracim default
		return def

	} else {
		// pokud mam promennou vracim ji ve spravnem typu
		if def.(type) == int {
			retInt, err := strconv.Atoi(s)
			if err != nil {
				return ret
			} else {
				return retInt
			}
		} else {
			return ret
		}
	}
}
