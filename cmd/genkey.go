package cmd

import (
	"github.com/autoabs/autoabs/signing"
	"os"
)

func GenKey() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	genkey := &signing.GenKey{
		Root:  wd,
		Name:  "AutoABS",
		Email: "build@autoabs.com",
	}

	err = genkey.Generate()
	if err != nil {
		panic(err)
	}

	err = genkey.Export()
	if err != nil {
		panic(err)
	}
}
