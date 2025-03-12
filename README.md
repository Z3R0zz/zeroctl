# zeroctl

`zeroctl` is a custom CLI tool designed to manage various tasks such as changing wallpapers, fetching weather data, and more. It includes a daemon mode to run scheduled tasks and listen for commands.

## Features

- **Wallpaper Management**: Randomize and set wallpapers.
- **Weather Data**: Fetch and cache weather data from OpenWeather API.
- **Task Scheduling**: Schedule and run periodic tasks.
- **System Stats**: Display system statistics and uptime.

## Installation

### Dependencies

Before building `zeroctl`, ensure you have the following dependencies installed on Arch Linux:

```sh
sudo pacman -S go git swww
yay -S python-pywal16
```

1. Clone the repository:
    ```sh
    git clone https://github.com/Z3R0zz/zeroctl.git
    cd zeroctl
    ```

2. Copy the example environment file and configure it:
    ```sh
    cp .env.example .env
    # Edit .env to set your environment variables
    ```

3. Build and install `zeroctl`:
    ```sh
    ./build.sh
    ```

4. Add the following lines to your hyprland.conf:
    ```sh
    exec-once = zeroctl daemon | tee -a /var/log/zeroctl.log
    exec-once = zeroctl wallpaper | tee -a /var/log/zeroctl.log # Optional
    ```

## Usage

### Running the Daemon

To start the `zeroctl` daemon, run:
```sh
zeroctl daemon
```

### Commands

- **wallpaper**: Randomize wallpaper
    ```sh
    zeroctl wallpaper
    ```

- **weather**: Get the current weather
    ```sh
    zeroctl weather
    ```

- **stats**: Get stats about zeroctl and its usage
    ```sh
    zeroctl stats
    ```

- **uptime**: Display how long the zeroctl daemon has been running
    ```sh
    zeroctl uptime
    ```

## Environment Variables

- `WALLPAPERS_DIR`: Directory containing wallpapers.
- `AFTER_WALLPAPER_COMMANDS`: Commands to run after setting the wallpaper.
- `OPENWEATHER_API_KEY`: API key for OpenWeather.
- `OPENWEATHER_CITY_ID`: City ID for OpenWeather.
- `OPENWEATHER_UNITS`: Units for weather data (e.g., metric).

## License

This project is licensed under the MIT License. Refer to the [LICENSE](LICENSE) file for more information.
