# apt-transport-s3-go

apt-transport-s3-go is a Go port of [apt-transport-s3](https://github.com/MayaraCloud/apt-transport-s3).

[![test](https://github.com/winebarrel/apt-transport-s3-go/actions/workflows/test.yml/badge.svg)](https://github.com/winebarrel/apt-transport-s3-go/actions/workflows/test.yml)

## Installation

```sh
# download from https://github.com/winebarrel/apt-transport-s3-go/releases
dpkg -i apt-transport-s3-go_x.x.x_amd64.deb
```

## Usage

```sh
echo 'deb s3://ap-northeast-1@my-bucket/repo/ xenial main' > /etc/apt/sources.list.d/s3.list
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
