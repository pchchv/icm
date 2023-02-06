package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	ui "github.com/gizak/termui/v3"
	termbox "github.com/nsf/termbox-go"
	"github.com/pchchv/icm/logger"
)

func panicExit() {
	if r := recover(); r != nil {
		Shutdown()
		panic(r)

		fmt.Printf("error: %s\n", r)

		os.Exit(1)
	}
}

func Shutdown() {
	log.Notice("shutting down")
	log.Exit()

	if termbox.IsInit {
		ui.Close()
	}
}

var (
	build     = "none"
	version   = "dev-build"
	goVersion = runtime.Version()

	log *logger.CTopLogger

	versionStr = fmt.Sprintf("icm version %v, build %v %v", version, build, goVersion)
)

func main() {
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

	defer panicExit()

	flag.Parse()
}
