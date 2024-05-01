# kmstatus

See [Documentation](doc.md) for more info.

## Overview

`kmstatus` displays system information.

Data is organized by a distinct segments, each representing a specific hardware or software component:

             󰂯 6%78% 0.8GHz  0.4% ░ 37°   6G/62G(9.6%) [eno1   2k 462] 2024-05-01 08:10:19
    -------------------------------------------------------------------------------------------
    processes bt audio   cpu            temp  memory        network         time

* Processes
* Bluetooth
* Audio
* CPU (average load and frequency)
* Temperature
* Memory
* Network
* Clock
* and a segment for a custom text

The order of segments as well as the templates are customizable via the [kmstatusrc.toml](internal/config/kmstatusrc.example.toml)

## How to run

Pull the repository:

    $ git clone https://github.io/maicher/kmstatus
    $ cd kmstatus

`kmstatus` can be run in a terminal or as a status bar for [DWM](https://github.com/maicher/dwm).

Install and run in terminal:

    # make build install
    $ kmstatus

Install and run as a status bar for [DWM](https://github.com/maicher/dwm):

    # make buildx install
    $ kmstatus -x

Uninstall:

    # make uninstall

## How does it work

Each segment has a separate parser to get data from system files or shell programs.
`kmstatus` parses data every given time interval and prints output periodically.

The running `kmstatus` can be additionally controlled by following commands:

Trigger an additional refresh:

    kmstatus --refresh

Set and unset text in the text segment:

    kmstatus --text TEXT
    kmstatus --text-unset


## Warning

`kmstatus` was tested only on Arch and Ubuntu Linux. With both Intel and AMD CPUs.

Feel free to contribute.
