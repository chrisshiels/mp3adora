// 'mp3adora.go'.
// Chris Shiels.


package main


import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
)


type mp3adora struct {
    mp3adorahandler mp3adorahandler
}


func newmp3adora(mp3adorahandler mp3adorahandler) *mp3adora {
    return &mp3adora{ mp3adorahandler: mp3adorahandler }
}


// See:  http://mutagen-specs.readthedocs.io/en/latest/apev2/apev2.html
//       http://wiki.hydrogenaud.io/index.php?title=APEv2_specification
func (m *mp3adora) parseape(reader io.Reader) (size int, err error) {
    var n int

    // First sixteen bytes are:
    // 0:        'A'.
    // 1:        'P'.
    // 2:        'E'.
    // 3:        'T'.
    // 4:        'A'.
    // 5:        'G'.
    // 6:        'E'.
    // 7:        'X'.
    // 8..11:    version.
    // 12..15:   size.

    bytes16 := make([]byte, 16)
    if n, err = io.ReadFull(reader, bytes16); n != 16 || err != nil {
        return 0, err
    }

    size = int(bytes16[12]) |
           int(bytes16[13]) << 8 |
           int(bytes16[14]) << 16 |
           int(bytes16[15]) << 24
    size += 32

    bytes := make([]byte, size)
    copy(bytes, bytes16)
    if n, err = io.ReadFull(reader, bytes[16:]); n != size - 16 || err != nil {
        return 0, err
    }

    if err = m.mp3adorahandler.processape(bytes); err != nil {
        return size, err
    }

    return size, nil
}


func (m *mp3adora) parseid3v1(reader io.Reader) (size int, err error) {
    var n int

    size = 128

    bytes := make([]byte, size)
    if n, err = io.ReadFull(reader, bytes); n != size || err != nil {
        return 0, err
    }

    if err = m.mp3adorahandler.processid3v1(bytes); err != nil {
        return size, err
    }

    return size, nil
}


func (m *mp3adora) parseid3v2(reader io.Reader) (size int, err error) {
    var n int

    // First ten bytes are:
    // 0:        'I'.
    // 1:        'D'.
    // 2:        '3'.
    // 3:        version.
    // 4:        revision.
    // 5:        flags.
    // 6..9:     size.

    bytes10 := make([]byte, 10)
    if n, err = io.ReadFull(reader, bytes10); n != 10 || err != nil {
        return 0, err
    }

    // See:  http://www.ulduzsoft.com/2012/07/parsing-id3v2-tags-in-the-mp3-files/
    size = int(bytes10[9]) |
           int(bytes10[8]) << 7 |
           int(bytes10[7]) << 14 |
           int(bytes10[6]) << 21
    size += 10

    bytes := make([]byte, size)
    copy(bytes, bytes10)
    if n, err = io.ReadFull(reader, bytes[10:]); n != size - 10 || err != nil {
        return 0, err
    }

    if err = m.mp3adorahandler.processid3v2(bytes); err != nil {
        return size, err
    }

    return size, nil
}


