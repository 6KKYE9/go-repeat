package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"go-repeat/internal/repeat"
)

func main() {
	n := flag.Int("n", 1, "重复次数，0 为不限")
	interval := flag.Duration("i", time.Second, "间隔（如 500ms、2s）")
	timeout := flag.Duration("timeout", 0, "每次执行超时")
	stopErr := flag.Bool("stop-on-err", false, "报错就停")
	silent := flag.Bool("silent", false, "不打印输出")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	cfg := repeat.Config{
		Count:     *n,
		Interval:  *interval,
		Timeout:   *timeout,
		StopOnErr: *stopErr,
		Silent:    *silent,
	}

	results, err := repeat.Do(cfg, args...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, r := range results {
		if !*silent {
			fmt.Print(r.Output)
		}
		if r.Error != nil {
			fmt.Fprintf(os.Stderr, "[%d] 失败 (耗时 %v, 退出码 %d)\n",
				r.Index, r.Duration.Round(time.Millisecond), r.ExitCode)
		}
	}

	nonZero := 0
	for _, r := range results {
		if r.ExitCode != 0 {
			nonZero++
		}
	}
	if nonZero > 0 {
		fmt.Fprintf(os.Stderr, "共 %d/%d 次非零退出\n", nonZero, len(results))
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `go-repeat — 重复执行命令

用法:
  go-repeat [选项] -- 命令及参数

选项:
  -n N         重复次数，默认 1，0 为不限（靠 -i 控制）
  -i 间隔      间隔时间，如 500ms、2s，默认 1s
  -timeout 超时 每次执行超时，如 5s
  -stop-on-err 任一次报错就停
  -silent      不打印输出

例子:
  go-repeat -n 5 -- curl http://localhost:8080/health
  go-repeat -n 0 -i 5s -- date
  go-repeat -n 3 -stop-on-err -- go build ./...
`)
}
