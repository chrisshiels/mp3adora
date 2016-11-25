// 'maintagalbum.go'.
// Chris Shiels.


package main


import (
    "flag"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path"
    "regexp"
    "strconv"
)


func tagalbum(stdin *os.File,
              stdout *os.File,
              stderr *os.File,
              verbose bool,
              directorypath string,
              encoding string,
              dryrun bool) (err error) {
    _, directoryname := path.Split(directorypath)

    regexpdirectory :=
        regexp.MustCompile(`^(.*) - ([0-9][0-9][0-9][0-9]) - (.*)$`)
    resultdirectory := regexpdirectory.FindStringSubmatch(directoryname)
    if resultdirectory == nil {
        return fmt.Errorf("Unable to parse directory name %s", directoryname)
    }

    artist := resultdirectory[1]
    year := resultdirectory[2]
    album := resultdirectory[3]

    fileinfos, err := ioutil.ReadDir(directorypath)
    if (err != nil) {
        return err
    }

    for _, fileinfo := range fileinfos {
        if path.Ext(fileinfo.Name()) != ".mp3" {
            fmt.Fprintf(stdout, "Skipping %s\n", fileinfo.Name())
            continue
        }
        fmt.Fprintf(stdout, "Processing %s\n", fileinfo.Name())

        regexpfile :=
            regexp.MustCompile(`^([0-9][0-9]) - ([^-]+) - (.*).mp3$`)
        resultfile := regexpfile.FindStringSubmatch(fileinfo.Name())
        if resultfile == nil {
            return fmt.Errorf("Unable to parse file name %s", fileinfo.Name())
        }

        track, _ := strconv.Atoi(resultfile[1])
        title := resultfile[3]

        id3v1 := newid3v1fromitems(title,
                                   artist,
                                   album,
                                   year,
                                   "",
                                   byte(track),
                                   255)

        file, err := os.Open(path.Join(directorypath, fileinfo.Name()))
        if err != nil {
            return err
        }
        defer file.Close()

        filenew, err := os.Create(path.Join(directorypath,
                                            fmt.Sprintf("%s.new",
                                                        fileinfo.Name())))
        if err != nil {
            return err
        }
        defer filenew.Close()

        mp3adoramp3framecopyhandler := newmp3adoramp3framecopyhandler(filenew)
        mp3adora := newmp3adora(mp3adoramp3framecopyhandler)

        if _, err = mp3adora.parse(file); err != nil {
            if err == io.ErrUnexpectedEOF {
                fmt.Fprintf(stderr,
                            "Warning:  Encounted truncated mp3frame.\n")
                return err
            } else {
                return err
            }
        }

        if _, err = filenew.Write(id3v1.bytes()); err != nil {
            return err
        }

        if err = os.Rename(filenew.Name(), file.Name()); err != nil {
            return err
        }
    }

    return nil
}


func maintagalbum(stdin *os.File,
                  stdout *os.File,
                  stderr *os.File,
                  verbose bool,
                  args []string) (exitstatus int) {
    flagset := flag.NewFlagSet("tagalbum", flag.ExitOnError)

    flagset.Usage = func() {
        fmt.Fprintln(stdout,
                     "Usage:  mp3adora [ -v ] tagalbum [ options ] directory ...")
        fmt.Fprintln(stdout)
        fmt.Fprintln(stdout, "Options:")
        flagset.PrintDefaults()
    }

    flagencoding := flagset.String("encoding",
                                   "utf-8",
                                   "Encoding")
    flagn := flagset.Bool("n",
                          false,
                          "Dry-run")

    // Note flagset.Parse() will also handle '-h' and '--help' and will exit
    // with exit status 2.
    flagset.Parse(args)

    if len(flagset.Args()) == 0 {
        flagset.Usage()
        return exitfailure
    }

    for i, directoryname := range flagset.Args() {
        if i > 0 {
            fmt.Fprintln(stdout)
        }
        fmt.Fprintf(stdout, "%s:\n", directoryname)

        if err := tagalbum(stdin,
                           stdout,
                           stderr,
                           verbose,
                           directoryname,
                           *flagencoding,
                           *flagn); err != nil {
            fmt.Fprintf(stderr, "mp3adora: %s\n", err)
            return exitfailure
        }
    }

    return exitsuccess
}
