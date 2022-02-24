# apt-transport-s3-go

apt-transport-s3-go is a Go port of [apt-transport-s3](https://github.com/MayaraCloud/apt-transport-s3).

[![test](https://github.com/winebarrel/apt-transport-s3-go/actions/workflows/test.yml/badge.svg)](https://github.com/winebarrel/apt-transport-s3-go/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/winebarrel/apt-transport-s3-go)](https://goreportcard.com/report/github.com/winebarrel/apt-transport-s3-go)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwinebarrel%2Fapt-transport-s3-go.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwinebarrel%2Fapt-transport-s3-go?ref=badge_shield)

## Installation

```sh
# download from https://github.com/winebarrel/apt-transport-s3-go/releases
dpkg -i apt-transport-s3-go_x.x.x_amd64.deb
```

## Usage

```sh
# aws s3 ls s3://my-bucket/repo/
#                           PRE dists/
#                           PRE pool/

echo 'Acquire::s3::region ap-northeast-1;' > /etc/apt/apt.conf.d/s3
echo 'deb s3://my-bucket/repo/ xenial main' > /etc/apt/sources.list.d/s3.list
apt update
apt install any-pkg
```

### Debug

```sh
export ATS3_LOG_LEVEL=debug
apt update
```

## Related Links

* [apt-transport-s3 License & Copyright](https://github.com/MayaraCloud/apt-transport-s3#license--copyright)
* [APT Method Interface](http://www.fifi.org/doc/libapt-pkg-doc/method.html/index.html#abstract)
* [google/apt-golang-s3: An s3 transport method for the apt package management system](https://github.com/google/apt-golang-s3)


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwinebarrel%2Fapt-transport-s3-go.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwinebarrel%2Fapt-transport-s3-go?ref=badge_large)