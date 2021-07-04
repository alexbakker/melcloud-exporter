package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	config     *Config
	flagConfig = flag.String("config", "config.yml", "the filename of the configuration file")

	client = NewClient()

	deviceLabelNames = []string{"building_id", "device_id", "device_name"}
	gaugeDevicePower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "melcloud_device_power",
		Help: "Whether the device is powered on",
	}, deviceLabelNames)
	gaugeDeviceTemperatureRoom = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "melcloud_device_temperature_room",
		Help: "The current temperature in the room a device is in",
	}, deviceLabelNames)
	gaugeDeviceTemperatureSet = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "melcloud_device_temperature_set",
		Help: "The temperature that the device targets to achieve in the room it is in",
	}, deviceLabelNames)
)

func main() {
	flag.Parse()

	var err error
	config, err = ReadConfig(*flagConfig)
	if err != nil {
		fmt.Printf("unable to read config: %s\n", err)
		os.Exit(1)
		return
	}

	if err = client.Login(config.MELCloud.Email, config.MELCloud.Password); err != nil {
		fmt.Printf("unable to log into melcloud: %s\n", err)
		os.Exit(1)
		return
	}

	if err = updateData(); err != nil {
		fmt.Printf("unable to update data: %s\n", err)
		os.Exit(1)
		return
	}

	go func() {
		tick := time.Tick(time.Duration(config.MELCloud.RefreshInterval) * time.Second)
		for {
			<-tick
			if err := updateData(); err != nil {
				fmt.Printf("unable to update data: %s\n", err)
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(config.Prometheus.Addr, nil)
}

func updateData() error {
	fmt.Println("updating data")

	devices, err := client.Devices()
	if err != nil {
		return err
	}

	for _, dev := range devices {
		labels := prometheus.Labels{
			"building_id": strconv.Itoa(dev.BuildingID),
			"device_id":   strconv.Itoa(dev.DeviceID),
			"device_name": dev.DeviceName,
		}

		power := 0
		if dev.Device.Power {
			power = 1
		}

		gaugeDevicePower.With(labels).Set(float64(power))
		gaugeDeviceTemperatureRoom.With(labels).Set(float64(dev.Device.RoomTemperature))
		gaugeDeviceTemperatureSet.With(labels).Set(float64(dev.Device.SetTemperature))
	}

	return nil
}
