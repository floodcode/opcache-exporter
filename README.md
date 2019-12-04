# OPcache Exporter for Prometheus

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

The MIT License, see [LICENSE](/LICENSE)