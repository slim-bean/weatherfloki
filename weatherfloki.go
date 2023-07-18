package main

import (
	"flag"
	"net"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/grafana/loki-client-go/loki"
	"github.com/prometheus/common/model"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)
	logger = level.NewFilter(logger, level.AllowDebug())

	sock, err := net.ListenPacket("udp", ":50222")
	if err != nil {
		level.Error(logger).Log("msg", "unable to listen", "error", err)
		os.Exit(1)
	}

	cfg := loki.Config{}
	// Sets defaults as well as anything from the command line
	cfg.RegisterFlags(flag.CommandLine)
	flag.Parse()

	c, err := loki.NewWithLogger(cfg, logger)
	if err != nil {
		level.Error(logger).Log("msg", "failed to create client", "err", err)
	}

	byteBuf := make([]byte, 2000)
	labels := model.LabelSet{
		model.LabelName("job"): model.LabelValue("weatherflow"),
	}
	nextHeartbeat := time.Now()
	for {
		// Read
		n, _, err := sock.ReadFrom(byteBuf)
		if err != nil {
			level.Error(logger).Log("msg", "Could not read from socket", "error", err)
			time.Sleep(1 * time.Second)
			continue
		}
		bytes := byteBuf[:n]
		payload := string(bytes)
		level.Debug(logger).Log("msg", "received packet", "payload", payload)

		// Send
		err = c.Handle(labels, time.Now(), payload)
		if err != nil {
			level.Error(logger).Log("msg", "when I wrote this it was impossible for the client to return an error?", "err", err)
			continue
		}

		// No errors, if enough time has passed log a heartbeat
		if time.Now().After(nextHeartbeat) {
			level.Info(logger).Log("msg", "heartbeat")
			nextHeartbeat = nextHeartbeat.Add(1 * time.Minute)
		}
	}
}
