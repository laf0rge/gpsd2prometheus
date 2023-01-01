// (C) 2022 by Harald Welte <laforge@gnumonks.org>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"flag"
	"github.com/stratoberry/go-gpsd"

	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	gpsd_tcp = flag.String("gpsd","localhost:2947", "remote gpsd host:port")
	listen_http = flag.String("http",":2112", "local HTTP listen host:port")

	gpsMode = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"fix_mode",
			Help:		"gpsd mode (0=NoValueSeen, 1=NoFix, 2=2D fix, 3=3D fix)",
		},
		[]string{
			"device",
		})

	gpsLat = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"latitude",
			Help:		"Latitude in degrees: +/- signifies North/South.",
		},
		[]string{
			"device",
		})
	gpsLon = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"longitude",
			Help:		"Longitude in degrees: +/- signifies East/West.",
		},
		[]string{
			"device",
		})
	gpsAlt = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"altitude",
		},
		[]string{
			"device",
		})

	gpsEpt = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"estimated_error_timestamp",
			Help:		"Estimated time stamp error in seconds. Certainty unknown.",
		},
		[]string{
			"device",
		})
	gpsEpx = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"estimated_error_longitude",
			Help:		"Longitude error estimate in meters. Certainty unknown.",
		},
		[]string{
			"device",
		})
	gpsEpy = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"estimated_error_latitude",
			Help:		"Latitude error estimate in meters. Certainty unknown.",
		},
		[]string{
			"device",
		})
	gpsEpv = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"estimated_error_altitude",
			Help:		"Estimated vertical error in meters. Certainty unknown.",
		},
		[]string{
			"device",
		})
	gpsTrack = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"track",
			Help:		"Course over ground, degrees from true north.",
		},
		[]string{
			"device",
		})
	gpsSpeed = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"speed",
			Help:		"Speed over ground, meters per second.",
		},
		[]string{
			"device",
		})
	gpsClimb = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"climb",
			Help:		"Climb (positive) or sink (negative) rate, meters per second.",
		},
		[]string{
			"device",
		})
	gpsEpd = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"estimated_error_direction",
			Help:		"Estimated track (direction) error in degrees. Certainty unknown.",
		},
		[]string{
			"device",
		})
	gpsEps = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"estimated_error_speed",
			Help:		"Estimated speed error in meters per second. Certainty unknown.",
		},
		[]string{
			"device",
		})
	gpsEpc = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name: 		"estimated_error_climb",
			Help:		"Estimated climb error in meters per second. Certainty unknown.",
		},
		[]string{
			"device",
		})





	gpsNumSats = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"space_vehicles_total",
			Help:		"Total number of space vehicles observed.",
		},
		[]string{
			"device",
		})
	gpsNumSatsUsed = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"space_vehicles_used",
			Help:		"Number of space vehicles used in fix.",
		},
		[]string{
			"device",
		})

	gpsSvAz = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"space_vehicle_azimuth",
			Help:		"Per-SV Azimuth, degrees from true north.",
		},
		[]string{
			"device",
			"prn",
		},
	)
	gpsSvEl = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"space_vehicle_elevation",
			Help:		"Per-SV Elevation in degrees.",
		},
		[]string{
			"device",
			"prn",
		},
	)
	gpsSvSs = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"space_vehicle_signal_noise_ratio",
			Help:		"Per-SV Signal to Noise ratio in dBHz.",
		},
		[]string{
			"device",
			"prn",
		},
	)

	gpsXdop = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"dilution_of_precision_longitude",
			Help:		"Longitudinal dilution of precision, a dimensionless factor which should be multiplied by a base UERE to get an error estimate.",
		},
		[]string{
			"device",
		},
	)
	gpsYdop = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"dilution_of_precision_latitude",
			Help:		"Latitudinal dilution of precision, a dimensionless factor which should be multiplied by a base UERE to get an error estimate.",
		},
		[]string{
			"device",
		},
	)
	gpsVdop = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"dilution_of_precision_altitude",
			Help:		"Vertical (altitude) dilution of precision, a dimensionless factor which should be multiplied by a base UERE to get an error estimate.",
		},
		[]string{
			"device",
		},
	)
	gpsTdop = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"dilution_of_precision_time",
			Help:		"Time dilution of precision, a dimensionless factor which should be multiplied by a base UERE to get an error estimate.",
		},
		[]string{
			"device",
		},
	)
	gpsHdop = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"dilution_of_precision_horizontal",
			Help:		"Horizontal dilution of precision, a dimensionless factor which should be multiplied by a base UERE to get a circular error estimate.",
		},
		[]string{
			"device",
		},
	)
	gpsPdop = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"dilution_of_precision_position",
			Help:		"Position (spherical/3D) dilution of precision, a dimensionless factor which should be multiplied by a base UERE to get an error estimate.",
		},
		[]string{
			"device",
		},
	)
	gpsGdop = promauto.NewGaugeVec(
		prometheus.GaugeOpts {
			Subsystem:	"gpsd",
			Name:		"dilution_of_precision_geometric",
			Help:		"Geometric (hyperspherical) dilution of precision, a combination of PDOP and TDOP. A dimensionless factor which should be multiplied by a base UERE to get an error estimate.",
		},
		[]string{
			"device",
		},
	)

)

