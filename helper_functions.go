package main

import (
	"fmt"
	"time"
	"encoding/hex"
	"encoding/binary"
	"math"
)






func hexToFloat32BigEndian(bytes []byte) (float32, error) {
	// Decode the hex string into bytes
	//bytes, err := hex.DecodeString(hexStr)
	//if err != nil {
	//	return 0, fmt.Errorf("invalid hex string: %v", err)
	//}

	// Check if the length of the bytes is exactly 4 for a float32 representation
	if len(bytes) != 4 {
		return 0, fmt.Errorf("hex string must represent exactly 4 bytes for float32")
	}

	fmt.Printf("bytes are %v \n", bytes)


	// Convert bytes to a uint32 using Big Endian byte order
	bits := binary.BigEndian.Uint32(bytes)

	// Convert the uint32 bits to a float32
	floatVal := math.Float32frombits(bits)

	return floatVal, nil
}


func swapEndianess(bytes []byte) {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
}




// findAllHexPatternPositions returns a slice of positions where the specific hex pattern is found in the buffer.
// If the pattern is not found, it returns an empty slice.
func findAllHexPatternPositions(buf []byte, n int, hexPattern string) ([]int, error) {
	// Decode the hex pattern into a byte slice
	searchBytes, err := hex.DecodeString(hexPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid hex pattern: %s", hexPattern)
	}

	searchLen := len(searchBytes)

	// Check if the buffer is smaller than the pattern length
	if n < searchLen {
		return nil, fmt.Errorf("pattern length (%d) is greater than buffer size (%d)", searchLen, n)
	}

	// Store all positions of the found pattern
	var positions []int

	// Loop through the buffer to find all occurrences of the pattern
	for i := 0; i <= n-searchLen; i++ {
		// Compare the bytes in the buffer to the search pattern
		if matchBytes(buf[i:i+searchLen], searchBytes) {
			positions = append(positions, i) // Store the starting position of the pattern
		}
	}

	return positions, nil
}




// matchBytes compares two byte slices for equality
func matchBytes(buf1, buf2 []byte) bool {
	for i := range buf1 {
		if buf1[i] != buf2[i] {
			return false
		}
	}
	return true
}






// 
// XOR checksum function 
func xorChecksum(data []byte) byte {
	var checksum byte = 0
	for _, b := range data {

		checksum ^= b
		fmt.Printf("checksum: %X \n", checksum)
		
		time.Sleep(1 * time.Second) // separate points by 1 second
	}
	return checksum
}





