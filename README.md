# SubGen

一个用 Go 编写的轻量后端，将远程基础 YAML 配置与本地多个订阅聚合，生成可直接供客户端订阅的最终配置。

## 快速开始
- 在项目根目录准备 `config.yaml`（见下方配置说明）。
- 启动服务：`go run ./main.go`（默认监听 `:7081`）。
- 获取配置：访问 `http://localhost:7081/config`。

## 配置说明
- `base_config_url`：远程基础配置地址，留空将使用内置默认地址。
- `subscriptions`：订阅列表，字段包括：
  - `name`：订阅名称（用于 `proxy-providers` 与 `use.use`）。
  - `url`：订阅链接。
  - `path`：可选，缓存文件路径。
  - `additional_prefix`：可选，覆写中的附加前缀。

示例：
```
base_config_url: https://gist.githubusercontent.com/liuran001/5ca84f7def53c70b554d3f765ff86a33/raw/
subscriptions:
  - name: 订阅一
    url: https://example.com/api/v1/client/subscribe?token=xxxx
  - name: 订阅二
    url: https://another.example/subscribe?token=yyyy
```

## 接口
- `/`：健康信息与使用提示。
- `/config`：返回生成后的 YAML（`Content-Type: text/yaml; charset=utf-8`）。

## Docker
- 运行容器：`docker run -p 7081:7081 -v ./config.yaml:/app/config.yaml fengshuu/subgen:latest`

