[![Release binary builder](https://github.com/dyammarcano/secure_message/actions/workflows/release.yml/badge.svg)](https://github.com/dyammarcano/secure_message/actions/workflows/release.yml)
[![ci-test](https://github.com/dyammarcano/secure_message/actions/workflows/ci.yml/badge.svg)](https://github.com/dyammarcano/secure_message/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dyammarcano/secure_message)](https://goreportcard.com/report/github.com/dyammarcano/secure_message)
[![Go Reference](https://pkg.go.dev/badge/github.com/dyammarcano/secure_message.svg)](https://pkg.go.dev/github.com/dyammarcano/secure_message)

# Secure message

`this project is in development`

## TODO 

- [ ] Add tests
- [ ] Add more documentation
- [ ] Add more examples
- [ ] Add more error handling
- [ ] Add more logging
- [ ] Add more flags
- [ ] Add more options
- [ ] Add more features

## Introduction

This is a simple command line tool that encrypts and decrypts messages. The message is encrypted with a simple algorithm and then encoded in base58.

## Description

The message is encrypted with a simple algorithm. The algorithm is as follows:

1. The message is split into blocks of 5 characters.
2. Each block is encrypted by adding the ASCII value of the first character to each of the other characters.
3. The encrypted blocks are then concatenated and returned.
4. The encrypted message is then encoded in base58.
5. The encoded message is then returned.

encrypted `hello world` is: `2nvKpky2oAeWHBzsbRkNzpFPMCbu3St8gxapX8akD1wBDMnzUoneXYw`

## How to install

```bash
$ go install github.com/dyammarcano/secure_message/cmd/sm@latest
```

## How to use

```bash
$ sm encrypt -i "decrypted.txt" -o "encrypted.txt"
$ sm decrypt -i "encrypted.txt" -o "decrypted.txt"
```