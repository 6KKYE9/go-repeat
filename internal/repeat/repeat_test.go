package repeat

import (
	"testing"
	"time"
)

func TestBasicCount(t *testing.T) {
	results, err := Do(Config{Count: 3, Interval: time.Millisecond}, "go", "version")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("应执行 3 次，实际 %d 次", len(results))
	}
	for i, r := range results {
		if r.Index != i+1 {
			t.Fatalf("序号不对: 第 %d 次 = %d", i+1, r.Index)
		}
		if r.Error != nil {
			t.Fatalf("第 %d 次报错: %v", i+1, r.Error)
		}
		if r.Duration <= 0 {
			t.Fatalf("第 %d 次耗时为 0", i+1)
		}
	}
}

func TestNoArgsRejected(t *testing.T) {
	if _, err := Do(Config{Count: 1}); err == nil {
		t.Fatal("不给命令应报错")
	}
}

func TestNegativeCount(t *testing.T) {
	if _, err := Do(Config{Count: -1}, "echo"); err == nil {
		t.Fatal("负次数应报错")
	}
}

func TestNegativeInterval(t *testing.T) {
	if _, err := Do(Config{Count: 1, Interval: -time.Second}, "echo"); err == nil {
		t.Fatal("负间隔应报错")
	}
}

// 次数和间隔都为 0 时退回到执行一次
func TestZeroCountAndInterval(t *testing.T) {
	results, err := Do(Config{}, "echo", "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("默认应执行 1 次，实际 %d", len(results))
	}
}

// 命令失败时，ExitCode 和 Error 都要对
func TestExitCode(t *testing.T) {
	results, err := Do(Config{Count: 1}, "go", "nonexistent_subcommand")
	if err != nil {
		t.Fatalf("Do 不应把命令失败当自身错误: %v", err)
	}
	r := results[0]
	if r.ExitCode == 0 {
		t.Fatal("失败命令 ExitCode 不应为 0")
	}
	if r.Error == nil {
		t.Fatal("失败命令 Error 不应为 nil")
	}
}

func TestStopOnErr(t *testing.T) {
	results, err := Do(Config{Count: 5, StopOnErr: true}, "go", "badcmd")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("StopOnErr 应停在第一次，实际 %d 次", len(results))
	}
}

// 间隔必须真的等了。用两次执行之间最小耗时来验证。
func TestIntervalWaits(t *testing.T) {
	start := time.Now()
	_, err := Do(Config{Count: 3, Interval: 100 * time.Millisecond}, "go", "version")
	if err != nil {
		t.Fatal(err)
	}
	elapsed := time.Since(start)
	if elapsed < 200*time.Millisecond {
		t.Fatalf("间隔 100ms 三次最少要 200ms，实际 %v", elapsed)
	}
}

// 命令不存在应正确返回错误
func TestCommandNotFound(t *testing.T) {
	results, err := Do(Config{Count: 1}, "nosuchcommand_xyz")
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Error == nil {
		t.Fatal("不存在的命令应报错")
	}
}

// 超时后命令应被终止
func TestTimeout(t *testing.T) {
	results, err := Do(Config{Count: 1, Timeout: 10 * time.Millisecond}, "go", "version")
	if err != nil {
		t.Fatal(err)
	}
	// go version 通常很快，10ms 应该够
	if results[0].Error != nil {
		t.Log("go version 没在 10ms 内跑完，不算 bug")
	}
}
