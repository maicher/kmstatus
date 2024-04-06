# kmstatus

See [Documentation](doc.md) for more info.

## Overview

`kmstatus` displays system information.
Information is split into segments, each representing a hardware or a software component:

* CPU average load and frequency
* Temperature

The order of segments as well as the templates are customizable via toml config file.

Each segment has a separate parser to get data from system files or shell programs.
Once `kmstatus` has started it parses data every given time interval (configurable by the `refreshinterval` option)
and prints output periodically.

The process can be additionally controlled by following commands:

    kmstatus --refreh
    kmstatus --text TEXT
    kmstatus --text-unset

## Installation

    # make buildx install

or

    # make build install

## Run

    $ kmstatus -x

or

    $ kmstatus
