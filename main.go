// 'main.go'.
// Chris Shiels.


package main


import (
    "fmt"
    "os"
)


const exitsuccess = 0
const exitfailure = 1


func _main(stdin *os.File,
           stdout *os.File,
           stderr *os.File,
           args []string) (exitstatus int) {
    mp3adorashowhandler := newmp3adorashowhandler(stdout, stderr)
    mp3adora := newmp3adora(mp3adorashowhandler)

    var file *os.File
    var err error
    if file, err = os.Open("file.mp3"); err != nil {
        fmt.Fprintf(stderr, "mp3adora: %s\n", err)
        return exitfailure
    }

    defer file.Close()

    var size int
    if size, err = mp3adora.parse(file); err != nil {
        fmt.Fprintf(stderr, "mp3adora: %s\n", err)
        return exitfailure
    }

    fmt.Printf("size:  %d\n", size)

    return exitsuccess
}


func main() {
    os.Exit(_main(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