func main() {
	var gps *gpsd.Session
	var err error

	flag.Parse()

	fmt.Println("Connecting to gpsd at", *gpsd_tcp, "...")
	if gps, err = gpsd.Dial(*gpsd_tcp); err != nil {
		panic(fmt.Sprintf("Failed to connect to GPSD: %s", err))
	}
	fmt.Println("Connected to gpsd!")

	gps.AddFilter("TPV", func(r interface{}) {
		tpv := r.(*gpsd.TPVReport)
		//fmt.Println("TPV", tpv.Mode, tpv.Time)
		gpsMode.WithLabelValues(tpv.Device).Set(float64(tpv.Mode))
		//time
		gpsEpt.WithLabelValues(tpv.Device).Set(tpv.Ept)
		gpsLat.WithLabelValues(tpv.Device).Set(tpv.Lat)
		gpsLon.WithLabelValues(tpv.Device).Set(tpv.Lon)
		gpsAlt.WithLabelValues(tpv.Device).Set(tpv.Alt)
		gpsEpx.WithLabelValues(tpv.Device).Set(tpv.Epx)
		gpsEpy.WithLabelValues(tpv.Device).Set(tpv.Epy)
		gpsEpv.WithLabelValues(tpv.Device).Set(tpv.Epv)
		gpsTrack.WithLabelValues(tpv.Device).Set(tpv.Track)
		gpsClimb.WithLabelValues(tpv.Device).Set(tpv.Climb)
		gpsEpd.WithLabelValues(tpv.Device).Set(tpv.Epd)
		gpsEps.WithLabelValues(tpv.Device).Set(tpv.Eps)
		gpsEpc.WithLabelValues(tpv.Device).Set(tpv.Epc)
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
		gpsXdop.WithLabelValues(sky.Device).Set(sky.Xdop)
		gpsYdop.WithLabelValues(sky.Device).Set(sky.Ydop)
		gpsVdop.WithLabelValues(sky.Device).Set(sky.Vdop)
		gpsTdop.WithLabelValues(sky.Device).Set(sky.Tdop)
		gpsHdop.WithLabelValues(sky.Device).Set(sky.Hdop)
		gpsPdop.WithLabelValues(sky.Device).Set(sky.Pdop)
		gpsGdop.WithLabelValues(sky.Device).Set(sky.Gdop)

		gpsNumSats.WithLabelValues(sky.Device).Set(float64(len(sky.Satellites)))
		gpsSvSs.Reset()
		num_sats_used := 0
		for i := 0; i < len(sky.Satellites); i++ {
			num_sats_used += 1
			prn_str := fmt.Sprintf("%.0f", sky.Satellites[i].PRN)
			gpsSvAz.WithLabelValues(sky.Device, prn_str).Set(sky.Satellites[i].Az)
			gpsSvEl.WithLabelValues(sky.Device, prn_str).Set(sky.Satellites[i].El)
			gpsSvSs.WithLabelValues(sky.Device, prn_str).Set(sky.Satellites[i].Ss)
		}
		gpsNumSatsUsed.WithLabelValues(sky.Device).Set(float64(num_sats_used))
	})

	done := gps.Watch()

	fmt.Println("Listening to HTTP requests at", *listen_http)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*listen_http, nil)

	<-done
}
