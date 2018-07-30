package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"strconv"
)

func buildGauge(check Check) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name: check.Name,
		Help: check.Help,
		ConstLabels: check.Labels,
	})
}

func metricsUpdater(gauges map[string]prometheus.Gauge, resultsChan chan CheckResult){
	for {
		checkResult := <- resultsChan
		gauge := gauges[checkResult.name]
		if checkResult.err == nil {
			gauge.Set(float64(checkResult.exitCode))
		}else{
			gauge.Set(255)
		}
	}
}

func mainExporter(conf *Configuration, resultsChan chan CheckResult) {
	gauges := make(map[string]prometheus.Gauge, len(conf.Checks))

	for _, check := range conf.Checks {
		gauge := buildGauge(check)
		prometheus.MustRegister(gauge)
		gauges[check.Name] = gauge
	}

	go metricsUpdater(gauges, resultsChan)

	http.Handle(conf.ExporterPath, promhttp.Handler())
	log.Fatal(http.ListenAndServe(conf.ListenIP + ":" + strconv.Itoa(conf.ListenPort), nil))
}
