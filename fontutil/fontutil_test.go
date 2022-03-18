package fontutil

import (
	"testing"
)

func TestCheckIfErrorWhenContentTooWide(t *testing.T) {
	bigTextContent := "sdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdssd"
	_, err := PrintCenterWhiteTextBlackImage(12.0, 264, 176, bigTextContent, true, false)

	if err == nil || err != ErrToBigForScreen {
		t.Fatalf("Content too wide, but not detected \n")
	}
}

func TestCheckIfErrorWhenContentTooTall(t *testing.T) {
	bigTextContent :=  `Mary has a little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb,
						little lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb,`
	_, err := PrintCenterWhiteTextBlackImage(12.0, 264, 176, bigTextContent, true, true)

	if err == nil || err != ErrToBigForScreen {
		t.Fatalf("Content too Tall, but not detected \n")
	}
}

func TestCheckIfErrorWhenSingleContentTooTall(t *testing.T) {
	bigTextContent :=  "Mary has a little lamb"
	_, err := PrintCenterWhiteTextBlackImage(100.0, 264, 176, bigTextContent, true, true)

	if err == nil || err != ErrToBigForScreen {
		t.Fatalf("Content too Tall, but not detected \n")
	}
}