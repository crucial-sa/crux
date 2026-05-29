# crux

The official CLI for [Crucial](https://crucial.sa) — manage your services from the terminal.

## Installation

### Install script (macOS/Linux)

```sh
curl -fsSL https://crucial.sa/install.sh | sh
```

### Direct download

Download the latest binary for your platform from the [releases page](https://github.com/crucial-sa/crux/releases).

## Running locally

We use [mise](https://mise.jdx.dev/) to manage tools, environment variables and task running. That means you only need to have mise installed and activated and mise will take care of the rest!

Running non-production build:

```sh
mise dev
```

Generating a new command:

```sh
mise gen:cmd
```

```
```

Bugs and feature requests go to [GitHub Issues](https://github.com/crucial-sa/crux/issues).
