package main

import (
	collectors "github.com/jay-lane/evewatch/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"gopkg.in/yaml.v2"

	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultlistenaddress    = ":9800"
	namespace               = "evewatch"
    accountName             = "Lanemire"
	HoursInDay              = 24
	trafficReportInterval   = 5 // mins
	lookbackDefaultDuration = 2 * 24 * time.Hour
	prefillDefaultDuration  = 10 * time.Minute
)

var (
	configFile           = kingpin.Flag("config.file", "GTM Metrics exporter configuration file. Default: ./gtm_metrics_config.yml").Default("evewatch_config.yml").String()

	// invalidMetricChars    = regexp.MustCompile("[^a-zA-Z0-9_:]")
	lookbackDuration = lookbackDefaultDuration
	prefillDuration  = prefillDefaultDuration
    account          = accountName
)

// Calculate window duration based on config and save in lookbackDuration global variable
func calcWindowDuration(window string) (time.Duration, error) {

	var datawin int
	var err error
	var multiplier time.Duration = time.Hour * time.Duration(HoursInDay) // assume days

	log.Debugf("Window: %s", window)
	if window == "" {
		return time.Second * 0, fmt.Errorf("Summary window not set")
	}
	iunit := window[len(window)-1:]
	if !strings.Contains("mhd", strings.ToLower(iunit)) {
		// no units. default days
		datawin, err = strconv.Atoi(window)
	} else {
		len := window[0 : len(window)-1]
		datawin, err = strconv.Atoi(len)
		if strings.ToLower(iunit) == "m" {
			multiplier = time.Minute
			if err == nil && datawin < trafficReportInterval {
				datawin = trafficReportInterval
			}
		} else if strings.ToLower(iunit) == "h" {
			multiplier = time.Hour
		}
	}
	if err != nil {
		log.Warnf("ERROR: %s", err.Error())
		return time.Second * 0, err
	}
	log.Debugf("multiplier: [%v} units: [%v]", multiplier, datawin)
	return multiplier * time.Duration(datawin), nil

}

func main() {

	/*
    log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print(namespace + "metrics_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infof("Config file: %s", *configFile)
	log.Infof("Starting EveWatch. %s", version.Info())
	log.Infof("Build context: %s", version.BuildContext())

    
	eveWatchConfig, err := loadConfig(*configFile) // save?
	if err != nil {
		log.Fatalf("Error loading eveWatch config file: %v", err)
	}
    */
    client, err := evesi.NewClient(*http.Client, userAgent string)

	log.Debugf("Exporter configuration: [%v]", eveWatchConfig)

    /*
	// Use custom registry
	r := prometheus.NewRegistry()
	r.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	r.MustRegister(prometheus.NewGoCollector())
	r.MustRegister(version.NewCollector(namespace + "metrics_exporter"))
	
    
    
    r.MustRegister(collectors.NewDatacenterTrafficCollector(r, gtmMetricsConfig, namespace, tstart, lookbackDuration))
	r.MustRegister(collectors.NewPropertyTrafficCollector(r, gtmMetricsConfig, namespace, tstart, lookbackDuration))
	r.MustRegister(collectors.NewLivenessTrafficCollector(r, gtmMetricsConfig, namespace, tstart, lookbackDuration))
    r.MustRegister(collectors.ItemCollector(r, eveWatchConfig,namespace,tstart,lookbackDefaultDuration))

	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{Registry: r}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>evewatch</title></head>
			<body>
			<h1>evewatch</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})
    */
	log.Info("Beginning to serve on address ", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}

func loadConfig(configFile string) (collectors.EveWatchConfig, error) {
	if fileExists(configFile) {
		// Load config from file
		configData, err := ioutil.ReadFile(configFile)
		if err != nil {
			return collectors.EveWatchConfig{}, err
		}
		log.Debugf("EveWatch config file: %s", string(configData))
		return loadConfigContent(configData)
	}

	log.Infof("Config file %v does not exist, using default values", configFile)
	return collectors.EveWatchConfig{}, nil

}

func loadConfigContent(configData []byte) (collectors.EveWatchConfig, error) {
	domains := make([]*collectors.DomainTraffic, 0)
	domains = append(domains, &collectors.DefaultDomainTraffic)
	config := collectors.EveWatchConfig{Domains: domains}
	err := yaml.Unmarshal(configData, &config)
	if err != nil {
		return config, err
	}

	log.Info("eveWatch config loaded")
	return config, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
