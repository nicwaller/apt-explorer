# apt-explorer

## Design Goals

Make it easy to:

- Browse and search any Apt repository using only a web browser. This opens up the possibility of exploring Apt repos from mobile devices (eg. iOS).
- Compare different mirrors of a given Apt repository. (`apt` and `apt-get` and `aptitude` aren't much help with that.)

## TODOs

- support multiple backends (HTTP/FTP/S3/filesystem)
- compare mirrors (at the level of archive or distribution?)
- fuzzy search packages
- show available versions
- respond with gzip-compressed JSON and/or NDJSON?
- respond to multiple Accept: formats (including application/json, text/plain, and text/html)
- Show a validation badge when the current view is valid with hashes and signatures, and click the badge to see an explanation of the trust hierarchy along with exact hashes
- Search across one or more components in a single operation (checkbox per component. all components by default)
- Decode a `deb` or `deb-src` line, showing what files would be retrieved. Use this as an entry point to start browsing a repo.
- If `ls-lR.gz` (or similar) is available, compare that with dists/ metadata to identify orphaned objects (not in any dist)
- When viewing a package or specific package version, show the `apt` and  `apt-get` and `aptitude` and `apt-cache` commands used to install that package.
  - `apt-get install <package name>=<version>`
  - ~also show the `add-apt-repository` command needed to make the repo available~ actually don't, it's not idempotent to run this
- for a given `.deb` or `.dpkg`, allow peeking inside to see the list of files contained, [just like Ubuntu Packages site](https://packages.ubuntu.com/jammy/amd64/squid/filelist)
- when viewing a package, show which architectures are available (and how that compares to the list of architectures supported in the distribution) like shown in the table at the bottom of [this page](https://packages.ubuntu.com/jammy/squid)
- maybe also get inspired by this sweet view https://pdb.finkproject.org/pdb/package.php/squid-unified
- Websockets for very fast response time
- have a configuration file for known repos, and also whether people are allowed to browse other repos

## Out of Scope

- Making a sweet terminal user interface (TUI) using something like [this](https://github.com/rivo/tview)

## See Also

https://github.com/google/apt-golang-s3
