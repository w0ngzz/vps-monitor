# VPS Monitor System

## 编译 Go Agent（自动）
提交 agent-go 目录内容将触发 GitHub Actions 自动构建，并产出二进制文件。

## 启动服务端
```bash
docker-compose up --build -d
```

## 配置 Agent
编辑 config.json，并上传 agent 二进制运行即可：
```bash
./agent-linux-amd64
```

## 查看服务器状态
打开浏览器访问：http://localhost:8000/status
