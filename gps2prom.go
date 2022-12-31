package main

import (
	"fmt"
	"github.com/stratoberry/go-gpsd"

	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	gpsMode = promauto.NewGauge(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"fix_mode",
			Help:		"gpsd mode (0=NoValueSeen, 1=NoFix, 2=2D fix, 3=3D fix)",
		})

	gpsNumSats = promauto.NewGauge(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"space_vehicles_total",
			Help:		"Total number of SV",
		})

	gpsNumSatsUsed = promauto.NewGauge(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"space_vehicles_used",
			Help:		"Used number of SV",
		})

	gpsSvSs = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"space_vehicle_signal_strength",
			Help:		"Per-SV signal strength",
		},
		[]string{
			"prn",
		},
	)

	gpsXdop = promauto.NewGauge(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"xdop",
		})
	gpsYdop = promauto.NewGauge(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"ydop",
		})
	gpsVdop = promauto.NewGauge(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"vdop",
		})
	gpsTdop = promauto.NewGauge(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"tdop",
		})
	gpsHdop = promauto.NewGauge(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"hdop",
		})
	gpsPdop = promauto.NewGauge(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"pdop",
		})
	gpsGdop = promauto.NewGauge(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"gdop",
		})
)

func main() {
	var gps *gpsd.Session
	var err error

	if gps, err = gpsd.Dial("apu-left:2947"); err != nil {
		panic(fmt.Sprintf("Failed to connect to GPSD: %s", err))
	}

	gps.AddFilter("TPV", func(r interface{}) {
		tpv := r.(*gpsd.TPVReport)
		//fmt.Println("TPV", tpv.Mode, tpv.Time)
		gpsMode.Set(float64(tpv.Mode))
	})

	/*
	gps.AddFilter("DEVICE", func(r interface{}) {
		dev := r.(*gpsd.DEVICEReport)
		fmt.Println("DEVICE", dev.Path, dev.Flags)
	})
	*/

	gps.AddFilter("SKY", func(r interface{}) {
		sky := r.(*gpsd.SKYReport)
		fmt.Println("SKY", sky.Satellites)
		gpsXdop.Set(sky.Xdop)
		gpsYdop.Set(sky.Ydop)
		gpsVdop.Set(sky.Vdop)
		gpsTdop.Set(sky.Tdop)
		gpsHdop.Set(sky.Hdop)
		gpsPdop.Set(sky.Pdop)
		gpsGdop.Set(sky.Gdop)

		gpsNumSats.Set(float64(len(sky.Satellites)))
		gpsSvSs.Reset()
		num_sats_used := 0
		for i := 0; i < len(sky.Satellites); i++ {
			num_sats_used += 1
			prn_str := fmt.Sprintf("%.0f", sky.Satellites[i].PRN)
			gpsSvSs.WithLabelValues(prn_str).Set(sky.Satellites[i].Ss)
		}
		gpsNumSatsUsed.Set(float64(num_sats_used))
	})


	fmt.Println("Hello, World!")

	done := gps.Watch()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

	<-done
}
