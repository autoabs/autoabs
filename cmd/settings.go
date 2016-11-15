package cmd

import (
	"flag"
	"github.com/pritunl/pritunl-auth/database"
	"github.com/pritunl/pritunl-auth/settings"
	"strconv"
)

// Manages system settigs stored in database
func Settings() {
	group := flag.Arg(1)
	key := flag.Arg(2)
	val := flag.Arg(3)
	db := database.GetDatabase()
	defer db.Close()

	var valParsed interface{}
	if x, err := strconv.Atoi(val); err == nil {
		valParsed = x
	} else {
		valParsed = val
	}

	err := settings.Set(db, group, key, valParsed)
	if err != nil {
		panic(err)
	}

	return
}
