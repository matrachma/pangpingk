package main

import (
	"crypto/x509"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/matrachma/pangpingk"
)

func main() {
	fset := flag.NewFlagSet(appName, flag.ExitOnError)
	fset.Usage = func() {
		printUsage(os.Stderr, usageShort)
	}
	tcpOnly := fset.Bool("tcponly", false, "")
	count := fset.Int("c", defaultIterations, "")
	timeout := fset.Float64("t", defaultTimeout, "")
	jsonOutput := fset.Bool("json", false, "")
	csvOutput := fset.Bool("csv", false, "")
	csvNoHeaderOutput := fset.Bool("csv-no-header", false, "")
	insecure := fset.Bool("insecure", false, "")
	ca := fset.String("ca", "", "")
	version := fset.Bool("version", false, "")
	help := fset.Bool("help", false, "")
	fset.Parse(os.Args[1:])

	if *version {
		printVersion(os.Stderr)
		os.Exit(0)
	}
	if *help {
		printUsage(os.Stderr, usageLong)
		os.Exit(0)
	}
	args := fset.Args()
	if len(args) != 1 {
		errlog.Printf("missing server address\n")
		printUsage(os.Stderr, usageShort)
		os.Exit(1)
	}
	serverAddr := args[0]
	if *count <= 0 {
		*count = 1
	}
	if *count > maxCount {
		errlog.Printf("number of allowed connections cannot exceed %d\n", maxCount)
		printUsage(os.Stderr, usageShort)
		os.Exit(1)
	}
	if *timeout <= 0 {
		*timeout = defaultTimeout
	}
	caCerts, err := loadCaCerts(*ca)
	if err != nil {
		errlog.Printf("%s\n", err)
		printUsage(os.Stderr, usageShort)
		os.Exit(1)
	}
	if *jsonOutput && (*csvOutput || *csvNoHeaderOutput) {
		errlog.Printf("choose only one output format\n")
		printUsage(os.Stderr, usageShort)
		os.Exit(1)
	}
	config := pangpingk.Config{
		Count:              *count,
		AvoidTLSHandshake:  *tcpOnly,
		InsecureSkipVerify: *insecure,
		RootCAs:            caCerts,
		Timeout:            *timeout,
	}
	result, err := pangpingk.Ping(serverAddr, &config)
	if err != nil {
		errlog.Printf("error connecting to '%s': %s\n", serverAddr, err)
		if !(*jsonOutput || *csvOutput || *csvNoHeaderOutput) {
			os.Exit(1)
		}
	}
	s := "TLS"
	if *tcpOnly {
		s = "TCP"
	}
	if !(*jsonOutput || *csvOutput || *csvNoHeaderOutput) && err == nil {
		outlog.Printf("%s connection to %s (%s) (%d connections)\n", s, serverAddr, result.IPAddr, *count)
		outlog.Printf("min/avg/max/stddev = %s/%s/%s/%s\n", result.MinStr(), result.AvgStr(), result.MaxStr(), result.StdStr())
		os.Exit(0)
	}

	// Format the result in JSON
	jsonRes := JsonResult{
		Datetime:   time.Now().Format("2006-01-02 15:04:05"),
		Host:       result.Host,
		IPAddr:     result.IPAddr,
		ServerAddr: result.Address,
		Connection: s,
		Min:        result.Min,
		Max:        result.Max,
		Count:      result.Count,
		Avg:        result.Avg,
		Std:        result.Std,
		Ping:       1,
	}
	if err != nil {
		jsonRes.Error = fmt.Sprintf("%s", err)
		jsonRes.Min = *timeout
		jsonRes.Max = *timeout
		jsonRes.Avg = *timeout
		jsonRes.Std = 0
		jsonRes.Ping = 0
	}

	if *jsonOutput {
		b, err := json.Marshal(jsonRes)
		if err != nil {
			errlog.Printf("error producing JSON: %s\n", err)
			os.Exit(1)
		}
		_, err = os.Stdout.Write(b)
		if err != nil {
			errlog.Printf("error writing JSON: %s\n", err)
			os.Exit(1)
		}
	} else {
		if *csvOutput {
			jsonRes.CSVHeader(os.Stdout)
		}
		jsonRes.CSVRow(os.Stdout)
	}
	os.Exit(0)
}

type JsonResult struct {
	Datetime   string  `json:"datetime"`
	Host       string  `json:"host"`
	IPAddr     string  `json:"ip"`
	ServerAddr string  `json:"address"`
	Connection string  `json:"connection"`
	Count      int     `json:"count"`
	Min        float64 `json:"min"`
	Max        float64 `json:"max"`
	Avg        float64 `json:"average"`
	Std        float64 `json:"stddev"`
	Error      string  `json:"error"`
	Ping       int     `json:"ping"`
}

func (*JsonResult) CSVHeader(w io.Writer) {
	cw := csv.NewWriter(w)
	err := cw.Write([]string{"datetime", "host", "ip", "address", "connection", "count", "min", "max",
		"average", "stddev", "error", "ping"})
	if err != nil {
		errlog.Printf("error writing CSV: %s\n", err)
		os.Exit(1)
	}
	cw.Flush()
}

func (jr *JsonResult) CSVRow(w io.Writer) {
	cw := csv.NewWriter(w)
	err := cw.Write([]string{jr.Datetime, jr.Host, jr.IPAddr, jr.ServerAddr, jr.Connection,
		fmt.Sprintf("%d", jr.Count),
		fmt.Sprintf("%.3f", jr.Min),
		fmt.Sprintf("%.3f", jr.Max),
		fmt.Sprintf("%.3f", jr.Avg),
		fmt.Sprintf("%.3f", jr.Std), jr.Error,
		fmt.Sprintf("%d", jr.Ping)})
	if err != nil {
		errlog.Printf("error writing CSV: %s\n", err)
		os.Exit(1)
	}
	cw.Flush()
}

func loadCaCerts(path string) (*x509.CertPool, error) {
	if path == "" {
		return nil, nil
	}
	caCerts, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error loading CA certficates from '%s': %s", path, err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caCerts) {
		return nil, fmt.Errorf("error creating pool of CA certficates: %s", err)
	}
	return pool, nil
}
