// 'encodings_test.go'.
// Chris Shiels.


package main


import (
    "testing"

    "golang.org/x/text/encoding/charmap"
)


func Test_knownencoding(t *testing.T) {
    e, err := find("iso8859-1")
    if ! (e == charmap.ISO8859_1 && err == nil) {
        t.Errorf("Test_knownencoding:  failed")
        return
    }
}


func Test_unknownencoding(t *testing.T) {
    e, err := find("iso8859-9")
    if ! (e == nil && err != nil) {
        t.Errorf("Test_unknownencoding:  failed")
        return
    }
}


func Test_convertsuccessful(t *testing.T) {
    e, err := find("iso8859-1")
    if ! (e == charmap.ISO8859_1 && err == nil) {
        t.Errorf("Test_convertsuccessful:  failed")
        return
    }

    s1, err := convert(e, "buenos días", '?')
    if ! (s1 == "buenos d\xedas" && err == nil) {
        t.Errorf("Test_convertsuccessful:  failed")
        return
    }
}


func Test_convertnotneeded(t *testing.T) {
    e, err := find("iso8859-1")
    if ! (e == charmap.ISO8859_1 && err == nil) {
        t.Errorf("Test_convertnotneeded:  failed")
        return
    }

    s1, err := convert(e, "buenos dias", '?')
    if ! (s1 == "buenos dias" && err == nil) {
        t.Errorf("Test_convertnotneeded:  failed")
        return
    }
}


func Test_convertunsupportedcharacter(t *testing.T) {
    e, err := find("iso8859-1")
    if ! (e == charmap.ISO8859_1 && err == nil) {
        t.Errorf("Test_convertunsupportedcharacter:  failed")
        return
    }

    s1, err := convert(e, "€1", '?')
    if ! (s1 == "?1" && err == nil) {
        t.Errorf("Test_convertunsupportedcharacter:  failed")
        return
    }
}
