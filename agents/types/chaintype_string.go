// Code generated by "stringer -type=ChainType -linecomment"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EVM-0]
}

const _ChainType_name = "EVM"

var _ChainType_index = [...]uint8{0, 3}

func (i ChainType) String() string {
	if i >= ChainType(len(_ChainType_index)-1) {
		return "ChainType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ChainType_name[_ChainType_index[i]:_ChainType_index[i+1]]
}
