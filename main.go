package main

import (
	"flag"
)

func nmain() {
	var (
		versionFlag     = flag.Bool("v", false, "output version information and exit")
		helpFlag        = flag.Bool("h", false, "display this help dialog")
		filterFlag      = flag.String("f", "", "filter containers")
		activeOnlyFlag  = flag.Bool("a", false, "show active containers only")
		sortFieldFlag   = flag.String("s", "", "select container sort field")
		reverseSortFlag = flag.Bool("r", false, "reverse container sort order")
		invertFlag      = flag.Bool("i", false, "invert default colors")
		connectorFlag   = flag.String("connector", "docker", "container connector to use")
	)
	flag.Parse()
}
