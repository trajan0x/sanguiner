// Code generated by "stringer -type=contractTypeImpl -linecomment"; DO NOT EDIT.

package testutil

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TestContractType-0]
	_ = x[TestContract2Type-1]
}

const _contractTypeImpl_name = "TestContractTypeTestContract2Type"

var _contractTypeImpl_index = [...]uint8{0, 16, 33}

func (i contractTypeImpl) String() string {
	if i < 0 || i >= contractTypeImpl(len(_contractTypeImpl_index)-1) {
		return "contractTypeImpl(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _contractTypeImpl_name[_contractTypeImpl_index[i]:_contractTypeImpl_index[i+1]]
}
