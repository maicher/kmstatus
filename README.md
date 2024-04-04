# KMST

`kmst` (statusbar) highly customizable statusbar.

See [Documentation](doc.md) for more info.

## Overview

`kmst` (statusbar) displays system information.
Information is split into segments, each representing a hardware or a software component:

* CPU average load and frequency
* Temperature

The order of segments as well as the templates are customizable via toml config file.

Each segment has a separate parser to get data from system files or shell programs.
Once `kmst` has started it parses data every given time interval (configurable by the `refreshinterval` option)
and prints output periodically.

The process can be additionally controlled by following commands:

    kmst --refreh
    kmst --text TEXT
    kmst --text-unset

## Installation

    # make buildx install

or

    # make build install

## Run

    $ kmst -x

or

    $ kmst
