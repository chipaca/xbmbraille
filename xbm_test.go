package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"chipaca.com/intmath"
	"github.com/go-test/deep"
)

// generated via
//
//	for x in 0 1 2 3; do for y in 0 1 2 3; do
//	  echo "{$x,$y}: {$( convert -size 4x4 xc:white -fill black -draw "point $x,$y" XBM:- | grep 0x00)"
//	done; done
var testPixels = map[[2]int][4]uint8{
	{0, 0}: {0x01, 0x00, 0x00, 0x00},
	{0, 1}: {0x00, 0x01, 0x00, 0x00},
	{0, 2}: {0x00, 0x00, 0x01, 0x00},
	{0, 3}: {0x00, 0x00, 0x00, 0x01},
	{1, 0}: {0x02, 0x00, 0x00, 0x00},
	{1, 1}: {0x00, 0x02, 0x00, 0x00},
	{1, 2}: {0x00, 0x00, 0x02, 0x00},
	{1, 3}: {0x00, 0x00, 0x00, 0x02},
	{2, 0}: {0x04, 0x00, 0x00, 0x00},
	{2, 1}: {0x00, 0x04, 0x00, 0x00},
	{2, 2}: {0x00, 0x00, 0x04, 0x00},
	{2, 3}: {0x00, 0x00, 0x00, 0x04},
	{3, 0}: {0x08, 0x00, 0x00, 0x00},
	{3, 1}: {0x00, 0x08, 0x00, 0x00},
	{3, 2}: {0x00, 0x00, 0x08, 0x00},
	{3, 3}: {0x00, 0x00, 0x00, 0x08},
}

func testXBM(xy [2]int, data [4]uint8) (string, string) {
	name := fmt.Sprintf("%dx%d", xy[0], xy[1])
	myXBM := fmt.Sprintf(`#define %v_width 4
#define %[1]v_height 4
static char %[1]v_bits[] = {
%#02x, %#02x, %#02x, %#02x, };`,
		name, data[0], data[1], data[2], data[3])
	return name, myXBM
}

// test that fromReader gets the right bits
func TestFromReader(t *testing.T) {
	for xy, data := range testPixels {
		name, myXBM := testXBM(xy, data)
		t.Run(name, func(t *testing.T) {
			xbm, err := fromReader(strings.NewReader(myXBM), true)
			if err != nil {
				t.Fatal(err)
			}
			if xbm.width != 4 {
				t.Errorf("wanted width 8, got %d", xbm.width)
			}
			if xbm.height != 4 {
				t.Errorf("wanted height 8, got %d", xbm.height)
			}
			if len(xbm.data) != 4 {
				t.Fatalf("wanted data of length 8, got %d", len(xbm.data))
			}
			if data != [4]uint8(xbm.data) {
				t.Errorf("wanted data of %v, got %v", data, xbm.data)
			}
		})
	}
}

func FuzzFromReader(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) == 0 {
			return
		}
		w := int(intmath.Sqrt(uint(len(data))))
		h := len(data) / w
		var buf bytes.Buffer
		fmt.Fprintf(&buf, `#define fuzz_width %d
#define fuzz_height %d
static char fizz_bits[] = {
`, w, h)
		udata := make([]uint8, len(data))
		for i := range data {
			udata[i] = uint8(data[i])
			fmt.Fprintf(&buf, "%#02x, ", udata[i])
		}
		fmt.Fprintln(&buf, "};")
		xbm, err := fromReader(&buf, true)
		if err != nil {
			t.Fatal(err)
		}
		if xbm.width != w {
			t.Errorf("wanted width %d, got %d", w, xbm.width)
		}
		if xbm.height != h {
			t.Errorf("wanted height %d, got %d", h, xbm.height)
		}
		if len(xbm.data) != len(data) {
			t.Fatalf("wanted data of length %d, got %d", len(data), len(xbm.data))
		}
		if diff := deep.Equal(udata, xbm.data); diff != nil {
			t.Errorf("wanted data of %v, got %v", udata, xbm.data)
			t.Error(diff)
		}
	})
}

