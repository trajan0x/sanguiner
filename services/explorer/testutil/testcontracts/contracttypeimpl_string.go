// Code generated by "stringer -type=contractTypeImpl -linecomment"; DO NOT EDIT.

package testcontracts

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TestBridgeConfigTypeV3-0]
	_ = x[TestSynapseBridgeType-1]
	_ = x[TestSwapFlashLoanType-2]
	_ = x[TestSynapseBridgeV1Type-3]
	_ = x[TestMessageBusType-4]
	_ = x[TestMetaSwapType-5]
	_ = x[TestCCTPType-6]
}

const _contractTypeImpl_name = "TestBridgeConfigTypeV3TestSynapseBridgeTypeTestSwapFlashLoanTypeTestSynapseBridgeV1TypeTestMessageBusTypeTestMetaSwapTypeTestCCTPType"

var _contractTypeImpl_index = [...]uint8{0, 22, 43, 64, 87, 105, 121, 133}

func (i contractTypeImpl) String() string {
	if i < 0 || i >= contractTypeImpl(len(_contractTypeImpl_index)-1) {
		return "contractTypeImpl(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _contractTypeImpl_name[_contractTypeImpl_index[i]:_contractTypeImpl_index[i+1]]
}
