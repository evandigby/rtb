package rtb

import (
	"testing"
)

// TestCpmToMicroCents ensures that we are converting from CPM to MicroCents correctly.
func TestCpmToMicroCents(t *testing.T) {
	cpm := 0.32
	expected := int64(32000000)
	actual := CpmToMicroCents(cpm)

	if expected != actual {
		t.Fail()
	}
}

// TestMicroCentsToCpm ensures that we are converting from MicroCents to CPM correctly.
func TestMicroCentsToCpm(t *testing.T) {
	microcents := int64(32000000)
	expected := 0.32
	actual := MicroCentsToCpm(microcents)

	if expected != actual {
		t.Fail()
	}
}

// TestMicroCentsToCpmRounding ensures that we are converting from MicroCents to CPM correctly, considering rounding error may occur if incorrectly calculated.
func TestMicroCentsToCpmRounding(t *testing.T) {
	microcents := int64(32100000)
	expected := 0.321
	actual := MicroCentsToCpm(microcents)

	if expected != actual {
		t.Fail()
	}
}

// TestMicroCentsPerImpression ensures that we are properly calculating the number of micro cents to subtract per impression
func TestMicroCentsPerImpression(t *testing.T) {
	microcents := int64(32000)
	expected := int64(32)
	actual := MicroCentsPerImpression(microcents)

	if expected != actual {
		t.Fail()
	}
}