var testDots = map[int][]string{
	0: {"⣾⣿", "⣽⣿", "⣻⣿", "⢿⣿"},
	1: {"⣷⣿", "⣯⣿", "⣟⣿", "⡿⣿"},
	2: {"⣿⣾", "⣿⣽", "⣿⣻", "⣿⢿"},
	3: {"⣿⣷", "⣿⣯", "⣿⣟", "⣿⡿"},
}

func TestBraille(t *testing.T) {
	for xy, data := range testPixels {
		name, myXBM := testXBM(xy, data)
		t.Run(name, func(t *testing.T) {
			xbm, err := fromReader(strings.NewReader(myXBM), false)
			if err != nil {
				t.Fatal(err)
			}
			gotDots := xbm.braille()
			wantDots := testDots[xy[0]][xy[1]] + "\n"
			if gotDots != wantDots {
				t.Errorf("wanted %q, got %q", wantDots, gotDots)
			}
		})
	}
}

var testNegDots = map[int][]string{
	0: {"⠁⠀", "⠂⠀", "⠄⠀", "⡀⠀"},
	1: {"⠈⠀", "⠐⠀", "⠠⠀", "⢀⠀"},
	2: {"⠀⠁", "⠀⠂", "⠀⠄", "⠀⡀"},
	3: {"⠀⠈", "⠀⠐", "⠀⠠", "⠀⢀"},
}

func TestNegBraille(t *testing.T) {
	for xy, data := range testPixels {
		name, myXBM := testXBM(xy, data)
		t.Run(name, func(t *testing.T) {
			xbm, err := fromReader(strings.NewReader(myXBM), true)
			if err != nil {
				t.Fatal(err)
			}
			gotDots := xbm.braille()
			wantDots := testNegDots[xy[0]][xy[1]] + "\n"
			if gotDots != wantDots {
				t.Errorf("wanted %q, got %q", wantDots, gotDots)
			}
		})
	}
}

func testOdd(size int, data []uint8, neg bool, want string) func(*testing.T) {
	return func(t *testing.T) {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "#define odd_width %d\n#define odd_height %[1]d\nstatic char odd_bits[] = {\n", size)
		for i := range data {
			fmt.Fprintf(&buf, "%#02x, ", data[i])
		}
		fmt.Fprintln(&buf, "};")
		xbm, err := fromReader(&buf, neg)
		if err != nil {
			t.Fatal(err)
		}
		got := strings.Trim(xbm.braille(), "\n")
		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	}
}

func TestOddSizes(t *testing.T) {
	type TC struct {
		data    []uint8
		wantPos string
		wantNeg string
	}
	for x, tc := range map[int]TC{
		2: {
			data:    []uint8{0x01, 0x02},
			wantPos: "⠊",
			wantNeg: "⠑",
		},
		3: {
			data:    []uint8{0x05, 0x02, 0x05},
			wantPos: "⠪⠂",
			wantNeg: "⠕⠅",
		},
		4: {
			data:    []uint8{0x05, 0x0A, 0x05, 0x0A},
			wantPos: "⡪⡪",
			wantNeg: "⢕⢕",
		},
		5: {
			data:    []uint8{0x15, 0x0A, 0x15, 0x0A, 0x15},
			wantPos: "⡪⡪⡂\n⠈⠈⠀",
			wantNeg: "⢕⢕⠅\n⠁⠁⠁",
		},
		6: {
			data:    []uint8{0x15, 0x2A, 0x15, 0x2A, 0x15, 0x2A},
			wantPos: "⡪⡪⡪\n⠊⠊⠊",
			wantNeg: "⢕⢕⢕\n⠑⠑⠑",
		},
		7: {
			data:    []uint8{0x55, 0x2A, 0x55, 0x2A, 0x55, 0x2A, 0x55},
			wantPos: "⡪⡪⡪⡂\n⠪⠪⠪⠂",
			wantNeg: "⢕⢕⢕⠅\n⠕⠕⠕⠅",
		},
	} {
		t.Run(fmt.Sprintf("%d/pos", x), testOdd(x, tc.data, false, tc.wantPos))
		t.Run(fmt.Sprintf("%d/neg", x), testOdd(x, tc.data, true, tc.wantNeg))
	}
}

