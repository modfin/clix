# clix 
> small library meant to turn a command line args into config structs

clix is a helper function for `github.com/urfave/cli`


## Installation 
```bash
go get github.com/modfin/clix@latest
```


## Usage V2

```go 
package main

import (
	"clix"
	"fmt"
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
			// go run example/example.go --a-int 123 --slice="in the" --slice neighborhood

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

```


## Usage V3

```go 
package main

import (
	"context"
	"fmt"
	"github.com/modfin/clix"
	cli "github.com/urfave/cli/v3"
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
	cmd := &cli.Command{
		Name: "test",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "a-str",
				Sources: cli.EnvVars("STR"),
			},
			&cli.StringFlag{
				Name: "a-int",
			},
			&cli.StringSliceFlag{
				Name: "slice",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {

			// Running
			// export STR="Something strange"
			// go run example.go --a-int 123 --slice="in the" --slice neighborhood

			// Takes the context and parses it recursively into a struct.
			// Since github.com/urfave/cli/v2 supports environment variables.
			// We can use both command line args and/or environment variables
			// to parse the input into a struct.

			cfg := clix.Parse[Cfg](clix.V3(cmd))
			fmt.Printf("%+v", cfg)

			// { Str:Something strange
			//   SubInt: {Int: 123}
			//   SubSlice: {SubStr:Something strange StrSlice:[in the neighborhood]}
			// }

			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
```