package decoder

import "fmt"

const (
	PRIVATE_USE_AREA_A_START = 0xF0000
	PRIVATE_USE_AREA_B_START = 0x100000
)

type EncodedByteBuffer struct {
	Bytes []int
}

func NewEncodedByteBuffer() *EncodedByteBuffer {
	return &EncodedByteBuffer{
		Bytes: []int{},
	}
}

func (eb *EncodedByteBuffer) Add(b int) {
	eb.Bytes = append(eb.Bytes, b)
}

func FromUtf16String(s string) *EncodedByteBuffer {
	encodedBuffer := NewEncodedByteBuffer()
	codePoints := []int{}

	// Converte a string em code points
	for _, r := range s {
		codePoints = append(codePoints, int(r))
	}

	for _, codePoint := range codePoints {
		// Special cases
		if codePoint >= PRIVATE_USE_AREA_B_START {
			// Single byte
			singleByteOffset := PRIVATE_USE_AREA_B_START + 0xEE
			if (codePoint & 0xFF) == 0xEE {
				actualValue := (int(codePoint) - singleByteOffset)
				actualValue = actualValue >> 8
				encodedBuffer.Add(int(actualValue))

				if actualValue > 255 {
					panic(fmt.Sprintf("Invalid code point: %d", codePoint))
				}
				continue
			}

			// Two bytes
			values := codePoint - PRIVATE_USE_AREA_B_START

			encodedBuffer.Add(255)
			encodedBuffer.Add(int(254 + (values & 0xFF)))

			if codePoint >= 0x100002 {
				panic(fmt.Sprintf("Invalid code point: %d", codePoint))
			}
			continue
		}

		// Normal case
		values := codePoint - PRIVATE_USE_AREA_A_START

		encodedBuffer.Add(int(values >> 8))
		encodedBuffer.Add(int(values & 0xFF))

		if codePoint >= 0xFFFFE {
			panic(fmt.Sprintf("Invalid code point: %d", codePoint))
		}
	}

	return encodedBuffer
}
