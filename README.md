# SSP Server

```sh
███████ ███████ ██████  ███████ ███████ ██████  ██    ██ ███████ ██████
██      ██      ██   ██ ██      ██      ██   ██ ██    ██ ██      ██   ██
███████ ███████ ██████  ███████ █████   ██████  ██    ██ █████   ██████
     ██      ██ ██           ██ ██      ██   ██  ██  ██  ██      ██   ██
███████ ███████ ██      ███████ ███████ ██   ██   ████   ███████ ██   ██
```

[LICENSE](LICENSE.md)

[![Build Status](https://github.com/sspserver/sspserver/workflows/Tests/badge.svg)](https://github.com/sspserver/sspserver/actions?workflow=Tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/sspserver/sspserver)](https://goreportcard.com/report/github.com/sspserver/sspserver)
[![GoDoc](https://godoc.org/github.com/sspserver/sspserver?status.svg)](https://godoc.org/github.com/sspserver/sspserver)
[![Coverage Status](https://coveralls.io/repos/github/sspserver/sspserver/badge.svg)](https://coveralls.io/github/sspserver/sspserver)

> **Attention:** Watch for us, the service in active development.

SSP (Supply-Side Platform) Advertisement Service is a comprehensive platform for monetizing digital properties like websites, mobile apps, or other software. It integrates seamlessly with various demand sources through RTB (Real-Time Bidding) and direct ad campaigns. With built-in event tracking, flexible ad-source management, and analytics, SSP helps developers, content creators, and businesses maximize revenue with minimal effort.

## Features

- **SSP Integration**: Supports RTB and direct ad campaigns with multiple ad sources.
- **Ad Storage and Management**: Manage ad campaigns, creatives, and formats in a scalable way.
- **Event Tracking**: Log clicks, impressions, wins, and other advertising events for detailed analytics.
- **Pixel Tracking**: Integrated pixel tracking for user behavior and ad performance monitoring.
- **Flexible Configuration**: Easy setup for different ad sources (RTB, direct, etc.) and user-targeting rules.
- **High-Performance HTTP Server**: Built with `fasthttp` for scalable and high-performance ad serving.
- **Personification & Targeting**: Custom user detection for advanced audience targeting.
- **Extensible Architecture**: Add features easily through extensions for tracking, DSP integration, etc.

## Motivation

The SSP Advertisement Service is designed to streamline monetization for developers, content creators, and businesses by simplifying the integration of ads into digital properties. The platform supports real-time bidding (RTB) and direct ad campaigns, helping you maximize your revenue through intelligent ad placements, advanced targeting, and performance tracking. With a robust set of analytics tools, you can optimize your ad strategies for better revenue generation. Whether you want to monetize a website, mobile app, or software, SSP Server provides all the necessary tools to do so efficiently.

## Getting Started

1. Clone the repository:

```bash
git clone https://github.com/sspserver/sspserver.git
cd sspserver
```

2. Install dependencies:

```bash
go mod tidy
```

3. Set up configuration in config.yaml.
4. Run the server:

```bash
go run main.go
```

For more detailed information on how to configure and run SSP Server, check out the Quick Start Guide.

## Key Components

- Ad Source Management: Manage ads from different sources, including RTB and direct campaigns.
- Event Stream Pipeline: Handles ad performance events like clicks, impressions, and wins.
- HTTP Server: Provides endpoints for external ad requests, utilizing fasthttp for high performance.
- Extensions: Extend the server with features such as pixel tracking, DSP integration, and analytics.

## Documentation

- [Quick Start](docs/quick-start.md)
- [Integration](docs/integration.md)

> Full documentation is a work in progress.

## Roadmap

- [ ] Complete documentation for configuration and usage.
- [ ] Switch template engine from quicktemplate to templ.
- [ ] Implement monetary usage interfaces (currently NoOp).
- [ ] Improve ad-selector optimizer interfaces.
- [ ] Extend the network layer client implementation.

## Contributing

We welcome contributions! Please see the CONTRIBUTING.md for more information on how to get started.

## License

SSP Server is licensed under the [GNU Affero General Public License](LICENSE.md).
