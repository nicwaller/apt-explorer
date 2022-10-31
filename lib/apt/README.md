# apt

A Go library for reading package indexes from Apt repositories.

Apt repositories are also known as Debian repositories; the terms are interchangeable. Debian introduced the usage of Apt repositories, and all Debian-based distributions (eg. Ubuntu) use the same [repository structure](repo-structure.md). 

## Features

- Support multiple transport backends (filesystem or HTTP)

## Planned Features

- Can operate in `fast` (cache-enabled) or `direct` (cache-disabled) modes
- Verify consistency of indexes with multiple representations (eg. `Packages.gz` and `Packages.bz2`)
- Verify trust of each file (by verifying hashes and GPG keys)
- Read `ls-lR.gz` (if available) to get all dists in archive

## Planned Design

- Highly space-efficient representation in memory (so you can keep with more/larger repos entirely in memory). Hashes as byte arrays, not strings. Deduplication of strings for paths, etc. 
- Local cache of downloaded indexes (using content-based addressing for trust and deduplication)
- High-performance local caches (using ProtoBuf, Avro, Parquet, or similar)

## Out of Scope

**Opaque Blobs**: This library treats everything in `pool` as an opaque blob of data. There will be no attempt to inspect the contents of `deb` files; use some other library for that. _(Admittedly, it would be very cool to use HTTP [Range](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Range) headers to read the control files from a Debian package to allow displaying the contents without retrieving the entire package.)_

**Mirrors:** This library deals with one Apt repo at a time. Use [apt-mirrorset](../apt-mirrorset) instead.
