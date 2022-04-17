# gobundle
**Caution:** `gobundle` is in alpha phase. It will work some cases and will fail
for many more cases. If `gobundle` does not work for some of your use case, feel free
to raise an issue with minimal sample to reprodue the issue. PRs are welcomed
and appreciated. Check [Development guidelines](#development-guidelines) for
more info on development.

---

A CLI tool to copy code from dependencies and bundle it into a single go file.

What issue does this solve?
<br>
In competitive programming, you are supposed to submit a single source code file.
This results in adding some boilerplate code in all the files.

Instead, you can write modular code with go modules and use `gobundle` at the
end to bundle entire code into a single file that you can submit.

Check [test\_project0/main.go](./test_files/test_project0/main.go) and code generated with `gobundle` at
[expected\_output0/main.go](./test_files/expected_output0/main.go) for example.

## Installation

### Binaries
Download binary from [latest release](https://github.com/njkevlani/gobundle/releases/latest).


### From source code
```shell
make install
```

Make sure your `$GOPATH/bin` is part of `$PATH` environment variable.

## Usage
```shell
$ gobundle -h
version: gobundle-dev
usage: gobundle [-h] file.go

file.go    path to input go file.
-h         show this message.
```

```shell
$ gobundle target_file.go
```

This will save output to `./build/main.go`.

## Development guidelines
To build `gobundle`:
```shell
make build
```

To run tests:
```shell
make test
```

To install `gobundle` from source:
```shell
make install
```
