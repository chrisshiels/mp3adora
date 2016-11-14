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


func show(stdin *os.File,
          stdout *os.File,
          stderr *os.File,
          filename string,
          verbose bool) (size int, err error) {
    mp3adorashowhandler := newmp3adorashowhandler(stdout, stderr)
    mp3adora := newmp3adora(mp3adorashowhandler)

    var file *os.File
    if filename != "" {
        if file, err = os.Open(filename); err != nil {
            return 0, err
        }
    } else {
        file = stdin
    }

    defer file.Close()

    if size, err = mp3adora.parse(file); err != nil {
        return 0, err
    }

    return size, nil
}


func mainshow(stdin *os.File,
              stdout *os.File,
              stderr *os.File,
              args []string,
              verbose bool) (exitstatus int) {
    flagset := flag.NewFlagSet("show", flag.ExitOnError)

    flagset.Usage = func() {
        fmt.Fprintln(stdout, "Usage:  mp3adora [ -v ] show [ filename ... ]")
        flag.PrintDefaults()
    }

    // Note flagset.Parse() will also handle '-h' and '--help' and will exit
    // with exit status 2.
    flagset.Parse(args)

    if len(flagset.Args()) == 0 {
        size, err := show(stdin, stdout, stderr, "", verbose)
        if err != nil {
            fmt.Fprintf(stderr, "mp3adora: %s\n", err)
            return exitfailure
        }
        fmt.Fprintf(stdout, "size:  %d\n", size)
    } else {
        for i, filename := range flagset.Args() {
            if i > 0 {
                fmt.Fprintln(stdout)
            }
            fmt.Fprintf(stdout, "%s:\n", filename)

            size, err := show(stdin, stdout, stderr, filename, verbose)
            if err != nil {
                fmt.Fprintf(stderr, "mp3adora: %s\n", err)
                return exitfailure
            }
            fmt.Fprintf(stdout, "size:  %d\n", size)
        }
    }

    return exitsuccess
}


func _main(stdin *os.File,
           stdout *os.File,
           stderr *os.File,
           args []string) (exitstatus int) {
    flagset := flag.NewFlagSet(args[0], flag.ExitOnError)

    flagset.Usage = func() {
        fmt.Fprintln(stdout, "Usage:  mp3adora [ -v ] command [ options ... ]")
        fmt.Fprintln(stdout)
        fmt.Fprintln(stdout, "Valid commands are:")
        fmt.Fprintln(stdout, "show")
        flag.PrintDefaults()
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
            return mainshow(stdin, stdout, stderr, flagset.Args()[1:], *flagv)
    }

    flagset.Usage()
    return exitfailure
}


func main() {
    os.Exit(_main(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
