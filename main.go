package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"strings"
	"encoding/hex"
	"encoding/binary"
	"math"

	"github.com/tarm/serial"
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





// SerialRead reads data from the serial port
func SerialRead(portName string) {
	// Configure the serial port
	config := &serial.Config{Name: portName, Baud: 9600, ReadTimeout: time.Second * 5}
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatalf("Failed to open port %s: %v\n", portName, err)
		return
	}
	defer port.Close()


	// Write data to the serial port - request signal data sending in continouos mode
	/*
	data_to_write := []byte("START\n")
	n_w, err_w := port.Write(data_to_write)
	if err_w != nil {
		log.Fatal(err_w)
	}	
	log.Printf("Sent %d bytes\n", n_w)
	*/


	// Write data to the serial port - request signal data through SEND function
	/*
	data_to_write := []byte("SEND\n")
	n_w, err_w := port.Write(data_to_write)
	if err_w != nil {
	        log.Fatal(err_w)
	    }	
	log.Printf("Sent %d bytes= \n", n_w)		
	*/


	// Buffer to store incoming data
	buf := make([]byte, 10240)

	fmt.Printf("Reading data from %s...\n", portName)

	for {
	    
	    // Write data to the serial port - request signal data through SEND function
	    
	    data_to_write := []byte("SEND\n")
	    n_w, err_w := port.Write(data_to_write)
	    if err_w != nil {
	        log.Fatal(err_w)
	    }	
	    log.Printf("Sent %d bytes\n", n_w)		


		// Read from the serial port
		n, err := port.Read(buf)
		if err != nil {
			log.Printf("Error reading from port %s: %v\n", portName, err)
			return
		}

		// Output the read data

		//fmt.Printf("Data from device= %s; is receved, data structure is=%s ; with size=%d \n", portName, string(buf[:n]),n)
		fmt.Printf("Data from device= %s; is receved, payload size=%d \n", portName,n)



		// Parse the received buffer 
		// Parsing frame_start and getting positions for all packets
		// 2 bytes

		frame_start := "AABB"

		positions, err := findAllHexPatternPositions(buf, len(buf), frame_start)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		return
		}

		var index int = 0
		// Output the result
		if len(positions) > 0 {

			for index < len(positions){
				fmt.Printf("Hex pattern '%s' found at positions: %v\n", frame_start, positions)

				// LENGTH
				// Parsing payload_length 
				poz1:= int(positions[index])+2
				poz2:= poz1 + 2
				
				fmt.Printf("poz1 for length '%d' \n", poz1)
				fmt.Printf("poz1 for length'%d' \n", poz2)

				payload_length := [2]byte(buf[poz1:poz2])
				fmt.Printf("payload length= '%X' \n", payload_length)	


				// TIMESTAMP
				// Parsing payload_timestamp
				
				poz1= int(positions[index])+4
				poz2= poz1+4

				fmt.Printf("poz1 for timestamp '%d' \n", poz1)
				fmt.Printf("poz1 for timestamp'%d' \n", poz2)


				payload_timestamp := []byte(buf[poz1:poz2])
				fmt.Printf("payload timestamp= '%X' \n", payload_timestamp)	

				// TIMESTAMP CONVERSION				

				// Step 2: Convert the byte slice to a uint32 value (Little Endian)
				if len(payload_timestamp) != 4 {
					fmt.Println("Invalid length for a 32-bit timestamp representation")
					return
				}

				// Interpret the bytes as a 32-bit unsigned integer
				timestamp := binary.BigEndian.Uint32(payload_timestamp)
				//timestamp_int := int32(timestamp)

				// Step 3: Convert the Unix timestamp to a time.Time object
				// Unix expects the number of seconds since 1970-01-01
				//timeValue := time.Unix(int64(timestamp), 0)

				// Print the converted timestamp
				fmt.Printf("VALUE 1  %v , VALUE 2 %T \n", timestamp, timestamp)
				//fmt.Printf("Readable timestamp: %s\n", timeValue.Format(time.RFC3339))
				
			

				// PAYLOAD PIR 1 STATUS
				// Parsing payload_pir_status
				poz1= int(positions[index])+8
				poz2= poz1+1

				fmt.Printf("poz1 for pir status '%d' \n", poz1)
				fmt.Printf("poz2 for pir status '%d' \n", poz2)

				payload_pir_status := []byte(buf[poz1:poz2])
				fmt.Printf("payload pir status= '%X' \n", payload_pir_status)	


				// PAYLOAD PIR1 VALUE
				// Parsing payload_pir_1_value
				poz1= int(positions[index])+9
				poz2= poz1+4

				fmt.Printf("poz1 for PIR 1 value '%d' \n", poz1)
				fmt.Printf("poz2 for PIR 1 value '%d' \n", poz2)

				
				payload_pir1_value := []byte(buf[poz1:poz2])
				swapEndianess(payload_pir1_value)


				fmt.Printf("payload PIR 1 value= '%X' \n", payload_pir1_value)	

				payloat_pir1_value_fl, err := hexToFloat32BigEndian(payload_pir1_value)
				if err != nil {
					fmt.Printf("Error converting hex to float32: %v\n", err)
				} else {
					fmt.Printf("PIR1 float val = %f \n", payloat_pir1_value_fl)
				}


				// PAYLOAD PIR2 VALUE
				// Parsing payload_pir_2_value
				poz3:= int(positions[index])+13
				poz4:= poz3+4

				fmt.Printf("poz1 for PIR 2 value '%d' \n", poz3)
				fmt.Printf("poz2 for PIR 2 value '%d' \n", poz4)

				
				payload_pir2_value := []byte(buf[poz3:poz4])
				swapEndianess(payload_pir2_value)

				fmt.Printf("payload PIR 2 value= '%X' \n", payload_pir2_value)	

				payloat_pir2_value_fl, err := hexToFloat32BigEndian(payload_pir2_value)
				if err != nil {
					fmt.Printf("Error converting hex to float32: %v\n", err)
				} else {
					fmt.Printf("PIR2 float val = %f \n", payloat_pir2_value_fl)
				}

				// PAYLOAD CHECKSUM
				// Parsing checksum_value

				poz_xor_1:=int(positions[index])+17
				poz_xor_2:= poz_xor_1 + 1

				fmt.Printf("poz xor 1 for XOR value '%d' \n", poz_xor_1)
				fmt.Printf("poz xor 2 for XOR value '%d' \n", poz_xor_2)
	
				payload_xor_checksum  := [1]byte(buf[poz_xor_1:poz_xor_2])
				fmt.Printf("cheksum= '%X' \n", payload_xor_checksum )	




				time.Sleep(2 * time.Second) // separate points by 1 second

			index++
			}	

		} else {
			fmt.Printf("Hex pattern '%s' not found in the buffer.\n", frame_start)
		}

		/*
		
		*/



		// timer to sleep for 2 seconds before we reloop
		time.Sleep(2 * time.Second) // separate points by 1 second
	}
}



