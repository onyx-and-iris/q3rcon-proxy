# Q3 Rcon Proxy

A modification of [lilproxy][lilproxy_url] that forwards only Q3 rcon/query packets. Useful for separating the rcon port from the game server port.

### Why

Unfortunately the Q3Rcon engine ties the rcon port to the game servers public port used for client connections. This proxy will allow you to run rcon through a separate whitelisted port.

### Use

Run one or multiple rcon proxies by setting an environment variable `Q3RCON_PROXY`

for example:

```bash
export Q3RCON_PROXY="20000:28960;20001:28961;20002:28962"
```

This would configure q3rcon-proxy to run 3 proxy servers listening on ports `20000`, `20001` and `20002` that redirect rcon requests to game servers on ports `28960`, `28961` and `28962` respectively.

Then just run the binary which you can compile yourself, download from `Releases` or use the included Dockerfile.

### Logging

Set the log level with environment variable `Q3RCON_LOGLEVEL`:

`0 = Panic, 1 = Fatal, 2 = Error, 3 = Warning, 4 = Info, 5 = Debug, 6 = Trace`

### Special Thanks

[Dylan][user_link] For writing [lilproxy][lilproxy_url].

[lilproxy_url]: https://github.com/dgparker/lilproxy
[user_link]: https://github.com/dgparker
