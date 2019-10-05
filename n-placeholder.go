package dorm

import (
	"bytes"
)

const (
	placeHolder         = "(?)"
	twoPlaceHolder      = "(?,?)"
	threePlaceHolder    = "(?,?,?)"
	fourPlaceHolder     = "(?,?,?.?)"
	fivePlaceHolder     = "(?,?,?,?,?)"
	sixPlaceHolder      = "(?,?,?,?,?,?)"
	sevenPlaceHolder    = "(?,?,?,?,?,?,?)"
	eightPlaceHolder    = "(?,?,?,?,?,?,?,?)"
	nightPlaceHolder    = "(?,?,?,?,?,?,?,?,?)"
	tenPlaceHolder      = "(?,?,?,?,?,?,?,?,?,?)"
	elevenPlaceHolder   = "(?,?,?,?,?,?,?,?,?,?,?)"
	twelvePlaceHolder   = "(?,?,?,?,?,?,?,?,?,?,?,?)"
	thirteenPlaceHolder = "(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	fourteenPlaceHolder = "(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	fifteenPlaceHolder  = "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	sixteenPlaceHolder  = "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
)

var placeholders = make([]string, 0, 512)

var purePlaceHolder = make([]string, 9)

func init() {
	purePlaceHolder[0] = "?"
	for i := 1; i <= 8; i++ {
		purePlaceHolder[i] = purePlaceHolder[i-1] + "," + purePlaceHolder[i-1]
	}

	placeholders = append(placeholders,
		"()",
		placeHolder,
		twoPlaceHolder,
		threePlaceHolder,
		fourPlaceHolder,
		fivePlaceHolder,
		sixPlaceHolder,
		sevenPlaceHolder,
		eightPlaceHolder,
		nightPlaceHolder,
		tenPlaceHolder,
		elevenPlaceHolder,
		twelvePlaceHolder,
		thirteenPlaceHolder,
		fourteenPlaceHolder,
		fifteenPlaceHolder,
		sixteenPlaceHolder,
	)
	placeholders = append(placeholders, make([]string, 512 - 17)...)
}
func generateNPlaceHolder(n int) string {
	var resBuf = bytes.NewBuffer(make([]byte, 0, (n<<1)+1))
	resBuf.WriteByte('(')

	for i := uint8(0); i <= 8; i++ {
		if (n & (1 << i)) != 0 {
			if n & (-n) != (1 << i) {
				resBuf.WriteByte(',')
			}
			resBuf.WriteString(purePlaceHolder[i])
		}
	}

	m := n & 0x1ff
	n = m ^ n

	for n > 0 {
		if m != 0 {
			resBuf.WriteByte(',')
		} else {
			m = 1
		}
		resBuf.WriteString(purePlaceHolder[8])
		n -= 0x100
	}

	resBuf.WriteByte(')')
	return resBuf.String()
}
func nPlaceHolder(n int) string {
	switch {
	case 0 <= n && n <= 16:
		return placeholders[n]
	case n < 0:
		return "()"
	case n >= 512:
		return generateNPlaceHolder(n)
	default:
		if len(placeholders[n]) != 0 {
			return placeholders[n]
		}
		placeholders[n] = generateNPlaceHolder(n)
		return placeholders[n]
	}
}
