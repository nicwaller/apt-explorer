# Apt Repository

Official documentation: https://wiki.debian.org/DebianRepository/Format

An Apt repo is a hierarchy of **Archive** -> **Distribution** -> **Component**.  Components hold a list of packages. Components are the center of everything.

## Archive

Each archive holds one or more distributions. The Ubuntu archive holds distributions for each version of Ubuntu.

> Examples:
> 
> - [http://archive.ubuntu.com/ubuntu/](http://archive.ubuntu.com/ubuntu/)
> - [http://ports.ubuntu.com/ubuntu-ports/](http://ports.ubuntu.com/ubuntu-ports/)
> - [http://us-west-2.ec2.ports.ubuntu.com/ubuntu-ports/](http://us-west-2.ec2.ports.ubuntu.com/ubuntu-ports/)

⚠️ There is no standard way to list the distributions contained in an archive. Some archives publish `ls-lR.gz` (a list of all files in the archive). Others provide HTML directory listings. Some provide neither.

If `ls-lR.gz` is available (not common), then you can probably get a list of dists by doing this:

```shell
ARCHIVE="http://archive.ubuntu.com/ubuntu"
curl -s -r 0-4096 "$ARCHIVE/ls-lR.gz" |
  gzip -d 2>/dev/null |
  awk '/^.\/dists:$/ {x=1} /^d/ {if (x==1) {print $9}} /^$/ {if (x==1) {exit 0}}'
```

💁 Ubuntu has separate archives for common architectures (`i386`,`amd64`) and "ported" architectures (`arm64`, `armhf`, `powerpc`, `ppc64el`).

💁‍ An archive typically has a single `pool` of artifacts that is shared by all the distributions in that archive. 

### Ubuntu PPA (Personal Package Archive)

An Ubuntu PPA is a standard Apt repository hosted on `ppa.launchpadcontent.net`. Ubuntu provides a special syntax for adding PPA repositories to /etc/apt/sources.list.d:

```shell
add-apt-repository ppa:deadsnakes/ppa
                       ^^^^^^^^^^^^^^
```

This creates a new file in /etc/apt/sources.list.d/ that contains a `deb` line as shown below. The distribution `jammy` is automatically detected from the current distribution by running `lsb_release -cs`.

```text
deb https://ppa.launchpadcontent.net/deadsnakes/ppa/ubuntu/ jammy main
                                     ^^^^^^^^^^^^^^         ^^^^^
```

> Examples:
> - `ppa:deadsnakes/ppa` -> [https://ppa.launchpadcontent.net/deadsnakes/ppa/ubuntu/](https://ppa.launchpadcontent.net/deadsnakes/ppa/ubuntu/)

## Distribution

A distribution is a single source for apt.

> Examples:
> - http://archive.ubuntu.com/ubuntu/dists/bionic/
> - http://archive.ubuntu.com/ubuntu/dists/jammy/
> - http://apt.corretto.aws/dists/stable/


The current contents of a distribution are defined by a `Release` file in the distribution directory.

_Tip: You can quickly compare distributions hosted on different mirrors by sending an HTTP `HEAD` request for the `Release` file and comparing the `ETag`, `Content-Length`, and `Expires` headers._ 

Each distribution is divided into one or more components.

## Component

> Examples: [main](http://archive.ubuntu.com/ubuntu/dists/trusty/main/), [universe](http://archive.ubuntu.com/ubuntu/dists/trusty/universe/), [multiverse](http://archive.ubuntu.com/ubuntu/dists/trusty/multiverse/), [restricted](http://archive.ubuntu.com/ubuntu/dists/trusty/restricted/)

Within each component there are several well-known paths.

Components have various kinds of contents: [Binaries](#binaries), [Sources](#sources), and more.

### Binary Packages

Binaries are the most commonly requested artifacts in Apt repos. Because binaries are architecture-specific, they are sorted into different directories of the form `binary-$ARCH`. This means that a host only needs to download the package index for compatible packages, and it can skip downloading the package indexes for other architectures.

_Tip: The different kinds of `binary-$ARCH` directories available depend on which architectures are listed in the distribution `Release` file._

Each directory contains a `Packages` file, but this file is often quite large, so typically only compressed versions will be available. Typical compression file formats are gzip, bzip2, and xzip.

> Example files:
> - binary-amd64/Packages.gz
> - binary-i386/Packages.bz2
> - binary-arm64/Packages.xz
> - binary-powerpc/Packages

See also:

https://wiki.debian.org/Packaging/Intro?action=show&redirect=IntroDebianPackaging

### Sources

Sources are not architecture-specific, so only one directory is needed. The index is found in a `Sources` file. These indexes can also be quite large, so they are typically offered in a compressed format just like the binary package indexes. 

> Example files:
> - source/Sources.gz
> - source/Sources.bz2

### Internationalized Text (i18n)

Translations are organized by [ISO 639-1 standard language codes](https://www.andiamo.co.uk/resources/iso-language-codes/) in the format `i18n/Translation-LANG`. Much like sources and binaries, these indexes are typically compressed.

> Example files:
> - i18n/Translation-de.bz2
> - i18n/Translation-en_GB.gz
> - i18n/Translation-fr

Each translation record has a `Description-md5` that is used to match translations with records from other indexes.

## See Also

https://www.ibiblio.org/gferg/ldp/giles/repository/repository-2.html

---

## sources.list

When using Ubuntu and other Debian-based systems, apt sources configured like this: 

```
deb [uri]                         [distribution] [component1] [component2] [...]
deb https://deb.debian.org/debian stable         main         contrib      non-free
```

This entry would correspond to the following URLs:

- http://ftp.debian.org/debian/dists/stable/main/
- http://ftp.debian.org/debian/dists/stable/contrib/
- http://ftp.debian.org/debian/dists/stable/non-free/
