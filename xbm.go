package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
)

type XBM struct {
	width  int
	height int
	name   string
	data   []uint8
	neg    bool
}

var (
	// TODO: also support 'version 10' (pre-1986) XBMs, which are weirder
	mainRx = regexp.MustCompile(`(?sm)^#define\s+\S+_width\s+(\d+)\s*$
#define\s+\S+_height\s+(\d+)\s*$
(?:#define\s+\S+_x_hot\s+\d+$
#define\s+\S+_y_hot\s+\d+$
)?static\s+(?:unsigned\s+)?char\s+(\S+)_bits\s*\[\]\s*=\s*{\s*$
((?:\s*0x[[:xdigit:]]{2}\s*,)*\s*0x[[:xdigit:]]{1,2})\s*,?\s*\}`)
	hexRx = regexp.MustCompile(`0x[[:xdigit:]]{1,2}`)
)

func fromReader(r io.Reader, neg bool) (*XBM, error) {
	lr := io.LimitedReader{N: 1024 * 1024, R: r}
	buf, err := io.ReadAll(&lr)
	if err != nil {
		return nil, err
	}
	subs := mainRx.FindSubmatch(buf)
	if len(subs) != 5 {
		return nil, fmt.Errorf("unable to understand this particular XBM file; please report")
	}

	xbm := &XBM{name: string(subs[3]), neg: neg}

	xbm.width, err = strconv.Atoi(string(subs[1]))
	if err != nil {
		return nil, fmt.Errorf("unable to parse the width: %w", err)
	}
	xbm.height, err = strconv.Atoi(string(subs[2]))
	if err != nil {
		return nil, fmt.Errorf("unable to parse the height: %w", err)
	}

	hexes := hexRx.FindAllSubmatch(subs[4], -1)
	xbm.data = make([]uint8, len(hexes))
	for i, hex := range hexes {
		d, err := strconv.ParseUint(string(hex[0]), 0, 8)
		if err != nil {
			return nil, err
		}
		xbm.data[i] = uint8(d)
	}
	d := uint8(xbm.width & 7)
	if neg || d == 0 {
		return xbm, nil
	}
	w := (xbm.width + 7) / 8
	n := -1
	d = ^((1 << d) - 1)
	for i := 0; i < xbm.height; i++ {
		n += w
		xbm.data[n] = xbm.data[n] | d
	}

	return xbm, nil
}

func (xbm XBM) braille() string {
	var filler byte = 0xff
	if xbm.neg {
		filler = 0
	}
	w := (xbm.width + 7) / 8
	u := (xbm.width+1)/2 - 4
	buf := make([]rune, 0, (w*4+1)*(xbm.height+3)/4)
	for i := 0; i < xbm.height; i += 4 {
		for j := 0; j < w; j++ {
			b := [4]byte{xbm.data[(i*w)+j], filler, filler, filler}
			switch xbm.height {
			default:
				b[3] = xbm.data[(i+3)*w+j]
				fallthrough
			case i + 3:
				b[2] = xbm.data[(i+2)*w+j]
				fallthrough
			case i + 2:
				b[1] = xbm.data[(i+1)*w+j]
				fallthrough
			case i + 1:
				// done
			}
			switch j*4 - u {
			default:
				buf = append(buf,
					bit2brails((b[0]&0b00000011)<<0+(b[1]&0b00000011)<<2+(b[2]&0b00000011)<<4+(b[3]&0b00000011)<<6, xbm.neg),
					bit2brails((b[0]&0b00001100)>>2+(b[1]&0b00001100)<<0+(b[2]&0b00001100)<<2+(b[3]&0b00001100)<<4, xbm.neg),
					bit2brails((b[0]&0b00110000)>>4+(b[1]&0b00110000)>>2+(b[2]&0b00110000)<<0+(b[3]&0b00110000)<<2, xbm.neg),
					bit2brails((b[0]&0b11000000)>>6+(b[1]&0b11000000)>>4+(b[2]&0b11000000)>>2+(b[3]&0b11000000)<<0, xbm.neg))
			case 1:
				buf = append(buf,
					bit2brails((b[0]&0b00000011)<<0+(b[1]&0b00000011)<<2+(b[2]&0b00000011)<<4+(b[3]&0b00000011)<<6, xbm.neg),
					bit2brails((b[0]&0b00001100)>>2+(b[1]&0b00001100)<<0+(b[2]&0b00001100)<<2+(b[3]&0b00001100)<<4, xbm.neg),
					bit2brails((b[0]&0b00110000)>>4+(b[1]&0b00110000)>>2+(b[2]&0b00110000)<<0+(b[3]&0b00110000)<<2, xbm.neg))
			case 2:
				buf = append(buf,
					bit2brails((b[0]&0b00000011)<<0+(b[1]&0b00000011)<<2+(b[2]&0b00000011)<<4+(b[3]&0b00000011)<<6, xbm.neg),
					bit2brails((b[0]&0b00001100)>>2+(b[1]&0b00001100)<<0+(b[2]&0b00001100)<<2+(b[3]&0b00001100)<<4, xbm.neg))
			case 3:
				buf = append(buf,
					bit2brails((b[0]&0b00000011)<<0+(b[1]&0b00000011)<<2+(b[2]&0b00000011)<<4+(b[3]&0b00000011)<<6, xbm.neg))
			}
		}
		buf = append(buf, '\n')
	}
	return string(buf)
}

const brailidxs = "" +
	"\xff\xfe\xf7\xf6\xfd\xfc\xf5\xf4\xef\xee\xe7\xe6\xed\xec\xe5\xe4" +
	"\xfb\xfa\xf3\xf2\xf9\xf8\xf1\xf0\xeb\xea\xe3\xe2\xe9\xe8\xe1\xe0" +
	"\xdf\xde\xd7\xd6\xdd\xdc\xd5\xd4\xcf\xce\xc7\xc6\xcd\xcc\xc5\xc4" +
	"\xdb\xda\xd3\xd2\xd9\xd8\xd1\xd0\xcb\xca\xc3\xc2\xc9\xc8\xc1\xc0" +
	"\xbf\xbe\xb7\xb6\xbd\xbc\xb5\xb4\xaf\xae\xa7\xa6\xad\xac\xa5\xa4" +
	"\xbb\xba\xb3\xb2\xb9\xb8\xb1\xb0\xab\xaa\xa3\xa2\xa9\xa8\xa1\xa0" +
	"\x9f\x9e\x97\x96\x9d\x9c\x95\x94\x8f\x8e\x87\x86\x8d\x8c\x85\x84" +
	"\x9b\x9a\x93\x92\x99\x98\x91\x90\x8b\x8a\x83\x82\x89\x88\x81\x80" +
	"\x7f~wv}|utongfmled{zsryxqpkjcbiha`_^WV]\\UTONGFMLED[ZSRYXQPKJCBI" +
	"HA@?>76=<54/.'&-,%$;:329810+*#\")(! \x1f\x1e\x17\x16\x1d\x1c\x15" +
	"\x14\x0f\x0e\x07\x06\r\x0c\x05\x04\x1b\x1a\x13\x12\x19\x18\x11\x10" +
	"\x0b\n\x03\x02\t\x08\x01\x00"

func bit2brails(b uint8, neg bool) rune {
	if neg {
		b = ^b
	}
	return 10240 + rune(brailidxs[b])
}
