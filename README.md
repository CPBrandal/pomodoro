# Pomodoro Timer

A simple, command-line Pomodoro timer built with Go. Focus on your work with customizable work/break intervals and preset management.

## Installation

Install with a single command:

curl -fsSL https://raw.githubusercontent.com/CPBrandal/pomodoro/main/install.sh | bash

Or download and run the install script:

curl -fsSL https://raw.githubusercontent.com/CPBrandal/pomodoro/main/install.sh -o install.sh
chmod +x install.sh
./install.sh

## Usage

### Start the Timer

Simply run:

pomodoro

### Features

- **Default Timer**: Quick start with 25-minute work sessions and 5-minute breaks
- **Custom Timers**: Create your own work/break intervals
- **Presets**: Save and reuse your favorite timer configurations
- **Last Used**: Quickly restart with your last used timer settings
- **Persistent Storage**: Your presets are saved and persist across sessions

### Command-Line Options

pomodoro # Start the timer
pomodoro --help # Show help message
pomodoro -h # Short form for help
pomodoro --uninstall # Uninstall the program and remove config
pomodoro -u # Short form for uninstall## Uninstallation

To uninstall the program and remove all configuration:

pomodoro -u

This will:

- Remove the configuration directory (`~/.pomodoro`)
- Provide instructions for removing the binary (requires sudo)

## Configuration

Your presets and settings are stored in `~/.pomodoro/presets.json`. You can edit this file directly if needed.

## Requirements

- macOS or Linux
- curl or wget (for installation)

## License

MIT

## Resources
ASCII art is from: https://asciiart.website/cat.php?category_id=88
