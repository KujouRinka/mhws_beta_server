## mhws_beta_server

Here’s a project simulating server requests for the Monster Hunter Wilds Beta Test version.

这是模拟 Monster Hunter Wilds Beta test 版本服务器请求的项目。

### Disclaimer
This project is intended solely for educational and learning purposes and must not be used for commercial purposes. I assume no legal responsibility for any misuse of this project.

免责声明：本项目内容不得用于商业用途，仅做学习交流。本人尚不承担由于不当使用该项目所产生的法律责任。

### Notes
The `cert` directory contains certificate files generated for mhws.io, but it is recommended that you generate your own certificates.

`cert` 目录下为 `mhws.io` 生成的证书文件，但建议你自己生成。

### Known Issues
没试过能不能从空档启动游戏，随缘吧，有需求的话后面再说

### Usage

简单的教程，不详细，看不懂的可以搜，也可以问我，但希望能够自己先思考

TL;DR

#### 需要的东西

- 绕过 steam 启动 steam 应用的东西. e.g. `steamclient`, `unsteam` 之类的
- 能够捕获进程 http/https 流量并进行修改的工具. e.g. `mitmproxy` + `Proxifier`, `Charles`

并不绝对要使用什么工具，找到自己最适合的工具就行

这里以 `mitmproxy` + `Proxifier` 为例介绍怎么操作，其他工具原理一致，换汤不换药

**注：更推荐使用 `Charles`，稳定一些**

#### 不详细步骤

建议按照以下顺序操作，虽然有些时候顺序并不重要：

- 安装 `python3`，建议 3.8 以上，并将 `python` 添加到你的系统 `PATH` 里，用于运行 `mitmproxy`（善用搜索工具）
- 安装 `go` 或者说 `golang`，要求 `1.23.1` 及以上，并将 `go` 添加到系统 `PATH` 里，用于运行 `mhw_beta_server`
  - 并不是强制要求 `1.23.1` 及以上，用低版本的可以改 `go.mod` 第三行
  - 这里没有提供二进制程序，就麻烦各位自行编译啦（后面会提到），中途可能会遇到的网络问题也能通过搜索引擎解决
  - 是不是感觉写得很糙，对！一是我懒，二是这样就没法大面积传播了
- 安装 `Proxifier`
- 安装 `mitmporxy`。打开一个**新的**终端，`cmd` 也好 `powershell` 也好，通过下面命令安装 `mitmporxy`:
  ```bash
  pip install mitmproxy
  ```
  - 可能会遇到网络问题，搜索引擎解决 pip 网络问题
- 配置 `hosts` 文件，在 `hosts` 中写入如下内存：
  ```text
  your_ip    your_host
  ```
  - 其中，`your_ip` 为伪装服务器 `mhws_beta_server` 的监听地址，`localhost` 也好，局域网地址也好，公网 IP 也行
  - `your_host` 为你的服务的域名，随便一个合法域名就行，这里默认的是 `mhws.io`，你也可以使用其他域名，但需要你自行签名证书并丢到 `cert` 目录下
  - 例如，可以像下面这样写：
    ```text
    127.0.0.1    mhws.io
    ```
- 配置 `Proxifier`
  - Profile -> Proxy Servers -> Add...，按以下配置好后 OK：
    ```text
    Address: 127.0.0.1
    Port: 8080
    Protocol: HTTPS
    ```
  - Profile -> Proxification Rules
    - 将 `Default` Rule Name 的 Action 更改为 `Direct`
    - 选择 Add，勾选 `Enabled`，按以下配置：
      ```text
      Name: Python
      Applications: python.exe
      Action: Direct
      ```
      OK
    - 选择 Add，勾选 `Enabled`，按以下配置：
      ```text
      Name: Wilds
      Application: monsterhunterwildsbeta.exe
      Action: Proxy HTTPS 127.0.0.1
      ```
      OK
  - 保持 `Proxifier` 开启状态
- 配置 `mitmproxy`
  - 安装 `mitmproxy` 证书（搜索引擎可查）
  - 找一个位置将以下内容写入 `main.py` 文件：
    ```python
    import mitmproxy.http
    
    class CapcomPatcher:
      def __init__(self):
        pass

      def request(self, flow: mitmproxy.http.HTTPFlow):
        flow.request.host = "your_host"
        return

    addons = [
      CapcomPatcher(),
    ]
    ```
    **注意**：上端代码中 `flow.request.host = "your_host"` 请将 your_host 替换为上面 `hosts` 修改时的值，默认为 `mhws.io`
  - 在该位置打开一个**新的**终端，运行以下命令并保持终端开启：
    ```bash
    mitmweb.exe --ssl-insecure -s ./main.py
    ```
- 启动 `mhw_beta_server`
  - 在这个项目的目录下打开一个终端，运行以下命令并保持终端开启：
    ```bash
    go run . your_ip
    ```
    **注意**：上端命令中 `your_ip` 请将 your_ip 替换为上面 `hosts` 修改时的值，可能为本地回环地址，局域网地址，或者公网 IP
- 绕过 steam 启动游戏，方法有很多，搜索引擎

### PS
`Proxifier` + `mitmproxy` 只是本人用的方案，方便调试而已，实际用起来挺麻烦的，平常使用推荐 `Charles`

最后，善用搜索！善用搜索！善用搜索！