func (m *mp3adora) parsemp3frame(reader io.Reader) (size int, err error) {
    var n int

    bytes4 := make([]byte, 4)
    if n, err = io.ReadFull(reader, bytes4); n != 4 || err != nil {
        return 0, err
    }

    header := binary.BigEndian.Uint32(bytes4)

    header >>= 9
    paddingbit := header & 0x01
    header >>= 1
    samplingrate := header & 0x03
    header >>= 2
    bitrate := header & 0x0f
    header >>= 5
    layer := header & 0x03
    header >>= 2
    audioversion := header & 0x03
    header >>= 2
    framesync := header & 0x7ff
    header >>= 11

    if framesync != 0x7ff || header != 0 {
        return 0, fmt.Errorf("Unable to parse mp3 frame header.")
    }

    bitratevalues := [][]int {
        //  V1,L1   V1,L2    V1,L3    V2,L1    V2, L2 & L3.
        {
            0,      0,       0,       0,       0,           // 0000.
        },
        {
            32,     32,      32,      32,      8,           // 0001.
        },
        {
            64,     48,      40,      48,      16,          // 0010.
        },
        {
            96,     56,      48,      56,      24,          // 0011.
        },
        {
            128,    64,      56,      64,      32,          // 0100.
        },
        {
            160,    80,      64,      80,      40,          // 0101.
        },
        {
            192,    96,      80,      96,      48,          // 0110.
        },
        {
            224,    112,     96,      112,     56,          // 0111.
        },
        {
            256,    128,     112,     128,     64,          // 1000.
        },
        {
            288,    160,     128,     144,     80,          // 1001.
        },
        {
            320,    192,     160,     160,     96,          // 1010.
        },
        {
            352,    224,     192,     176,     112,         // 1011.
        },
        {
            384,    256,     224,     192,     128,         // 1100.
        },
        {
            416,    320,     256,     224,     144,         // 1101.
        },
        {
            448,    384,     320,     256,     160,         // 1110.
        },
        {
            -1,     -1,      -1,      -1,      -1,          // 1111.
        },
    }

    samplingratevalues := [][]int {
        //  MPEG1   MPEG2    MPEG2.5
        {
            44100,  22050,   11025,                         // 00.
        },
        {
            48000,  24000,   12000,                         // 01.
        },
        {
            32000,  16000,    8000,                         // 10.
        },
        {
                0,      0,       0,                         // 11.
        },
    }

    var bitratecolumn int
    if audioversion == 0x03 && layer == 0x03 {
        bitratecolumn = 0
    } else if audioversion == 0x03 && layer == 0x02 {
        bitratecolumn = 1
    } else if audioversion == 0x03 && layer == 0x01 {
        bitratecolumn = 2
    } else if audioversion == 0x02 && layer == 0x03 {
        bitratecolumn = 3
    } else if (audioversion == 0x02 && layer == 0x02) ||
              (audioversion == 0x02 && layer == 0x01) {
        bitratecolumn = 4
    } else {
        return 0, fmt.Errorf("Unable to find mp3 bitrate.")
    }

    var samplingratecolumn int
    if audioversion == 0x03 {
        samplingratecolumn = 0
    } else if audioversion == 0x02 {
        samplingratecolumn = 1
    } else if audioversion == 0x00 {
        samplingratecolumn = 2
    } else {
        return 0, fmt.Errorf("Unable to find mp3 sampling rate.")
    }

    size = 144 * bitratevalues[bitrate][bitratecolumn]
    size *= 1000
    size /= samplingratevalues[samplingrate][samplingratecolumn]
    size += int(paddingbit)

    bytes := make([]byte, size)
    copy(bytes, bytes4)
    if n, err = io.ReadFull(reader, bytes[4:]); n != size - 4 || err != nil {
        return 0, err
    }

    if err = m.mp3adorahandler.processmp3frame(bytes); err != nil {
        return size, err
    }

    return size, nil
}


func (m *mp3adora) parse(reader io.Reader) (size int, err error) {
    bufferedreader := bufio.NewReader(reader)

    var bytes []byte
    var sizeframe int
    for true {
        if bytes, err = bufferedreader.Peek(3); err != nil {
            break
        }

        if string(bytes) == "TAG" {
            if sizeframe, err = m.parseid3v1(bufferedreader); err != nil {
                break
            }
            size += sizeframe
            continue
        }

        if string(bytes) == "ID3" {
            if sizeframe, err = m.parseid3v2(bufferedreader); err != nil {
                break
            }
            size += sizeframe
            continue
        }

        if bytes, err = bufferedreader.Peek(4); err != nil {
            break
        }

        if bytes[0] == 0xff && bytes[1] & 0xe0 == 0xe0 {
            if sizeframe, err = m.parsemp3frame(bufferedreader); err != nil {
                break
            }
            size += sizeframe
            continue
        }

        if bytes, err = bufferedreader.Peek(8); err != nil {
            break
        }

        if string(bytes) == "APETAGEX" {
            if sizeframe, err = m.parseape(bufferedreader); err != nil {
                break
            }
            size += sizeframe
            continue
        }

        bytes1 := make([]byte, 1)
        if n, err := io.ReadFull(bufferedreader, bytes1); n != 1 || err != nil {
            return 0, err
        }
        if err = m.mp3adorahandler.processunrecognised(bytes1[0]); err != nil {
            return size, err
        }
        size++
    }

    if err != io.EOF {
        return size, err
    }

    return size, nil
}
