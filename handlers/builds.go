package handlers

import (
	"bytes"
	"fmt"
	"github.com/autoabs/autoabs/build"
	"github.com/autoabs/autoabs/database"
	"github.com/gin-gonic/gin"
	"text/tabwriter"
)

func buildsGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	buf := &bytes.Buffer{}
	tw := tabwriter.NewWriter(buf, 0, 0, 3, ' ', 0)

	builds, err := build.GetAll(db)
	if err != nil {
		return
	}

	fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		"BUILD ID",
		"NAME",
		"STATE",
		"VERSION",
		"RELEASE",
		"REPO",
		"ARCH",
		"START",
		"STOP",
	)

	build.Sort(builds)

	for _, bild := range builds {
		start := "-"
		stop := "-"

		if !bild.Start.IsZero() {
			start = bild.Start.String()
		}

		if !bild.Stop.IsZero() {
			stop = bild.Stop.String()
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			bild.Id.Hex(),
			bild.Name,
			bild.State,
			bild.Version,
			bild.Release,
			bild.Repo,
			bild.Arch,
			start,
			stop,
		)
	}

	tw.Flush()

	c.String(200, string(buf.Bytes()))
}
