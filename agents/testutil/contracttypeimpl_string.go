// Code generated by "stringer -type=contractTypeImpl -linecomment"; DO NOT EDIT.

package testutil

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OriginType-0]
	_ = x[MessageHarnessType-1]
	_ = x[OriginHarnessType-2]
	_ = x[AttestationHarnessType-3]
	_ = x[TipsHarnessType-4]
	_ = x[HeaderHarnessType-5]
	_ = x[DestinationHarnessType-6]
	_ = x[AttestationCollectorType-7]
	_ = x[DestinationType-8]
	_ = x[AgentsTestContractType-9]
	_ = x[TestClientType-10]
}

const _contractTypeImpl_name = "OriginMessageHarnessOriginHarnessAttestationHarnessTypeTipsHarnessTypeHeaderHarnessTypeDestinationHarnessAttestationCollectorDestinationAgentsTestContractTestClient"

var _contractTypeImpl_index = [...]uint8{0, 6, 20, 33, 55, 70, 87, 105, 125, 136, 154, 164}

func (i contractTypeImpl) String() string {
	if i < 0 || i >= contractTypeImpl(len(_contractTypeImpl_index)-1) {
		return "contractTypeImpl(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _contractTypeImpl_name[_contractTypeImpl_index[i]:_contractTypeImpl_index[i+1]]
}
