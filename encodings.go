// 'encodings.go'.
// Chris Shiels.


package main


import (
    "fmt"
    "strings"

    "golang.org/x/text/encoding"
    "golang.org/x/text/encoding/charmap"
)


var charmaps = map[string]encoding.Encoding {
    "iso8859-1":  charmap.ISO8859_1,
    "iso8859-2":  charmap.ISO8859_2,
    "iso8859-3":  charmap.ISO8859_3,
    "iso8859-4":  charmap.ISO8859_4,
    "iso8859-5":  charmap.ISO8859_5,
    "iso8859-6":  charmap.ISO8859_6,
    "iso8859-7":  charmap.ISO8859_7,
    "iso8859-8":  charmap.ISO8859_8,
    "iso8859-10": charmap.ISO8859_10,
    "iso8859-13": charmap.ISO8859_13,
    "iso8859-14": charmap.ISO8859_14,
    "iso8859-15": charmap.ISO8859_15,
    "iso8859-16": charmap.ISO8859_16,
}


func find(name string) (e encoding.Encoding, err error) {
    e, ok := charmaps[name]
    if !ok {
        return nil, fmt.Errorf("Unrecognised encoding %s", name)
    }
    return e, nil
}


func convert(e encoding.Encoding,
             s string,
             replacement byte) (s1 string, err error) {
    s1, err = encoding.ReplaceUnsupported(e.NewEncoder()).String(s)
    if err != nil {
        return "", err
    }

    // Fix encoding.ReplaceUnsupported()'s encoding specific replacement.
    return strings.Replace(s1, "\x1a", string(replacement), -1), nil
}
