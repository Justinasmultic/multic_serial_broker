package main

import (
	"context"
	"log"
	"os"
	"time"
	"math/rand"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)



func init_influx_connection () {


}

func write_influx_datapoint () {

	token := os.Getenv("INFLUXDB_TOKEN")
	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)
	org := "multic"
	bucket := "ikea_pilot"
	writeAPI := client.WriteAPIBlocking(org, bucket)


	randomized_val:=rand.Intn(100)
	randomized_val2:=rand.Intn(100)	

	tags := map[string]string{
		"sensor": "sensor-1",
		"location": "Klaipedos Baldai",
		"recipe": "Juodas alamo 13",
	}


	fields := map[string]interface{}{
		"pir1_value": randomized_val,
		"pir2_value": randomized_val2,
	}

	point := write.NewPoint("sensor", tags, fields, time.Now())	


	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			log.Fatal(err)
	}


}