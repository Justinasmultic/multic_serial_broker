package main

import (
	"fmt"
	"math/rand"
	"encoding/binary"
	"time"
	"os"
	"log"
	"context"
	"strings"
	"errors"
	"github.com/tarm/serial"
	//"go.bug.st/serial.v1"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

)



// MonitorSerialDevices continuously checks for new serial devices and starts a reader goroutine for each
// 
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
					//go SerialReadSinglePayload(portName) // Start reading from the new port in a separate goroutine, read single payload
					go SerialReadContinuousPayload(portName) // Start reading from the new port in a separate goroutine, read single payload
				}
			}
		}

		time.Sleep(5 * time.Second) // Check for new devices every 5 seconds
	}
}




func unrecoverableError(err error) bool {
    // Define conditions for an unrecoverable error
    return errors.Is(err, os.ErrClosed) 
}





// SerialReadContinuousPayload reads data from the serial port
// READ continuous payload stream by initiating it through START command

func SerialReadContinuousPayload(portName string) {
	



	// Configure the serial port
	config := &serial.Config{Name: portName, Baud: 115200}
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatalf("Failed to open port %s: %v\n", portName, err)
		return
	}
	//defer port.Close()


	// Write data to the serial port - request signal data sending in continouos mode
	
	data_to_write := []byte("START\n")
	n_w, err_w := port.Write(data_to_write)
	if err_w != nil {
		log.Fatal(err_w)
	}	
	log.Printf("Sent %d bytes\n", n_w)

	//defer port.Close()

	time.Sleep(1 * time.Second)


	fifo := make(chan []byte, 10240) 

	
// GOROUTINE to read continuos stream of data every 1ms
	go func (portforroutine *serial.Port) {
	for {
		// Read from the serial port
		buf := make([]byte, 10240)
		n, err := portforroutine.Read(buf)
		if err != nil {
			log.Printf("Error reading from port %s: %v\n", portName, err)
			return
		}

        if n > 0 {
            //fmt.Printf("Writing %d bytes to FIFO\n", n)
            fifo <- buf[:n]
        }


		// Output the read data

		
		//fmt.Printf("Data from device= %s; is receved, payload size=%d \n", portName,n)
		//fmt.Printf("\nDATA RECEIVED TIMESTAMP %s \n", time.Now())

		//time.Sleep(10 * time.Millisecond) // separate points by 1 second
		}
	} (port)



// GOROUTINE to read from FIFO buffer and print how much has been read
    go func() {
        for data := range fifo {
            fmt.Printf("Read from FIFO %d bytes: \n", len(data))
            //fmt.Printf("FIFO size %d bytes: \n", len(fifo))
            //fmt.Printf("FIFO message : %b \n", data)


	        positions, err := findPos(data)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			return
			}

			go parseAndStore(portName, positions, data)

			//fmt.Printf("\nAll positions are: %d  \n", positions)

        }
        //time.Sleep(10 * time.Millisecond) // separate points by 1 second
    }()



    // Goroutine 
    /*
    go func() {

		// Parse the received buffer 
		// Parsing frame_start and getting positions for all packets
		// 2 bytes

		//frame_start := "AABB"

		positions, err := findPos(data)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		return
		}	

		fmt.Printf("\nAll positions are: %d  \n", positions)

		time.Sleep(1 * time.Second) // separate points by 1 second
    }()
	*/

		// Output the read data

		//fmt.Printf("Data from device= %s; is receved, data structure is=%s ; with size=%d \n", portName, string(buf[:n]),n)
		//fmt.Printf("Data from device= %s; is receved, payload size=%d \n", portName,n)
		//fmt.Printf("\nDATA RECEIVED TIMESTAMP %s \n", time.Now())

		//time.Sleep(10 * time.Millisecond) // separate points by 1 second


