# WorklogManager

[![Go Version](https://img.shields.io/github/go-mod/go-version/xederro/WorklogManager)](https://golang.org/)
[![License](https://img.shields.io/github/license/xederro/WorklogManager)](LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/xederro/WorklogManager/build.yml)](https://github.com/xederro/WorklogManager/actions)
[![Release](https://img.shields.io/github/v/release/xederro/WorklogManager)](https://github.com/xederro/WorklogManager/releases)

## Overview

**WorklogManager** is a Terminal User Interface (TUI) application built with Go and [Charm](https://charm.sh/), designed to simplify the process of measuring and logging time spent on JIRA tasks. This tool integrates with JIRA through the REST API (v2), allowing users to seamlessly track the time spent on tasks and log it directly into JIRA.

## Motivation

As a developer who struggled with accurately logging the time spent on various tasks in JIRA, I developed WorklogManager to address this challenge. The goal of this project is to provide a lightweight, easy-to-use solution that helps track work time in a more efficient manner.

## Features

- **Track Time Spent on Tasks:** Easily start, pause, and stop timers for specific JIRA tasks.
- **Log Work Directly to JIRA:** Log the tracked time to the appropriate JIRA issues using the REST API (v2).
- **TUI Interface:** Navigate and interact with the app entirely from the terminal, using a simple and intuitive interface.

## Installation

To install WorklogManager, you need to have Go installed on your machine. You can then clone the repository and build the application:

```bash
git clone https://github.com/xederro/WorklogManager.git
cd WorklogManager
go build
```

## Usage

Run the built executable from the terminal:

```bash
JIRA_URL="https://<jira_server>/rest/api/2/" && ./WorklogManager
```

## Demo

*A GIF demonstrating the main functionalities of WorklogManager will be placed here.*

## Contributing

Contributions are welcome! If you have suggestions for new features or improvements, feel free to submit an issue or pull request.

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/YourFeature`).
3. Commit your changes (`git commit -am 'Add some feature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


---

Thank you for using WorklogManager! Happy time tracking!