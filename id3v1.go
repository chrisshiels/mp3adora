// 'id3v1.go'.
// Chris Shiels.


package main


import (
    "fmt"
)


// See:  https://en.wikipedia.org/wiki/ID3#ID3v1
// "- The ID3v1 tag occupies 128 bytes, beginning with the string TAG
//    128 bytes from the end of the file.
//    Strings are either space- or zero-padded.
//    Unset string entries are filled using an empty string."
type id3v1 struct {
    header string
    title string
    artist string
    album string
    year string
    comment string
    track byte
    genre byte
}


func newid3v1frombytes(bytes []byte) (i *id3v1, err error) {
    i = new(id3v1)
    i.header = string(bytes[0:3])

    if i.header != "TAG" {
        return nil, fmt.Errorf("Unable to find id3v1 header.")
    }

    i.title = string(bytes[3:33])
    i.artist = string(bytes[33:63])
    i.album = string(bytes[63:93])
    i.year = string(bytes[93:97])
    i.comment = string(bytes[97:125])
    i.track = bytes[126]
    i.genre = bytes[127]

    return i, nil
}


func newid3v1fromitems(header string,
                       title string,
                       artist string,
                       album string,
                       year string,
                       comment string,
                       track byte,
                       genre byte) (i *id3v1) {
    return &id3v1{ header: "TAG",
                   title: title,
                   artist: artist,
                   album: album,
                   year: year,
                   comment: comment,
                   track: track,
                   genre: genre }
}


func (i *id3v1)bytes() []byte {
    bytes := make([]byte, 128)
    copy(bytes[0:3], "TAG")
    copy(bytes[3:33], i.title)
    copy(bytes[33:63], i.artist)
    copy(bytes[63:93], i.album)
    copy(bytes[93:97], i.year)
    copy(bytes[97:125], i.comment)
    bytes[125] = 0
    bytes[126] = i.track
    bytes[127] = i.genre
    return bytes
}
