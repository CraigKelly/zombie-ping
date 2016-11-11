# Zombie Ping!

![Version Status](https://img.shields.io/badge/status-experimental-red.svg)

Experimental Software! Use at your own risk.

This is a simple URL monitoring tool that uses browser notifications to let
you know if anything changes. Current target is Firefox only.

## Zero to Running

```
git clone https://github.com/CraigKelly/zombie-ping.git
cd zombie-ping
make
./zombie-ping
```

## Configuration

TODO: document config file

## Dependencies

* `make` for building
* `go` for the server
* `pysassc` (available via apt) for scss

## Hacking

main.go is the server side

The static directory contains the client (and thus the files served), with one
exception: the CSS files are generated from the SCSS files in the scss
directory.

## Testing

```
make test
```

## License

No license yet - but it **will** be OSS
