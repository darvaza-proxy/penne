# Penne, darvaza's DNS Resolver

[![Go Reference][godoc-badge]][godoc]
[![Go Report Card][goreport-badge]][goreport]

[godoc]: https://pkg.go.dev/darvaza.org/penne
[godoc-badge]: https://pkg.go.dev/badge/darvaza.org/penne.svg
[goreport]: https://goreportcard.com/report/darvaza.org/penne
[goreport-badge]: https://goreportcard.com/badge/darvaza.org/penne

_Penne_ is a config-driven pipeline oriented DNS resolver that allows complex
workflows to be defined in a simple way.
_Penne_ is built using the [darvaza sidecar engine][sidecar] and
the [darvaza resolver interface][resolver].

[core]: https://pkg.go.dev/darvaza.org/core
[resolver]: https://pkg.go.dev/darvaza.org/resolver
[sidecar]: https://pkg.go.dev/darvaza.org/sidecar
[slog]: https://pkg.go.dev/darvaza.org/slog

[split-horizon]: https://en.wikipedia.org/wiki/Split-horizon_DNS

## Horizons

_Penne_ is designed upon the idea of [_split horizons_][split-horizon],
where DNS answers depend on the IP address of the client.

A _Horizon_ is a named set of network patterns (aka `CIDR`) that can optionally
choose a custom `Resolver`,
and can annotate or filter requests before passing them to the next _Horizon_ on
a chain.

## Resolvers

On the config file you define a series of _resolvers_ in charge of
handling DNS requests.
Each _Resolver_ has a unique _name_. Names are not case sensitive and allow unicode text.

_Resolvers_ have three operation modes:

* _Iterative_ goes to the root servers and iterates through authoritative
  servers until the answer is found.
* _Forwarder_  connects to a specific server to get the answer, optionally
  allowing recursion to be performed remotely.
* and _Chained_, where requests are passed to the _Next_ resolver, optionally modified.

_Resolvers_ act as middlewares, optionally restricted to specific domains (suffixes).

_Resolvers_ can also be configured to discard various entries (like `AAAA` for example)
and execute request rewrites.

## Server

_TBD ..._

### Installation

_TBD ..._

### Configuration

_TBD ..._

### Run as service

_TBD ..._

## Web Interface

_TBD ..._

### Frontend

_TBD ..._

## See also

* [JPI Technologies' Open Source Software](https://oss.jpi.io/)
* [Split-horizon DNS (wikipedia)][split-horizon]
* [darvaza.org/core][core]
* [darvaza.org/resolver][resolver]
* [darvaza.org/sidecar][sidecar]
* [darvaza.org/slog][slog]
