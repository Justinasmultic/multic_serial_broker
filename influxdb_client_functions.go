package main

import (
	"context"
	"log"
	"os"
	"time"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)



func init_influx_connection () {


}

func write_influx_datapoint (sensor_name string, sensor_ir_status []byte, pir1_value float32, pir2_value float32) {

	token := os.Getenv("INFLUX_TOKEN")

	fmt.Printf("token= %s \n", token)

	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)
	org := "multic"
	bucket := "ikea_pilot"
	writeAPI := client.WriteAPIBlocking(org, bucket)


	//randomized_val:=rand.Intn(100)
	//randomized_val2:=rand.Intn(100)	

	tags := map[string]string{
		"sensor_name": sensor_name,
		"location": "Klaipedos Baldai",
		"recipe": "Stalas sokiams 666",
	}



	fields := map[string]interface{}{
		"pir1_value": pir1_value,
		"pir2_value": pir2_value,
		"sensor_ir_status": sensor_ir_status,
	}

	fmt.Printf("writing sensor name %s pir1 %f and pir2 %f and sensor ir status %X  \n", sensor_name, pir1_value, pir2_value, sensor_ir_status)


	point := write.NewPoint("sensor_pir_measurement", tags, fields, time.Now())	


	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			log.Fatal(err)
	}


}