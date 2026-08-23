# lever-phase — 二元合金杠杆定律核算

lever-phase 是二元合金凝固相分数核算工具：输入合金成分（质量分数）与温度，在 Pb–Sn 型固溶 + 共晶简化相图上判定相区，并按杠杆定律计算各相成分与质量分数；内置 `example/alloy-60.json` 算例。

## 构建 / 运行 / 测试

```text
go build ./...     # 编译
go run . -http :8080   # 启动 Web 控制台，页面可加载 example/alloy-60.json
go test ./...      # 测试
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
