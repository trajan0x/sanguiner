// Code generated by "stringer -type=contractTypeImpl -linecomment"; DO NOT EDIT.

package testutil

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SynapseCCTPType-1]
	_ = x[MockMessageTransmitterType-2]
	_ = x[MockTokenMessengerType-3]
}

const _contractTypeImpl_name = "SynapseCCTPMockMessageTransmitterMockTokenMessenger"

var _contractTypeImpl_index = [...]uint8{0, 11, 33, 51}

func (i contractTypeImpl) String() string {
	i -= 1
	if i < 0 || i >= contractTypeImpl(len(_contractTypeImpl_index)-1) {
		return "contractTypeImpl(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _contractTypeImpl_name[_contractTypeImpl_index[i]:_contractTypeImpl_index[i+1]]
}
