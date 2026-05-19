## mhws_beta_server

**[中文文档](README.zh.md)**

Here’s a project simulating server requests for the Monster Hunter Wilds Beta Test version.

> [!NOTE]  
> This project is intended for OBT1, and there are currently no plans to update it for OBT2. If you want to use it for the latter, please refer to [WiseHorror/mhws_beta_server](https://github.com/WiseHorror/mhws_beta_server).

### Disclaimer
This project is intended solely for educational and learning purposes and must not be used for commercial purposes. Please delete this software within 24 hours of downloading. We do not take any responsibility for any illegal usage or distribution of this software.

### Notes
The `cert` directory contains certificate files generated for mhws.io

### Known Issues
~~Haven't tested if it can be launched from scratch; maybe later if needed.~~
### Usage

A simple tutorial, not very detailed. If you don't understand, you can search or ask me, but please think for yourself first.
- Start `mhws_beta_server`
  ```bash
  go run . --cert-domain hjm.rebe.capcom.com,40912.playfabapi.com --api-host hjm.rebe.capcom.com
  ```
  For more information, use `--help` command
- Add `./cert/root.crt` as a trusted root certificate
- Add the following to your `hosts` file and flush the system DNS cache:
  ```text
  127.0.0.1 hjm.rebe.capcom.com
  127.0.0.1 40912.playfabapi.com
  ```
  **Note**：After adding, you may not be able to connect to other Capcom games. Restore the `hosts` file and flush the cache to fix; if necessary, try restarting your system.

#### Below is the old tutorial, for reference only, validity not guaranteed.

TL;DR

#### What's needed

- Something to bypass Steam and launch Steam applications, e.g. `steamclient`, `unsteam`
- A tool that can capture and modify process HTTP/HTTPS traffic. e.g. `mitmproxy` + `Proxifier`, `Charles`

You are not required to use any specific tool; find the one that suits you best.

Here we use `mitmproxy` + `Proxifier` as an example. Other tools work on the same principle, just different implementations.

**`Charles` is recommended as it is more stable.**

#### Steps (not detailed)

It is recommended to follow the steps in order, though sometimes the order is not important:

- Install `python3` (recommended version 3.8 or higher) and add `python` to your system `PATH` so you can run `mitmproxy` (make good use of search engines).
- Install `go` (a.k.a. `golang`), version `1.23.1` or higher, and add `go` to your system `PATH` to run `mhw_beta_server`.
  - Version `1.23.1` or higher is not strictly required; for lower versions you can modify the third line of `go.mod`.
  - No pre‑compiled binaries are provided here, so you'll need to compile it yourself (more on that later). Any network issues you encounter along the way can also be resolved with a search engine.
- Install `Proxifier`.
- Install `mitmproxy`. Open a **new** terminal (whether `cmd` or `powershell`) and run the following command:
  ```bash
  pip install mitmproxy
  ```
  - You may run into network issues – use a search engine to solve pip network problems.
- Configure the `hosts` file by adding the following line:
  ```text
  your_ip    your_host
  ```
  - Here, `your_ip` is the listening address of the fake server `mhws_beta_server` – it can be `localhost`, a LAN address, or a public IP.
  - `your_host` is the domain name for your service. Any valid domain name works; the default used here is `mhws.io`. You can use a different domain, but then you'll need to sign a certificate yourself and place it in the `cert` directory.
  - For example, you could write:
    ```text
    127.0.0.1    mhws.io
    ```
- Configure `Proxifier`:
  - Profile -> Proxy Servers -> Add... and fill in as follows, then click OK:
    ```text
    Address: 127.0.0.1
    Port: 8080
    Protocol: HTTPS
    ```
  - Profile -> Proxification Rules
    - Change the Action of the `Default` rule to `Direct`.
    - Click Add, check `Enabled`, and configure as follows:
      ```text
      Name: Python
      Applications: python.exe
      Action: Direct
      ```
      OK
    - Click Add, check `Enabled`, and configure as follows:
      ```text
      Name: Wilds
      Application: monsterhunterwildsbeta.exe
      Action: Proxy HTTPS 127.0.0.1
      ```
      OK
  - Keep `Proxifier` running.

  **Note**: Make sure that in the Rules list, `Python` is above `Wilds`, and `Wilds` is above `Localhost`.
- Configure `mitmproxy`:
  - Install the `mitmproxy` certificate (searchable via a search engine).
  - In a location of your choice, write the following content into a `main.py` file:
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
    **Note**: In the line `flow.request.host = "your_host"`, replace `your_host` with the value you used in the `hosts` file (default is `mhws.io`).
  - Open a **new** terminal in that location and run the following command, keeping the terminal open:
    ```bash
    mitmweb.exe --ssl-insecure -s ./main.py
    ```
- Start `mhw_beta_server`:
  - Open a terminal in the project directory and run the following command, keeping the terminal open:
    ```bash
    go run . your_ip
    ```
    **Note**: In the command above, replace `your_ip` with the value you used in the `hosts` file – it could be the loopback address, a LAN address, or a public IP.
- Bypass Steam to launch the game – there are many methods, search for them.

### PS
`Proxifier` + `mitmproxy` is just the solution I personally use because it's convenient for debugging. In practice it's quite cumbersome; for everyday use I recommend `Charles`.

Finally, make good use of search engines! Make good use of search engines! Make good use of search engines!

### Thanks
[@atYuguo](https://github.com/atYuguo)

[@EdLovecraft](https://github.com/EdLovecraft)

[@Evilmass](https://github.com/Evilmass)

[@pangliang](https://github.com/pangliang)
