package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/jedib0t/go-pretty/v6/table"
)

var (
	// flag defaults
	metric      = "NetUnblendedCost"
	granularity = string(types.GranularityMonthly)
	start       = time.Now().AddDate(0, -5, -(time.Now().Day() - 1)).Format(time.DateOnly)
	end         = time.Now().AddDate(0, 1, -(time.Now().Day() - 1)).Format(time.DateOnly)
)

func main() {
	// slightly better usage output
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Pretty print AWS costs over time.\n\n"+
			"Given no arguments, the output displays monthly net unblended cost for\n"+
			"the previous 5 months through the end of the current month. Flags\n"+
			"can be used to adjust the cost metric, granularity, and date range.\n\n",
		)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags]\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&metric, "metric", metric, "Cost metric")
	flag.StringVar(&granularity, "granularity", granularity, "Cost granularity")
	flag.StringVar(&start, "start", start, "Usage start date")
	flag.StringVar(&end, "end", end, "Usage end date")
	flag.Parse()

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("loading SDK config: %s", err)
	}

	client := costexplorer.NewFromConfig(cfg)
	out, err := client.GetCostAndUsage(ctx, &costexplorer.GetCostAndUsageInput{
		Metrics:     []string{metric},
		Granularity: types.Granularity(granularity),
		TimePeriod: &types.DateInterval{
			Start: aws.String(start),
			End:   aws.String(end),
		},
	})
	if err != nil {
		log.Fatalf("getting cost and usage: %s", err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.AppendHeader(table.Row{"Start", "End", "Amount", "% Change"})

	var prevAmount float64
	var total float64
	for _, r := range out.ResultsByTime {
		amount, _ := strconv.ParseFloat(aws.ToString(r.Total[metric].Amount), 64)
		amountf := fmt.Sprintf("$ %.2f", amount)
		if r.Estimated {
			amountf += " (est)"
		}

		pctChangef := ""
		if prevAmount > 0 {
			pctChange := (1 - (prevAmount / amount)) * 100
			pctChangef = fmt.Sprintf("%.2f %%", pctChange)
		}
		prevAmount = amount
		total = total + amount

		t.AppendRow([]interface{}{
			aws.ToString(r.TimePeriod.Start),
			aws.ToString(r.TimePeriod.End),
			amountf,
			pctChangef,
		})
	}

	avg := fmt.Sprintf("$ %.2f", (total / float64(len(out.ResultsByTime))))
	t.AppendFooter(table.Row{"", "Average", avg, ""})

	t.Render()
}
