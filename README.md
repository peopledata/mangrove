# mangrove

数据消费者发布数据需求，签署数据合约以及数据资产管理。

## 功能说明

该项目为数据消费者发布数据需求的后台管理系统，实现了后台管理系统功能，包括接口和前端系统；还包括数据消费端 Papaya 需要的需求数据接口。

当管理员在后台发布一条数据需求的时候会自动在 `goerli` 链上部署一条智能合约，部署完成后 Papaya 前端页面可以看到该需求数据，然后可以将自己的数据 mint 成一个 NFT 资产，将 NFT 资产授权给需求方使用。

## 项目依赖

- 一个 MySQL 数据库
- 一个能够连接以太坊测试网络 goerli 的钱包
- 还需要能够访问以太坊 RPC 地址（我们使用的是 alchemy 服务）

## 配置文件

```yaml
app:
  name: "patronus"
  mode: "debug" # debug、test、release
  port: 8081
  log_level: "debug"
  start_time: "2023-02-06"
  machine_id: 1

mysql:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "root321"
  database: "patronus_demo"
  max_open_conns: 100
  max_idle_conns: 10

nft:
  network: "goerli"
  infura_api_key: "a62e439c8c1048b6a1f983e5d9a0e72d"
  goerli_private_key: "ce76af7fc0aca89f3a769b2ebaa236d21faa173361b9103502400116863dc71f"
  etherscan_api_key: "2VII71ZA9Q56RTWWGGK7F42WVPHHQP2KUG"
```

- 其中 app.mode 表示运行模式，可选值为：debug、test、release
- mysql 下面是数据库相关的配置
- nft 下面包括两个 key，infura_api_key 需要到 infura.io 网站注册获取 key，goerli_private_key 是用于部署合约的钱包私钥账号

`ui` 目录下面是前端项目。
