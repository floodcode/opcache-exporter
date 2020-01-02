# OPcache Exporter for Prometheus

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/crowdin/opcache-exporter)

This is a simple server that scrapes OPcache status and exports it via HTTP for Prometheus consumption.

## Development

### Building

```
make build
```

### Running

```
$ opcache_exporter [<flags>]

Flags:
  -h, --help                    Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=":9101"  
                                Address to listen on for web interface and telemetry.
      --web.telemetry-path="/metrics"  
                                Path under which to expose metrics.
      --opcache.fcgi-uri="127.0.0.1:9000"  
                                Connection string to FastCGI server.
      --opcache.script-path=""  Path to PHP script which echoes json-encoded OPcache status
      --opcache.script-dir=""   Path to directory where temporary PHP file will be created
```

## License
<pre>
Copyright © 2020 Crowdin

The Crowdin OPcache exporter is licensed under the MIT License. 
See the LICENSE file distributed with this work for additional 
information regarding copyright ownership.

Except as contained in the LICENSE file, the name(s) of the above copyright
holders shall not be used in advertising or otherwise to promote the sale,
use or other dealings in this Software without prior written authorization.
</pre>
