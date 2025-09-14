# Mintop

Mintop is a simple, terminal-based process monitor built with Go and the Bubble Tea framework. It provides a real-time view of system processes and resource usage.

## Features

- **Process Monitoring**: View a list of running processes with details such as PID, PPID, name, CPU usage, memory usage, and owner.
- **System Information**: Displays key system metrics including:
    - CPU Usage
    - Memory and Swap Usage
    - System Load Average
- **Interactive Table**: The process list is displayed in an interactive table that can be scrolled.

## Installation

To build and run Mintop, you need to have Go installed on your system.

1. **Clone the repository:**
   ```bash
   git clone https://github.com/ashwineaso/mintop.git
   cd mintop
   ```

2. **Build the application:**
   ```bash
   go build
   ```

3. **Run Mintop:**
   ```bash
   ./mintop
   ```

## How it Works

Mintop uses the following libraries to gather system information and build the terminal UI:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea): A powerful framework for building terminal-based applications.
- [gopsutil](https://github.com/shirou/gopsutil): A cross-platform library for retrieving process and system utilization information.
- [Lipgloss](https://github.com/charmbracelet/lipgloss): A library for styling terminal output.

The application is structured around the Model-View-Update architecture provided by Bubble Tea. The `internal` package contains the core logic for fetching data, updating the model, and rendering the view.

## Contributing

Contributions are welcome! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines on how to contribute to the project.
