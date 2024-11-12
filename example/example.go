package main

import (
	"fmt"
	"github.com/modfin/clix"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

type Cfg struct {
	Str      string `cli:"a-str"`
	SubInt   SubInt
	SubSlice SubSlice
}
type SubInt struct {
	Int int `cli:"a-int"`
}
type SubSlice struct {
	SubStr   string   `cli:"a-str"`
	StrSlice []string `cli:"slice"`
}

func main() {
	app := &cli.App{
		Name: "test",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "a-str",
				EnvVars: []string{"STR"},
			},
			&cli.StringFlag{
				Name: "a-int",
			},
			&cli.StringSliceFlag{
				Name: "slice",
			},
		},
		Action: func(context *cli.Context) error {

			// Running
			// export STR="Something strange"
			// go run example/asd.go --a-int 123 --slice="in the" --slice neighborhood

			// Takes the context and parses it recursively into a struct.
			// Since github.com/urfave/cli/v2 supports environment variables.
			// We can use both command line args and/or environment variables
			// to parse the input into a struct.
			cfg := clix.Parse[Cfg](context)
			fmt.Printf("%+v", cfg)
			// { Str:Something strange
			//   SubInt: {Int: 123}
			//   SubSlice: {SubStr:Something strange StrSlice:[in the neighborhood]}
			// }
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
