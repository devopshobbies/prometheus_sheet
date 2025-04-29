package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current CPU temperature in Celsius degrees",
	})
)

func init() {
	// ثبت متریک با Prometheus
	prometheus.MustRegister(cpuTemp)
}

func getCPUTemperature() (float64, error) {
	// راه‌های مختلف خواندن دمای CPU بر اساس سیستم عامل
	if _, err := os.Stat("/sys/class/thermal/thermal_zone0/temp"); err == nil {
		// برای لینوکس (رایج‌ترین روش)
		data, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
		if err != nil {
			return 0, err
		}
		tempStr := strings.TrimSpace(string(data))
		tempMilliC, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			return 0, err
		}
		return tempMilliC / 1000.0, nil
	} else if _, err := os.Stat("/proc/acpi/thermal_zone/THM0/temperature"); err == nil {
		// برای برخی سیستم‌های لینوکس قدیمی
		data, err := os.ReadFile("/proc/acpi/thermal_zone/THM0/temperature")
		if err != nil {
			return 0, err
		}
		parts := strings.Fields(string(data))
		if len(parts) < 2 {
			return 0, fmt.Errorf("unexpected temperature format")
		}
		tempStr := strings.TrimSuffix(parts[1], "C")
		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			return 0, err
		}
		return temp, nil
	}

	return 0, fmt.Errorf("no known temperature file found")
}

func updateMetrics() {
	for {
		temp, err := getCPUTemperature()
		if err != nil {
			log.Printf("Error getting CPU temperature: %v", err)
		} else {
			cpuTemp.Set(temp)
			log.Printf("CPU temperature updated: %.2f°C", temp)
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	// شروع به روزرسانی متریک‌ها در background
	go updateMetrics()

	// راه‌اندازی HTTP server برای Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	port := "8000"
	log.Printf("Starting CPU temperature exporter on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}