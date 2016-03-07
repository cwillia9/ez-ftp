package authentication

import "testing"

func TestComputeHmac256(t *testing.T) {
	res := ComputeHmac256("message", "secret")

	expectation := "i19IcCmVwVmMVz2x4hhmqbgl1KeU0WnXBgoDYFeWNgs="

	if res != expectation {
		t.Error("Expected", expectation, "got", res)
	}

}

func TestComputeHmac1(t *testing.T) {
	res := ComputeHmac1("message", "secret")

	expectation := "DK9kn+7klT2Hv5A6wRdsReAo3xY="

	if res != expectation {
		t.Error("Expected", expectation, "got", res)
	}

}
