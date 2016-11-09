// 'mp3adorashowhandler.go'.
// Chris Shiels.


package main


import (
    "fmt"
    "io"
)


type mp3adorashowhandler struct {
    stdout io.Writer
    stderr io.Writer
}


func newmp3adorashowhandler(stdout io.Writer,
                            stderr io.Writer) *mp3adorashowhandler {
    return &mp3adorashowhandler{ stdout: stdout,
                                 stderr: stderr }
}


func (h *mp3adorashowhandler) processape(bytes []byte) (err error) {
    fmt.Fprintf(h.stdout, "ape:       %d bytes:  %v\n", len(bytes), bytes)
    return nil
}


func (h *mp3adorashowhandler) processid3v1(bytes []byte) (err error) {
    var i *id3v1
    if i, err = newid3v1frombytes(bytes); err != nil {
        return err
    }

    fmt.Fprintf(h.stdout, "id3v1:     %d bytes:  ", len(bytes))
    fmt.Fprintf(h.stdout, "header: %s, ", i.header)
    fmt.Fprintf(h.stdout, "title: %s, ", i.title)
    fmt.Fprintf(h.stdout, "artist: %s, ", i.artist)
    fmt.Fprintf(h.stdout, "album: %s, ", i.album)
    fmt.Fprintf(h.stdout, "year: %s, ", i.year)
    fmt.Fprintf(h.stdout, "comment: %s, ", i.comment)
    fmt.Fprintf(h.stdout, "track: %d, ", i.track)
    fmt.Fprintf(h.stdout, "genre: %d\n", i.genre)

    return nil
}


func (h *mp3adorashowhandler) processid3v2(bytes []byte) (err error) {
    fmt.Fprintf(h.stdout, "id3v2:     %d bytes:  %v\n", len(bytes), bytes)
    return nil
}


func (h *mp3adorashowhandler) processmp3frame(bytes []byte) (err error) {
    var m *mp3header
    if m, err = newmp3headerfrombytes(bytes); err != nil {
        return err
    }

    fmt.Fprintf(h.stdout, "mp3frame:  %d bytes:  ", m.size)
    fmt.Fprintf(h.stdout, "audioversion: %1.2f, ", m.audioversion)
    fmt.Fprintf(h.stdout, "layer: %d, ", m.layer)
    fmt.Fprintf(h.stdout, "protection: %t, ", m.protection)
    fmt.Fprintf(h.stdout, "bitrate: %d, ", m.bitrate)
    fmt.Fprintf(h.stdout, "samplingrate: %d, ", m.samplingrate)
    fmt.Fprintf(h.stdout, "padding: %t, ", m.padding)
    fmt.Fprintf(h.stdout, "private: %t, ", m.private)
    fmt.Fprintf(h.stdout, "channelmode: %d, ", m.channelmode)
    fmt.Fprintf(h.stdout, "modeextension: %d, ", m.modeextension)
    fmt.Fprintf(h.stdout, "copyright: %t, ", m.copyright)
    fmt.Fprintf(h.stdout, "original: %t, ", m.original)
    fmt.Fprintf(h.stdout, "emphasis: %d\n", m.emphasis)

    return nil
}


func (h *mp3adorashowhandler) processunrecognised(byte byte) (err error) {
    fmt.Fprintf(h.stdout, "unrecognised:  1 byte:  %v\n", byte)
    return nil
}