// MonitorSerialDevices continuously checks for new serial devices and starts a reader goroutine for each
func MonitorSerialDevices() {
	knownPorts := make(map[string]bool)

	for {
		// Get the list of serial ports
		entries, err := os.ReadDir("/dev")
		if err != nil {
			log.Fatalf("Failed to read /dev directory: %v\n", err)
		}

		// Look for new serial ports (usually tty.* on macOS)
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			if strings.HasPrefix(entry.Name(), "tty.usb") { // Adjust this as per your OS naming conventions
				portName := "/dev/" + entry.Name()

				// Start a new goroutine for each new port found
				if !knownPorts[portName] {
					knownPorts[portName] = true
					go SerialRead(portName) // Start reading from the new port in a separate goroutine
				}
			}
		}

		time.Sleep(5 * time.Second) // Check for new devices every 5 seconds
	}
}






func main() {
	fmt.Println("Starting to monitor serial connections...")

	// Start monitoring serial connections
	MonitorSerialDevices()



	/*
		entries, err := os.ReadDir("/dev")
	if err != nil {
		log.Fatalf("Failed to read /dev directory: %v\n", err)
	}

	// Iterate over each entry in the directory
	for _, entry := range entries {
		// Check if the entry's name starts with "tty."
		if len(entry.Name()) >= 4 && entry.Name()[:4] == "tty." {
			fmt.Printf("Found serial device: /dev/%s\n", entry.Name())
		}
		
	}
*/

}
