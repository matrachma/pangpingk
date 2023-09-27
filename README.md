# pangpingk — a tool for measuring various networks related latency

## Overview
`pangpingk` is stand for "Pang PING Keun" in Sundanese which mean "please to PING".
This is a command line tool to measure the time required to perform TCP/TLS connections with a network server.

It concurrently establishes several network connections, performs the TCP/TLS handshake on every one of them, measures the time spent handshaking and reports a summary on the observed results.

## Installation
Download a **binary release** for your target operating system from the [releases page](https://github.com/matrachma/pangpingk/releases).

Alternatively, if you prefer to **build from sources**, you need the [Go programming environment](https://golang.org). Do:

```
go get -u github.com/matrachma/pangpingk/...
```

## Credits

This tool was inspired by work of Fabio Hernandez's tlsping.

## License
MIT License

Copyright (c) 2023 matrachma

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
