package repeat

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// Config 控制如何重复执行命令。
type Config struct {
	Count     int           // 重复次数，0 表示不限（靠 Interval 控制）
	Interval  time.Duration // 间隔
	Timeout   time.Duration // 每次执行超时，0 不限
	Silent    bool          // 不打印每次的输出
	StopOnErr bool          // 报错就停
}

// Result 一次执行的结果。
type Result struct {
	Index    int
	ExitCode int
	Output   string
	Error    error
	Duration time.Duration
}

// Do 重复执行命令。返回所有执行结果。
// args[0] 是可执行文件，后面是参数。
func Do(cfg Config, args ...string) ([]Result, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("没有指定要执行的命令")
	}
	if cfg.Count < 0 {
		return nil, fmt.Errorf("次数不能为负")
	}
	if cfg.Interval < 0 {
		return nil, fmt.Errorf("间隔不能为负")
	}
	if cfg.Count == 0 && cfg.Interval == 0 {
		cfg.Count = 1
	}

	var results []Result
	prog := args[0]
	progArgs := args[1:]

	for i := 0; cfg.Count == 0 || i < cfg.Count; i++ {
		start := time.Now()

		ctx := context.Background()
		if cfg.Timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
			defer cancel()
		}

		cmd := exec.CommandContext(ctx, prog, progArgs...)
		out, err := cmd.CombinedOutput()

		r := Result{
			Index:    i + 1,
			Output:   string(out),
			Duration: time.Since(start),
		}
		if err != nil {
			r.Error = err
			if ee, ok := err.(*exec.ExitError); ok {
				r.ExitCode = ee.ExitCode()
			}
		}
		results = append(results, r)

		if cfg.StopOnErr && r.Error != nil {
			break
		}

		if cfg.Count == 0 || i+1 < cfg.Count {
			time.Sleep(cfg.Interval)
		}
	}
	return results, nil
}