/*		

		// Parse the received buffer 
		// Parsing frame_start and getting positions for all packets
		// 2 bytes

		frame_start := "AABB"

		positions, err := findAllHexPatternPositions(buf, len(buf), frame_start)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		return
		}

		

		// init current timestamp for this cycle
		
		//this_cycle_time := time.Now() 
		//timenow := time.Now() 
		//this_cycle_time := timenow.Format("2006-01-02 15:04:05.000000")
		//timenow.UnixMicro()

		//fmt.Printf("CYCLE TIMESTAMP IS = %s \n", this_cycle_time)
		//time.Sleep(3 * time.Second) // separate points by 1 second


		var index int = 0
		// Output the result
		if len(positions) > 0 {

			for index < len(positions){
				
				//fmt.Printf("\nHex pattern '%s' found at positions: %v\n", frame_start, positions)

				// LENGTH
				// Parsing payload_length 
				//poz_length:= int(positions[index])+2
				//poz_length2:= poz_length + 2
				
				//fmt.Printf("poz1 for length '%d' \n", poz_length)
				//fmt.Printf("poz1 for length'%d' \n", poz_length2)

				//payload_length := [2]byte(buf[poz_length:poz_length2])
				//fmt.Printf("payload length= '%X' \n", payload_length)	


				// TIMESTAMP
				// Parsing payload_timestamp
				
				poz1:= int(positions[index])+4
				poz2:= poz1+4

				//fmt.Printf("poz1 for timestamp '%d' \n", poz1)
				//fmt.Printf("poz1 for timestamp'%d' \n", poz2)


				payload_timestamp := []byte(buf[poz1:poz2])
				//fmt.Printf("payload timestamp= '%X' \n", payload_timestamp)	

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
				//fmt.Printf("VALUE 1  %v , VALUE 2 %T \n", timestamp, timestamp)
				//fmt.Printf("Readable timestamp: %s\n", timeValue.Format(time.RFC3339))
				
			

				// PAYLOAD PIR 1 STATUS
				// Parsing payload_pir_status
				poz1= int(positions[index])+8
				poz2= poz1+1

				//fmt.Printf("poz1 for pir status '%d' \n", poz1)
				//fmt.Printf("poz2 for pir status '%d' \n", poz2)

				payload_pir_status := []byte(buf[poz1:poz2])
				//fmt.Printf("\n PIR status= '%x' \n", payload_pir_status)	


				// PAYLOAD PIR1 VALUE
				// Parsing payload_pir_1_value
				poz1= int(positions[index])+9
				poz2= poz1+4

				//fmt.Printf("poz1 for PIR 1 value '%d' \n", poz1)
				//fmt.Printf("poz2 for PIR 1 value '%d' \n", poz2)

				
				payload_pir1_value := []byte(buf[poz1:poz2])
				swapEndianess(payload_pir1_value)


				//fmt.Printf("payload PIR 1 value= '%X' \n", payload_pir1_value)	

				payload_pir1_value_fl, err := hexToFloat32BigEndian(payload_pir1_value)

				if err != nil {
					fmt.Printf("Error converting hex to float32: %v\n", err)
				} else {
					//fmt.Printf("PIR1 float val = %f \n", payload_pir1_value_fl)
				}


				// PAYLOAD PIR2 VALUE
				// Parsing payload_pir_2_value
				poz3:= int(positions[index])+13
				poz4:= poz3+4

				//fmt.Printf("poz1 for PIR 2 value '%d' \n", poz3)
				//fmt.Printf("poz2 for PIR 2 value '%d' \n", poz4)

				
				payload_pir2_value := []byte(buf[poz3:poz4])
				swapEndianess(payload_pir2_value)

				//fmt.Printf("payload PIR 2 value= '%X' \n", payload_pir2_value)	

				payloat_pir2_value_fl, err := hexToFloat32BigEndian(payload_pir2_value)
				if err != nil {
					fmt.Printf("Error converting hex to float32: %v\n", err)
				} else {
					//fmt.Printf("PIR2 float val = %f \n", payloat_pir2_value_fl)
				}

				// PAYLOAD CHECKSUM
				// Parsing checksum_value

				//poz_xor_1:=int(positions[index])+17
				//poz_xor_2:= poz_xor_1 + 1

				//fmt.Printf("poz xor 1 for XOR value '%d' \n", poz_xor_1)
				//fmt.Printf("poz xor 2 for XOR value '%d' \n", poz_xor_2)
	
				//payload_xor_checksum  := [1]byte(buf[poz_xor_1:poz_xor_2])
				//fmt.Printf("cheksum from payload = '%X' \n", payload_xor_checksum )	


				// do XOR for this range 
				//calculated_xor_checksum := xorChecksum(buf[int(positions[index]):poz_xor_1])
				//fmt.Printf("calculated cheksum = '%X' \n", calculated_xor_checksum )

				//fmt.Printf("full payload = '%X' \n", string (buf[poz_length:poz_xor_2]))


				

				payload_pir1_value_fl = rand.Float32()*(5-200)
				payloat_pir2_value_fl = rand.Float32()*(10-100)

				this_cycle_time := time.Now() 

				//fmt.Printf("\nTIMESTAMP WRITTEN IS = %s \n", this_cycle_time)

				write_influx_datapoint (portName, byteHexToInt(payload_pir_status), this_cycle_time , timestamp, payload_pir1_value_fl, payloat_pir2_value_fl)


				//time.Sleep(2 * time.Second) // separate points by 1 second

			index++
			}	

		} else {
			fmt.Printf("Hex pattern '%s' not found in the buffer.\n", frame_start)
		}

		// timer to sleep for 2 seconds before we reloop
		//time.Sleep(1 * time.Second) // separate points by 1 second





*/		

}







