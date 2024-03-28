# KMST

`kmst` (statusbar) highly customizable statusbar.

See [Documentation](doc.md) for more info.

## Overview

`kmst` (statusbar) displays system information.
Information is split into segments, each representing a hardware or a software component:

* CPU average load and frequency
* Temperature

Each segment has a separate parser to get data from system files or shell programs.
Data parsing is performed every given time period (configurable by the `parseinterval` option).
Parsing can be additionally performed in reaction to system signal (configurable by the `parseonsig` option).
The order of segments as well as the templates are customizable.

## Installation

    # make buildx install

or

    # make build install

## Run

    $ kmst -x

or

    $ kmst
