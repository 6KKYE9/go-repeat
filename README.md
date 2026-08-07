# go-repeat

重复执行命令。支持次数/间隔/超时/报错停。零依赖。

## 装

```
go build -o go-repeat .
```

## 用

```
go-repeat -n 5 -- curl localhost:8080/health     # 测 5 次健康检查
go-repeat -n 0 -i 5s -- date                     # 每 5 秒一次，不限次数
go-repeat -n 3 -stop-on-err -- go build ./...    # 构建，报错就停
go-repeat -n 10 -timeout 2s -- curl slow.api     # 每次最多等 2 秒
go-repeat -n 100 -silent -- go test ./...        # 静默跑，只看汇总
```

## 选项

| 选项 | 说明 |
|---|---|
| `-n N` | 重复次数，默认 1，给 0 就不限次数 |
| `-i 间隔` | 间隔，如 `500ms`、`2s`，默认 1s |
| `-timeout 超时` | 每次最多等多久 |
| `-stop-on-err` | 任一次报非零退出码就停 |
| `-silent` | 不打印每次的输出 |

## 说明

- `--` 后面跟要重复执行的命令和参数。
- 次数和间隔都为 0 时退回到执行一次。
- 出错时退出码会被捕捉并汇总。