func parseAndStore (portName string, positions []int, buf []byte) {

		var index int = 0
		// Output the result
		if len(positions) > 0 {

			for index < len(positions) {
				
				//fmt.Printf("\nHex pattern '%s' found at positions: %v\n", frame_start, positions)

				// LENGTH
				// Parsing payload_length 
				//poz_length:= int(positions[index])+2
				//poz_length2:= poz_length + 2
				
				//fmt.Printf("poz1 for length '%d' \n", poz_length)
				//fmt.Printf("poz1 for length'%d' \n", poz_length2)

				//payload_length := [2]byte(buf[poz_length:poz_length2])
				//fmt.Printf("payload length= '%X' \n", payload_length)	


				// TIMESTAMP
				// Parsing payload_timestamp
				
				poz1:= int(positions[index])+4
				poz2:= poz1+4

				//fmt.Printf("poz1 for timestamp '%d' \n", poz1)
				//fmt.Printf("poz1 for timestamp'%d' \n", poz2)


				payload_timestamp := []byte(buf[poz1:poz2])
				//fmt.Printf("payload timestamp= '%X' \n", payload_timestamp)	

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
				//fmt.Printf("VALUE 1  %v , VALUE 2 %T \n", timestamp, timestamp)
				//fmt.Printf("Readable timestamp: %s\n", timeValue.Format(time.RFC3339))
				
			

				// PAYLOAD PIR 1 STATUS
				// Parsing payload_pir_status
				poz1= int(positions[index])+8
				poz2= poz1+1

				//fmt.Printf("poz1 for pir status '%d' \n", poz1)
				//fmt.Printf("poz2 for pir status '%d' \n", poz2)

				payload_pir_status := []byte(buf[poz1:poz2])
				//fmt.Printf("\n PIR status= '%x' \n", payload_pir_status)	


				// PAYLOAD PIR1 VALUE
				// Parsing payload_pir_1_value
				poz1= int(positions[index])+9
				poz2= poz1+4

				//fmt.Printf("poz1 for PIR 1 value '%d' \n", poz1)
				//fmt.Printf("poz2 for PIR 1 value '%d' \n", poz2)

				
				payload_pir1_value := []byte(buf[poz1:poz2])
				swapEndianess(payload_pir1_value)


				//fmt.Printf("payload PIR 1 value= '%X' \n", payload_pir1_value)	

				payload_pir1_value_fl, err := hexToFloat32BigEndian(payload_pir1_value)

				if err != nil {
					fmt.Printf("Error converting hex to float32: %v\n", err)
				} else {
					//fmt.Printf("PIR1 float val = %f \n", payload_pir1_value_fl)
				}


				// PAYLOAD PIR2 VALUE
				// Parsing payload_pir_2_value
				poz3:= int(positions[index])+13
				poz4:= poz3+4

				//fmt.Printf("poz1 for PIR 2 value '%d' \n", poz3)
				//fmt.Printf("poz2 for PIR 2 value '%d' \n", poz4)

				
				payload_pir2_value := []byte(buf[poz3:poz4])
				swapEndianess(payload_pir2_value)

				//fmt.Printf("payload PIR 2 value= '%X' \n", payload_pir2_value)	

				payloat_pir2_value_fl, err := hexToFloat32BigEndian(payload_pir2_value)
				if err != nil {
					fmt.Printf("Error converting hex to float32: %v\n", err)
				} else {
					//fmt.Printf("PIR2 float val = %f \n", payloat_pir2_value_fl)
				}

				// PAYLOAD CHECKSUM
				// Parsing checksum_value

				//poz_xor_1:=int(positions[index])+17
				//poz_xor_2:= poz_xor_1 + 1

				//fmt.Printf("poz xor 1 for XOR value '%d' \n", poz_xor_1)
				//fmt.Printf("poz xor 2 for XOR value '%d' \n", poz_xor_2)
	
				//payload_xor_checksum  := [1]byte(buf[poz_xor_1:poz_xor_2])
				//fmt.Printf("cheksum from payload = '%X' \n", payload_xor_checksum )	


				// do XOR for this range 
				//calculated_xor_checksum := xorChecksum(buf[int(positions[index]):poz_xor_1])
				//fmt.Printf("calculated cheksum = '%X' \n", calculated_xor_checksum )

				//fmt.Printf("full payload = '%X' \n", string (buf[poz_length:poz_xor_2]))
				


				payload_pir1_value_fl = rand.Float32()*(5-200)
				payloat_pir2_value_fl = rand.Float32()*(10-100)

				this_cycle_time := time.Now() 

				//fmt.Printf("\nTIMESTAMP WRITTEN IS = %s \n", this_cycle_time)

				write_influx_datapoint (portName, byteHexToInt(payload_pir_status), this_cycle_time , timestamp, payload_pir1_value_fl, payloat_pir2_value_fl)


				//time.Sleep(2 * time.Second) // separate points by 1 second

			index++

			}	

		}
}




func write_influx_datapoint (sensor_name string, sensor_ir_status int, this_cycle_time time.Time, timestamp uint32, pir1_value float32, pir2_value float32) {

	token := os.Getenv("INFLUX_TOKEN")

	//fmt.Printf("token= %s \n", token)

	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)
	org := "multic"
	bucket := "ikea_pilot1"
	writeAPI := client.WriteAPIBlocking(org, bucket)
	

	tags := map[string]string{
		"sensor_name": sensor_name,
		"location": "Klaipedos Baldai",
		"recipe": "Stalas sokiams ",
	}


	fields := map[string]interface{}{
		"pir1_value": pir1_value,
		"pir2_value": pir2_value,
		"sensor_ir_status": sensor_ir_status,
		"sensor_timestamp": timestamp,
	}

	//fmt.Printf("writing sensor name %s pir1 %f and pir2 %f and sensor ir status %d and TIMESTAMP is %s  \n", sensor_name, pir1_value, pir2_value, sensor_ir_status, this_cycle_time)


	point := write.NewPoint("sensor_pir_measurement", tags, fields, this_cycle_time)	


	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			log.Fatal(err)
	}

	//writeAPI.Flush()
	//time.Sleep(1 * time.Second) // separate points by 1 second


}
