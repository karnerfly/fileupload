package middlewares

import "testing"

func TestIsValidImageId(t *testing.T) {
	r := isValidImageId("4293")
	if !r {
		t.Fail()
	}
}

func TestIsValidFileName(t *testing.T) {
	r := isValidFileName("demo_imgfl0_12fa3.mp4")
	if !r {
		t.Fail()
	}
}
