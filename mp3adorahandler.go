// 'mp3adorahandler.go'.
// Chris Shiels.


package main


import (
)


type mp3adorahandler interface {
    processape(bytes []byte) (err error)
    processid3v1(bytes []byte) (err error)
    processid3v2(bytes []byte) (err error)
    processmp3frame(bytes []byte) (err error)
    processunrecognised(byte byte) (err error)
}
