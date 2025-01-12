# awscost

Pretty print AWS costs over time.

Given no arguments, the output displays monthly net unblended cost for the previous 5 months through the end of the current month.
Flags can be used to adjust the cost metric, granularity, and date range.

```console
% awscost
╭────────────┬────────────┬───────────────┬───────────╮
│ START      │ END        │ AMOUNT        │ % CHANGE  │
├────────────┼────────────┼───────────────┼───────────┤
│ 2024-08-01 │ 2024-09-01 │ $ 75.92       │           │
│ 2024-09-01 │ 2024-10-01 │ $ 186.15      │ 59.22 %   │
│ 2024-10-01 │ 2024-11-01 │ $ 21.54       │ -764.39 % │
│ 2024-11-01 │ 2024-12-01 │ $ 13.32       │ -61.67 %  │
│ 2024-12-01 │ 2025-01-01 │ $ 17.69       │ 24.69 %   │
│ 2025-01-01 │ 2025-02-01 │ $ 42.43 (est) │ 58.31 %   │
├────────────┼────────────┼───────────────┼───────────┤
│            │ AVERAGE    │ $ 59.51       │           │
╰────────────┴────────────┴───────────────┴───────────╯
```

## Installation

Via `go install`:

```console
go install github.com/jar-b/awscost@latest
```

## Usage

```console
% awscost -h
Pretty print AWS costs over time.

Given no arguments, the output displays monthly net unblended cost for
the previous 5 months through the end of the current month. Flags
can be used to adjust the cost metric, granularity, and date range.

Usage: awscost [flags]

Flags:
  -end string
        Usage end date (default "YYYY-MM-DD")
  -granularity string
        Cost granularity (default "MONTHLY")
  -metric string
        Cost metric (default "NetUnblendedCost")
  -start string
        Usage start date (default "YYYY-MM-DD")
```

