# Cannula

"A cannula is a tube that can be inserted into the body, often for the delivery or removal of fluid or for the gathering of data"

<img width="30%" src="https://raw.githubusercontent.com/retailnext/cannula/master/gopher_cannula.png">

*Go Gopher designed by [Renee French](http://reneefrench.blogspot.com/), original png created by [Takuya Ueda](http://u.hinoichi.net) licensed under [CC 3.0 Attribution](http://creativecommons.org/licenses/by/3.0/deed.ja)*

Cannula is Go debug package that exposes debug information about your application to a unix socket. Currently it exposes the same data as "net/http/pprof" and "expvar". You can also register your own debug handlers to expose additional information and provide additional debug actions.

## Why use Cannula instead of net/http/pprof directly?

'net/http/pprof' registers its debug handlers with the default http ServeMux. This means that unless you have a proxy in front of your webserver blocking requests to /debug those handlers will be exposed to the internet. Cannula exposes the same data as net/http/pprof on a unix socket to make it much harder to accidentally expose your debug handlers to anything but localhost.

## Connecting to the Cannula socket

If you are using a recent version of curl (>= 7.40) you can use the `--unix-socket` to connect directly to the cannula socket. Otherwise you can use the `cannula-proxy` tool which creates a localhost proxy from an ephemeral tcp port to the unix socket. This is also how you can use `go tool pprof` directly against the running go application.
