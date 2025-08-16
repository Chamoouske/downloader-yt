# Downloader-YT

Downloader-YT is a powerful and flexible YouTube video downloader written in Go. It provides multiple interfaces to download your favorite videos, including a command-line interface (CLI), a dedicated Termux interface for Android users, and a web server for remote downloads.

## Features

- **Multiple Interfaces**: Choose the interface that best suits your needs:
  - **CLI**: A simple and straightforward command-line tool for quick downloads.
  - **Termux**: A specialized CLI for Termux on Android, with support for native notifications.
  - **Web Server**: A robust web server that exposes an API for remote video downloads.
- **YouTube Integration**: Utilizes the `kkdai/youtube/v2` library to efficiently fetch video information and download streams.
- **Notifications**: Stay informed about your downloads with built-in notification support:
  - **Termux Notifier**: Receive native Android notifications when using the Termux interface.
  - **Server Notifier**: Get webhook notifications sent to a configured URL.
- **In-Memory Database**: Stores video metadata in a fast and efficient in-memory database.
- **Progress Bar**: Keep track of your downloads with a real-time progress bar in the terminal.
- **Flexible Configuration**: Easily configure the application using a `config.json` file, with support for environment variables.
- **Comprehensive Logging**: Logs are sent to both the console and a file for easy debugging and monitoring.
- **Clean Architecture**: The project follows the principles of Clean Architecture, ensuring a modular, scalable, and maintainable codebase.
- **Docker Support**: Run the application in a containerized environment with the included `Dockerfile` and `docker-compose.yml` files.

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.18 or later)
- [Docker](https://docs.docker.com/get-docker/) (optional, for running in a container)

### Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/your-username/downloader-yt.git
   cd downloader-yt
   ```

2. **Install dependencies:**
   ```sh
   go mod tidy
   ```

3. **Build the binaries:**
   ```sh
   make build
   ```
   This will create the binaries for the CLI, Termux, and web server in the `bin/` directory.

## Usage

### CLI

To download a video using the CLI, use the following command:

```sh
./bin/cli -v "YOUTUBE_VIDEO_URL"
```

### Termux

For Termux users, the command is the same:

```sh
./bin/termux -v "YOUTUBE_VIDEO_URL"
```

### Web Server

To start the web server, run:

```sh
./bin/web -p 8080
```

You can then download a video by sending a GET request to the `/video/download` endpoint:

```
http://localhost:8080/video/download?url=YOUTUBE_VIDEO_URL&requester=USER_ID
```

The downloaded video can be accessed at:

```
http://localhost:8080/video/VIDEO_ID
```

## Configuration

The application can be configured using a `config.json` file located in the `.config/` directory. You can also use environment variables to override the default settings.

| Setting      | Environment Variable | Default Value                                            | Description                               |
|--------------|----------------------|----------------------------------------------------------|-------------------------------------------|
| `port`       | `PORT`               | `8080`                                                   | The port for the web server to listen on. |
| `log_dir`    | `LOG_DIR`            | `./.logs`                                                | The directory to store log files.         |
| `video_dir`  | `VIDEO_DIR`          | `./videos`                                               | The directory to store downloaded videos. |
| `config_dir` | `CONFIG_DIR`         | `./.config`                                              | The directory for the configuration file. |
| `url_webhook`| `WEBHOOK`            | `http://host.docker.internal:5677/webhook/downloader-yt` | The URL for webhook notifications.        |

## Docker

To run the web server in a Docker container, use the following commands:

1. **Build the Docker image:**
   ```sh
   docker-compose build
   ```

2. **Start the container:**
   ```sh
   docker-compose up -d
   ```

The web server will be accessible at `http://localhost:8080`.

## Project Structure

The project is organized into the following directories:

- `bin/`: Contains the compiled binaries.
- `cmd/`: Contains the `main` packages for the different entry points (CLI, Termux, and web server).
- `internal/`: Contains the core application logic, divided into `domain`, `infra`, and `usecase` layers.
- `pkg/`: Contains reusable packages for configuration, logging, and utilities.
- `videos/`: The default directory for downloaded videos.
- `.logs/`: The default directory for log files.
- `.config/`: The default directory for the configuration file.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue to discuss your ideas.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
