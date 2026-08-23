# lever-phase — 二元合金杠杆定律核算

lever-phase 是一个凝固相分数核算工具：输入合金成分（质量分数）与温度，内核在 Pb–Sn 型「固溶 + 共晶」简化相图上判定相区，并按杠杆定律计算各相成分与质量分数。能力边界：只做杠杆定律核算，不做等温闪蒸（vapor-flash）与逐板精馏（distill-bin）；不处理气液平衡、不做仓库出入库。

## 用法

启动 Web 控制台：

```bash
go run . -http :8080
```

打开 http://localhost:8080，页面可一键加载 `example/alloy-60.json`（60 wt% Sn、458 K，位于 α + 液两相区的核算点），也可以直接调 API：

```bash
curl -s -X POST http://localhost:8080/api/point -H 'Content-Type: application/json' \
  --data-binary @example/alloy-60.json
```

固定成分扫温度、取液相分数 fL(T) 点列：

```bash
curl -s -X POST http://localhost:8080/api/scan -H 'Content-Type: application/json' \
  -d '{"c":0.45,"tmin":300,"tmax":700,"n":81}'
```

页面上的 fL(T) 曲线点列全部来自 `/api/scan` 的后端求解。

## 关键约定

- **相图**：经典 Pb–Sn 型简化相图。端元熔点 TA=600.61 K（Pb）、TB=505.08 K（Sn），共晶温度 TE=456.0 K、共晶成分 cE=0.619（Sn 质量分数），TE 处 α 最大溶解度 cAlphaMax=0.192、β 最小 Sn 含量 cBetaMin=0.975。
- **成分轴**：全程质量分数（0～1），不使用摩尔分数。
- **杠杆定律**：两相区 wα=(cβ−c)/(cβ−cα)、wβ=1−wα，恒有 wα+wβ=1，且 0≤w≤1；杠杆两端取的是当前温度的相界成分。合金成分恰等于某相界成分时，该相分数为 1、另一相为 0。
- **共晶反应**：温度刚低于 TE 时剩余液相全部凝固为共晶组织（α+β），质量守恒仍成立；c=cE 且 T 刚高于 TE 时全液相。
- **校验**：成分越界、温度低于 0 K、相图结点乱序都返回错误，不做静默错值。

## 构建与测试

```bash
go build ./...
go test ./...
```

## 许可

MIT。
