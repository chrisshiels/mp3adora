// 'main.go'.
// Chris Shiels.


package main


import (
    "flag"
    "fmt"
    "os"
)


const exitsuccess = 0
const exitfailure = 1


func _main(stdin *os.File,
           stdout *os.File,
           stderr *os.File,
           args []string) (exitstatus int) {
    flagset := flag.NewFlagSet(args[0], flag.ExitOnError)

    flagset.Usage = func() {
        fmt.Fprintln(stdout, "Usage:  mp3adora [ -v ] command options ...")
        fmt.Fprintln(stdout)
        fmt.Fprintln(stdout, "Commands:")
        fmt.Fprintln(stdout, "show        Parse contents of mp3 files")
        fmt.Fprintln(stdout)
        fmt.Fprintln(stdout, "Options:")
        flagset.PrintDefaults()
    }

    flagv := flagset.Bool("v",
                          false,
                          "Verbose")

    // Note flagset.Parse() will also handle '-h' and '--help' and will exit
    // with exit status 2.
    flagset.Parse(args[1:])

    if len(flagset.Args()) == 0 {
        flagset.Usage()
        return exitfailure
    }

    switch {
        case flagset.Args()[0] == "show":
            return mainshow(stdin,
                            stdout,
                            stderr,
                            *flagv,
                            flagset.Args()[1:])
    }

    flagset.Usage()
    return exitfailure
}


func main() {
    os.Exit(_main(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
