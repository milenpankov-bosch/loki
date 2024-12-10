package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	dto "github.com/prometheus/client_model/go"

	"github.com/grafana/loki/v3/pkg/loki"
	"github.com/grafana/loki/v3/pkg/util/cfg"
)

func main() {

	var config loki.ConfigWrapper

	if err := cfg.DynamicUnmarshal(&config, os.Args[1:], flag.CommandLine); err != nil {
		fmt.Fprintf(os.Stderr, "failed parsing config: %v\n", err)
		os.Exit(1)
	}

	fs := flag.NewFlagSet("_", flag.ExitOnError)
	cfg := &loki.Config{}
	cfg.RegisterFlags(fs)

	proc, _ := loki.New(config.Config)
	err := proc.Prepare()
	if err != nil {
		panic(err)
	}

	metrics, _ := proc.Registerer.Gather()
	slices.SortFunc(metrics, func(a, b *dto.MetricFamily) int { return strings.Compare(*a.Name, *b.Name) })

	for _, m := range metrics {
		fmt.Println(m.GetName(), m.GetType())
		fmt.Println("   ", m.GetHelp())
	}
}
