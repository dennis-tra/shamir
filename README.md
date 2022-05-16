# Shamir

[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

A thin CLI frontend for [Hashicorp Vault's](https://www.vaultproject.io/) [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing) [implementation](https://github.com/hashicorp/vault/tree/main/shamir). This allows you to split a secret into `x` shares, and then combine them back into a single secret using any `y` of those shares with `y <= x`.

```shell
$ shamir split secret.txt # default: split into 5 shares where you need any 3 to restore the secret (numbers configurable)
```

```shell
$ shamir restore shares.txt # shares.txt should contain at least 3 newline separated shares from above
```

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Background](#background)
- [Install](#install)
- [Usage](#usage)
  - [Split](#split)
  - [Restore](#restore)
- [Related Efforts](#related-efforts)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Background

Extra care needs to be taken when dealing with secrets. Therefore, this tiny tool is designed to be:

- **credible** - It uses [Hashicorp Vault's](https://github.com/hashicorp/vault/tree/main/shamir) Shamir's Secret Sharing implementation.
- **minimal** - It only uses Go's standard library besides the above dependency.
- **approachable** - The 131 LoC should be quick and easy to audit.

Further it's:

- **configurable** - You can configure the number of shares and the number of shares needed to restore the secret.
- **composable** - You can pipe [stdin](https://en.wikipedia.org/wiki/Standard_streams) to it and use it in scripts.

## Install

When you are dealing with secrets I would recommend compiling the 131 LoC yourself instead of relying on a binary distribution:

```shell
go install github.com/dennis-tra/shamir@latest
```

Make sure the `$GOPATH/bin` is in your `PATH` variable to access the installed `shamir` executable.

## Usage

### Split

Let's imagine you have confidential data in a file called `secret.txt`. You can then run any of the following commands:

```shell
$ shamir split secret.txt
$ shamir split -shares 10 -threshold 5 secret.txt
$ cat secret.txt | shamir split
```

Example:

```shell
$ echo "My very secret secret." | shamir split -shares 4 -threshold 3
gU3GKbSg3CpSHtC+04y8OH9mtIdiq2tm
GXmZfZhoqRAgzGO+fULXEXfDusDJuCcX
ByQs4+phvdU2zXzMjYvjA+7qLLTke8Uk
9dV1XA0pJV2RDzLYh6qwKzjxJ+iBrd9W
```

Each line corresponds to one share of which you need any three to restore the original message.

### Restore

Let's imagine you have a file called `shares.txt` which contains more than `threshold` shares of your secret **separated by newlines**. E.g., you can then run any of the following commands:

```shell
$ shamir restore shares.txt
$ cat shares.txt | shamir restore
```

Example:

```shell
$ echo "9dV1XA0pJV2RDzLYh6qwKzjxJ+iBrd9W\nByQs4+phvdU2zXzMjYvjA+7qLLTke8Uk" | shamir restore # not enough shares
VL_��n�!�m5��Π8
$ echo "9dV1XA0pJV2RDzLYh6qwKzjxJ+iBrd9W\nByQs4+phvdU2zXzMjYvjA+7qLLTke8Uk\ngU3GKbSg3CpSHtC+04y8OH9mtIdiq2tm" | shamir restore
My very secret secret.
```

Note the `\n` characters in the `echo` command to separate the shares from above.

## Related Efforts

- [kinvolk/go-shamir](https://github.com/kinvolk/go-shamir) - A small CLI tool for Shamir's Secret Sharing written in Go, using Vault's Shamir implementation

## Maintainers

[@dennis-tra](https://github.com/dennis-tra).

## Contributing

Feel free to dive in! [Open an issue](https://github.com/RichardLitt/standard-readme/issues/new) or submit PRs.

## License

[MIT](LICENSE) © Dennis Trautwein