var brailidxv = [...]rune{
	'⣿', '⣾', '⣷', '⣶', '⣽', '⣼', '⣵', '⣴', '⣯', '⣮', '⣧', '⣦', '⣭', '⣬', '⣥', '⣤',
	'⣻', '⣺', '⣳', '⣲', '⣹', '⣸', '⣱', '⣰', '⣫', '⣪', '⣣', '⣢', '⣩', '⣨', '⣡', '⣠',
	'⣟', '⣞', '⣗', '⣖', '⣝', '⣜', '⣕', '⣔', '⣏', '⣎', '⣇', '⣆', '⣍', '⣌', '⣅', '⣄',
	'⣛', '⣚', '⣓', '⣒', '⣙', '⣘', '⣑', '⣐', '⣋', '⣊', '⣃', '⣂', '⣉', '⣈', '⣁', '⣀',
	'⢿', '⢾', '⢷', '⢶', '⢽', '⢼', '⢵', '⢴', '⢯', '⢮', '⢧', '⢦', '⢭', '⢬', '⢥', '⢤',
	'⢻', '⢺', '⢳', '⢲', '⢹', '⢸', '⢱', '⢰', '⢫', '⢪', '⢣', '⢢', '⢩', '⢨', '⢡', '⢠',
	'⢟', '⢞', '⢗', '⢖', '⢝', '⢜', '⢕', '⢔', '⢏', '⢎', '⢇', '⢆', '⢍', '⢌', '⢅', '⢄',
	'⢛', '⢚', '⢓', '⢒', '⢙', '⢘', '⢑', '⢐', '⢋', '⢊', '⢃', '⢂', '⢉', '⢈', '⢁', '⢀',
	'⡿', '⡾', '⡷', '⡶', '⡽', '⡼', '⡵', '⡴', '⡯', '⡮', '⡧', '⡦', '⡭', '⡬', '⡥', '⡤',
	'⡻', '⡺', '⡳', '⡲', '⡹', '⡸', '⡱', '⡰', '⡫', '⡪', '⡣', '⡢', '⡩', '⡨', '⡡', '⡠',
	'⡟', '⡞', '⡗', '⡖', '⡝', '⡜', '⡕', '⡔', '⡏', '⡎', '⡇', '⡆', '⡍', '⡌', '⡅', '⡄',
	'⡛', '⡚', '⡓', '⡒', '⡙', '⡘', '⡑', '⡐', '⡋', '⡊', '⡃', '⡂', '⡉', '⡈', '⡁', '⡀',
	'⠿', '⠾', '⠷', '⠶', '⠽', '⠼', '⠵', '⠴', '⠯', '⠮', '⠧', '⠦', '⠭', '⠬', '⠥', '⠤',
	'⠻', '⠺', '⠳', '⠲', '⠹', '⠸', '⠱', '⠰', '⠫', '⠪', '⠣', '⠢', '⠩', '⠨', '⠡', '⠠',
	'⠟', '⠞', '⠗', '⠖', '⠝', '⠜', '⠕', '⠔', '⠏', '⠎', '⠇', '⠆', '⠍', '⠌', '⠅', '⠄',
	'⠛', '⠚', '⠓', '⠒', '⠙', '⠘', '⠑', '⠐', '⠋', '⠊', '⠃', '⠂', '⠉', '⠈', '⠁', '⠀',
}

func bit2brailv(b uint8, neg bool) rune {
	if neg {
		b = ^b
	}
	return brailidxv[b]
}

