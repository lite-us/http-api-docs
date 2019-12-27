# http-api-docs


> A generator for go-btfs API endpoints documentation.

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Captain](#captain)
- [Contribute](#contribute)
- [License](#license)

## Install

In order to build this project, you need to first install Go, clone this repo, and finally run `make install`:

```sh
> git clone https://github.com/TRON-US/http-api-docs "$(go env GOPATH)/src/github.com/TRON-US/http-api-docs"
> cd "$(go env GOPATH)/src/github.com/TRON-US/http-api-docs"
> make install
```

## Usage

After installing you can run:

```
> http-api-docs
```

This should spit out a Markdown document. This is exactly the `api.md` documentation at https://github.com/TRON-US/go-btfs/tree/master/docs/content/reference/api/http.md, so you can redirect the output to just overwrite that file.

## Captain

This project is captained by @hsanjuan.

## Contribute

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT (C) Protocol Labs, Inc.
