// 'mp3header.go'.
// Chris Shiels.


package main


import (
    "encoding/binary"
    "fmt"
)


type mp3header struct {
    audioversion float32
    layer int
    protection bool
    bitrate int
    samplingrate int
    padding bool
    private bool
    channelmode int
    modeextension int
    copyright bool
    original bool
    emphasis int
    size int
}


func newmp3headerfrombytes(bytes []byte) (m *mp3header, err error) {

    // Sign:  Length:  Position:  Description:
    //        (bits)   (bits)
    // A      11       (31-21)    Frame sync - all bits must be set.
    // B      2        (20,19)    MPEG Audio version ID.
    // C      2        (18,17)    Layer description.
    // D      1        (16)       Protection bit.
    // E      4        (15,12)    Bitrate index.
    // F      2        (11,10)    Sampling rate frequency index.
    // G      1        (9)        Padding bit.
    // H      1        (8)        Private bit - this one is only informative.
    // I      2        (7,6)      Channel mode.
    // J      2        (5,4)      Mode extension - only used in joint stereo.
    // K      1        (3)        Copyright.
    // L      1        (2)        Original.
    // M      2        (1,0)      Emphasis.

    m = new(mp3header)

    header := binary.BigEndian.Uint32(bytes[0:4])

    if header >> 21 != 0x7ff {
        return nil, fmt.Errorf("Unable to find mp3 frame header.")
    }

    emphasis := header & 0x03
    header >>= 2

    original := header & 0x01
    header >>= 1

    copyright := header & 0x01
    header >>= 1

    modeextension := header & 0x03
    header >>= 2

    channelmode := header & 0x03
    header >>= 2

    private := header & 0x01
    header >>= 1

    paddingbit := header & 0x01
    header >>= 1

    samplingrate := header & 0x03
    header >>= 2

    bitrate := header & 0x0f
    header >>= 4

    protection := header & 0x01
    header >>= 1

    layer := header & 0x03
    header >>= 2

    audioversion := header & 0x03
    header >>= 2

    framesync := header & 0x7ff
    header >>= 11

    if framesync != 0x7ff || header != 0 {
        return nil, fmt.Errorf("Unable to find mp3 frame header.")
    }

    audioversionvalues := []float32 {
        2.5,                                                // 00.
        -1,                                                 // 01.
        2,                                                  // 10.
        1,                                                  // 11.
    }

    layervalues := []int {
        -1,                                                 // 00.
        3,                                                  // 01.
        2,                                                  // 10.
        1,                                                  // 11.
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
        return nil, fmt.Errorf("Unable to find mp3 bitrate.")
    }

    var samplingratecolumn int
    if audioversion == 0x03 {
        samplingratecolumn = 0
    } else if audioversion == 0x02 {
        samplingratecolumn = 1
    } else if audioversion == 0x00 {
        samplingratecolumn = 2
    } else {
        return nil, fmt.Errorf("Unable to find mp3 sampling rate.")
    }

    m.audioversion = audioversionvalues[audioversion]
    m.layer = layervalues[layer]
    m.protection = protection == 0x01
    m.bitrate = bitratevalues[bitrate][bitratecolumn]
    m.samplingrate = samplingratevalues[samplingrate][samplingratecolumn]
    m.padding = paddingbit == 0x01
    m.private = private == 0x01
    m.channelmode = int(channelmode)
    m.modeextension = int(modeextension)
    m.copyright = copyright == 0x01
    m.original = original == 0x01
    m.emphasis = int(emphasis)
    m.size = 144 * bitratevalues[bitrate][bitratecolumn]
    m.size *= 1000
    m.size /= samplingratevalues[samplingrate][samplingratecolumn]
    m.size += int(paddingbit)

    return m, nil
}