const brailidxr = `
⣿ ⣾ ⣷ ⣶ ⣽ ⣼ ⣵ ⣴ ⣯ ⣮ ⣧ ⣦ ⣭ ⣬ ⣥ ⣤ ⣻ ⣺ ⣳ ⣲ ⣹ ⣸ ⣱ ⣰ ⣫ ⣪ ⣣ ⣢ ⣩ ⣨ ⣡ ⣠
⣟ ⣞ ⣗ ⣖ ⣝ ⣜ ⣕ ⣔ ⣏ ⣎ ⣇ ⣆ ⣍ ⣌ ⣅ ⣄ ⣛ ⣚ ⣓ ⣒ ⣙ ⣘ ⣑ ⣐ ⣋ ⣊ ⣃ ⣂ ⣉ ⣈ ⣁ ⣀
⢿ ⢾ ⢷ ⢶ ⢽ ⢼ ⢵ ⢴ ⢯ ⢮ ⢧ ⢦ ⢭ ⢬ ⢥ ⢤ ⢻ ⢺ ⢳ ⢲ ⢹ ⢸ ⢱ ⢰ ⢫ ⢪ ⢣ ⢢ ⢩ ⢨ ⢡ ⢠
⢟ ⢞ ⢗ ⢖ ⢝ ⢜ ⢕ ⢔ ⢏ ⢎ ⢇ ⢆ ⢍ ⢌ ⢅ ⢄ ⢛ ⢚ ⢓ ⢒ ⢙ ⢘ ⢑ ⢐ ⢋ ⢊ ⢃ ⢂ ⢉ ⢈ ⢁ ⢀
⡿ ⡾ ⡷ ⡶ ⡽ ⡼ ⡵ ⡴ ⡯ ⡮ ⡧ ⡦ ⡭ ⡬ ⡥ ⡤ ⡻ ⡺ ⡳ ⡲ ⡹ ⡸ ⡱ ⡰ ⡫ ⡪ ⡣ ⡢ ⡩ ⡨ ⡡ ⡠
⡟ ⡞ ⡗ ⡖ ⡝ ⡜ ⡕ ⡔ ⡏ ⡎ ⡇ ⡆ ⡍ ⡌ ⡅ ⡄ ⡛ ⡚ ⡓ ⡒ ⡙ ⡘ ⡑ ⡐ ⡋ ⡊ ⡃ ⡂ ⡉ ⡈ ⡁ ⡀
⠿ ⠾ ⠷ ⠶ ⠽ ⠼ ⠵ ⠴ ⠯ ⠮ ⠧ ⠦ ⠭ ⠬ ⠥ ⠤ ⠻ ⠺ ⠳ ⠲ ⠹ ⠸ ⠱ ⠰ ⠫ ⠪ ⠣ ⠢ ⠩ ⠨ ⠡ ⠠
⠟ ⠞ ⠗ ⠖ ⠝ ⠜ ⠕ ⠔ ⠏ ⠎ ⠇ ⠆ ⠍ ⠌ ⠅ ⠄ ⠛ ⠚ ⠓ ⠒ ⠙ ⠘ ⠑ ⠐ ⠋ ⠊ ⠃ ⠂ ⠉ ⠈ ⠁ ⠀
`

func bit2brailr(b uint8, neg bool) rune {
	if neg {
		b = ^b
	}
	n := 4*int(b) + 1
	if n < 0 || n > len(brailidxr) {
		return '⠀'
	}
	r, _ := utf8.DecodeRuneInString(brailidxr[n:])
	return r
}

func benchmarkF(b *testing.B, f func(uint8, bool) rune) {
	for i := 0; i < b.N; i++ {
		for u := uint8(0); u < 255; u++ {
			_ = f(u, true)
			_ = f(u, false)
		}
	}
}

func BenchmarkStringAndMath(b *testing.B) { benchmarkF(b, bit2brails) }
func BenchmarkRuneSlice(b *testing.B)     { benchmarkF(b, bit2brailv) }
func BenchmarkRuneConst(b *testing.B)     { benchmarkF(b, bit2brailr) }
