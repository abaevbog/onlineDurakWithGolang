// Code generated by "stringer -type=Suit,Rank"; DO NOT EDIT.

package deck

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Spade-0]
	_ = x[Diamond-1]
	_ = x[Club-2]
	_ = x[Heart-3]
}

const _Suit_name = "SpadeDiamondClubHeart"

var _Suit_index = [...]uint8{0, 5, 12, 16, 21, 26}

func (i Suit) String() string {
	if i >= Suit(len(_Suit_index)-1) {
		return "Suit(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Suit_name[_Suit_index[i]:_Suit_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Six-6]
	_ = x[Seven-7]
	_ = x[Eight-8]
	_ = x[Nine-9]
	_ = x[Ten-10]
	_ = x[Jack-11]
	_ = x[Queen-12]
	_ = x[King-13]
	_ = x[Ace-14]
}

const _Rank_name = "SixSevenEightNineTenJackQueenKingAce"

var _Rank_index = [...]uint8{0, 3, 8, 13, 17, 20, 24, 29, 33, 36}

func (i Rank) String() string {
	i -= 6
	if i >= Rank(len(_Rank_index)-1) {
		return "Rank(" + strconv.FormatInt(int64(i+6), 10) + ")"
	}
	return _Rank_name[_Rank_index[i]:_Rank_index[i+1]]
}
