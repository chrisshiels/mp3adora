// 'mp3adoramp3framecopyhandler.go'.
// Chris Shiels.


package main


import (
    "fmt"
    "io"
)


type mp3adoramp3framecopyhandler struct {
    out io.Writer
}


func newmp3adoramp3framecopyhandler(out io.Writer) *mp3adoramp3framecopyhandler {
    return &mp3adoramp3framecopyhandler{ out: out }
}


func (h *mp3adoramp3framecopyhandler) processape(bytes []byte) (err error) {
    return nil
}


func (h *mp3adoramp3framecopyhandler) processid3v1(bytes []byte) (err error) {
    return nil
}


func (h *mp3adoramp3framecopyhandler) processid3v2(bytes []byte) (err error) {
    return nil
}


func (h *mp3adoramp3framecopyhandler) processmp3frame(bytes []byte) (err error) {
    if _, err := h.out.Write(bytes); err != nil {
        return err
    }
    return nil
}


func (h *mp3adoramp3framecopyhandler) processunrecognised(byte byte) (err error) {
    fmt.Fprintf(h.out, "Unrecognised %d %c\n", byte, byte)
    return nil
}
