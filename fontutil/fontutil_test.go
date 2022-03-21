package fontutil

import (
	"strings"
	"testing"
)

func TestCheckIfErrorWhenContentTooWide(t *testing.T) {
	bigTextContent := "sdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdsdssd"
	_, _, err := PrintCenterWhiteTextBlackImage(12.0, 264, 176, bigTextContent, true, false)

	if err == nil || err != ErrTooBigForScreen {
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
						little lamb, little lamb, little lamb, little sheep, 
						baba lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb,`
	_, spilledText, err := PrintCenterWhiteTextBlackImage(12.0, 264, 176, bigTextContent, true, true)

	if len(spilledText) == 0 {
		t.Fatalf("SpilledText is empty when it's not supposed to")
	}

	wants := strings.Fields(`baba lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb, 
						little lamb, little lamb, little lamb, little lamb,`)

	if !stringSlicesAreEqual(wants, spilledText){
		t.Fatalf("SpilledText is unexpected")
	}

	if err == nil || err != ErrContinueNextScreen{
		t.Fatalf("Content spilled, but err is not continue next screen \n")
	}

	//t.Logf("spilledText %q \n", spilledText)
}

func TestCheckIfErrorWhenSingleContentTooTall(t *testing.T) {
	bigTextContent :=  "Mary has a little lamb"
	_, _, err := PrintCenterWhiteTextBlackImage(200.0, 264, 176, bigTextContent, true, true)

	if err == nil || err != ErrTooBigForScreen {
		t.Fatalf("Content too Tall, but not detected \n")
	}
}

func stringSlicesAreEqual(sa, sb []string) bool{
	if len(sa) != len(sb){
		return false
	}

	idx := 0
	for _, content := range(sa){

		if sb[idx] != content {
			return false
		}
		idx += 1
	}

	return true

}