// Code generated by "stringer -type=ContractType -linecomment"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OriginContract-1]
	_ = x[DestinationContract-2]
	_ = x[SummitContract-3]
	_ = x[Other-4]
}

const _ContractType_name = "OriginContractDestinationContractSummitContractOther"

var _ContractType_index = [...]uint8{0, 14, 33, 47, 52}

func (i ContractType) String() string {
	i -= 1
	if i < 0 || i >= ContractType(len(_ContractType_index)-1) {
		return "ContractType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ContractType_name[_ContractType_index[i]:_ContractType_index[i+1]]
}
