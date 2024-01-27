# Q3 Rcon Proxy

A modification of [lilproxy][lilproxy_url] that forwards only Q3 rcon/query packets. Useful for separating the rcon port from the game server port.

### Use

Run one or multiple rcon proxies by setting an environment variable `Q3RCON_PROXY`

for example:

```bash
export Q3RCON_PROXY="20000:28960;20001:28961;20002:28962"
```

This would run 3 proxy servers listening on ports `20000`, `20001` and `20002` that redirect rcon requests to game servers on ports `28960`, `28961` and `28962` respectively.

### Why

Avoid sending plaintext rcon requests (that include the password) to public ports. Instead send them to whitelisted ports.

Gives you the option to disable remote rcon entirely and have the server accept requests only from localhost.

### Special Thanks

[Dylan][user_link] For writing this proxy.

[lilproxy_url]: https://github.com/dgparker/lilproxy
[user_link]: https://github.com/dgparker