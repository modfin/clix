module github.com/modfin/clix/example/cli3

go 1.24.0

require (
	github.com/modfin/clix v0.0.0
	github.com/urfave/cli/v3 v3.0.0-beta1
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli/v2 v2.27.6 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
)
replace (
	github.com/modfin/clix => "../.."
)