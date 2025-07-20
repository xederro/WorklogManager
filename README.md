# WorklogManager

[![Go Version](https://img.shields.io/github/go-mod/go-version/xederro/WorklogManager)](https://golang.org/)
[![License](https://img.shields.io/github/license/xederro/WorklogManager)](LICENSE)
[![Build Status](https://img.shields.io/github/downloads/xederro/WorklogManager/total)](https://github.com/xederro/WorklogManager/releases)
[![Release](https://img.shields.io/github/v/release/xederro/WorklogManager)](https://github.com/xederro/WorklogManager/releases)

## Overview

**WorklogManager** is a Terminal User Interface (TUI) application built with Go and [Charm](https://charm.sh/), designed to simplify the process of measuring and logging time spent on JIRA tasks. This tool integrates with JIRA through the REST API (v2), allowing users to seamlessly track the time spent on tasks and log it directly into JIRA.

## Motivation

As a developer who struggled with accurately logging the time spent on various tasks in JIRA, I developed WorklogManager to address this challenge. The goal of this project is to provide a lightweight, easy-to-use solution that helps track work time in a more efficient manner.

## Features

- **Track Time Spent on Tasks:** Start, pause, and stop timers for specific JIRA tasks.
- **Log Work Directly to JIRA:** Log the tracked time to the appropriate JIRA issues using the REST API (v2).
- **Create Worklogs with Google AI:** Use Google AI to better prepare your worklog.
- **TUI Interface:** Navigate and interact with the app entirely from the terminal, using a simple and intuitive interface.

## Building

To build WorklogManager, you need to have Go >= 1.24, pkl and sqlc installed on your machine. You can then clone the repository and build the application:

```bash
git clone https://github.com/xederro/WorklogManager.git
cd WorklogManager
go generate ./generate.go
go build
```

## Installation
You can download the latest release of WorklogManager from the [Releases page](https://github.com/xederro/WorklogManager/releases). 
The release includes pre-built binaries for both Linux and Windows.

## Usage

Run the built executable from the terminal:

### Linux
```bash
./WorklogManager [-config <path_to_config_file>]
```

### Windows
```bash
./WorklogManager.exe [-config <path_to_config_file>]
```

### Configuration
WorklogManager requires a configuration file to connect to your JIRA instance. You can create a configuration file in PKL format. Below is an example of how to structure your configuration file:

```pkl
amends "package://github.com/xederro/WorklogManager/releases/download/0.2.2/WorklogManager@0.2.2#/Config.pkl"

db_path = "db/db.sqlite"
jira {
  url = "https://your-jira-instance.atlassian.net"
  default_worklog_comment = "Work"
  server_type = "cloud"
  cloud_config {
    email = "email@example.com"
    api_token = "your_api_token"
  }
  on_premise_config = null
  requests {
    new {
      jql = "assignee = currentUser() AND statusCategory != Done ORDER BY updated DESC"
      refetch_interval = 10.min
    }
    new {
      jql = "assignee = currentUser() ORDER BY updated DESC"
      refetch_interval = 15.min
    }
  }
}
use_ai = true
google_ai {
  APIKey = "your_google_api_key"
  default_model = "gemini-2.5-flash"
  default_prompt = "You are Senior Jira Worklog Writer. Write a short worklog description from the following information:"
}
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
