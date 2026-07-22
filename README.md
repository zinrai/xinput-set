# xinput-set

A declarative configuration tool for xinput devices. Define your input device settings in YAML and apply them automatically.

## Features

- Automatic device detection
- Declarative YAML configuration
- Dry-run mode for testing
- Clear command output showing what's being executed
- Stops on error (device not found, validation failed)
- Support for any xinput subcommand

## Requirements

- `xinput` command-line tool
- X Window System

## Usage

Apply configuration from default config.yaml

```bash
$ xinput-set
```

Specify a configuration file

```bash
$ xinput-set -config examples/thinkpad-usb.yaml
```

Dry-run mode (show commands without executing)

```bash
$ xinput-set -config examples/thinkpad-usb.yaml -dry-run
```

Show version information

```bash
$ xinput-set -version
```

### Configuration

See [examples/thinkpad-usb.yaml](examples/thinkpad-usb.yaml) for a complete example.

The configuration uses a simple structure:

- **profiles**: Named device configurations
- **detection**: How to find and validate devices
- **actions**: List of xinput commands to execute

For `set-prop`, the `args` field is split by taking everything before the first numeric token as the property name, so a property name containing spaces (for example `libinput Accel Speed`) is kept as one argument.

## Example output

### Normal execution

```bash
$ xinput-set -config examples/thinkpad-usb.yaml
Processing profile: thinkpad_trackpoint
Description: Lenovo ThinkPad USB Keyboard with TrackPoint
Executing: xinput list
Found 2 device(s) matching "Lenovo ThinkPad Compact USB Keyboard with TrackPoint"
Executing: xinput list-props 13
  Device id=13: missing required properties [libinput Accel Speed]
Executing: xinput list-props 15
  Device id=15: validation passed
Selected device: id=15
Executing: xinput set-prop 15 "libinput Accel Speed" 0.5
Executing: xinput set-prop 15 "Coordinate Transformation Matrix" 1 0 0 0 1 0 0 0 0.5
Executing: xinput set-button-map 15 1 0 3 4 5 6 7
Configuration applied successfully.
```

### Dry-run mode

```bash
$ xinput-set -config examples/thinkpad-usb.yaml -dry-run
Processing profile: thinkpad_trackpoint
Description: Lenovo ThinkPad USB Keyboard with TrackPoint
Executing: xinput list
Found 2 device(s) matching "Lenovo ThinkPad Compact USB Keyboard with TrackPoint"
Executing: xinput list-props 13
  Device id=13: missing required properties [libinput Accel Speed]
Executing: xinput list-props 15
  Device id=15: validation passed
Selected device: id=15
[DRY-RUN] Would execute: xinput set-prop 15 "libinput Accel Speed" 0.5
[DRY-RUN] Would execute: xinput set-prop 15 "Coordinate Transformation Matrix" 1 0 0 0 1 0 0 0 0.5
[DRY-RUN] Would execute: xinput set-button-map 15 1 0 3 4 5 6 7
```

### Error handling

```bash
$ ./xinput-set -config examples/thinkpad-usb.yaml -dry-run
Processing profile: thinkpad_trackpoint
Description: Lenovo ThinkPad USB Keyboard with TrackPoint
Executing: xinput list
Error: no valid device found matching "Lenovo ThinkPad Compact USB Keyboard with TrackPoint"
```

## Finding device names and properties

To find your device name:

```bash
$ xinput list
```

To find available properties for a device:

```bash
$ xinput list-props [device-id]
```

## License

This project is licensed under the [MIT License](./LICENSE).
