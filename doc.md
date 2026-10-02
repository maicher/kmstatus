# NAME

kmstatus - dynamic status bar

# SYNOPSIS

**kmstatus**
[**--config** *path*]
[**--socketpath** *path*]
[**--xwindow**]

**kmstatus**
[**--socketpath** *path*]
**--refresh** | **--text** *text* | **--text-unset**

**kmstatus**
**--print-template** | **--doc** | **--version** | **--help**

# DESCRIPTION

kmstatus displays system information as a single line of text.
It prints the line to the standard output, or sets it as the name of the X root window, which is where dwm reads its status text from.

The line consists of segments, each representing a specific hardware or software component (CPU, memory, network, audio, etc.).
Each segment has a parser, which periodically reads data from system files or external programs, and a template, which formats the data.
Segments, their order, refresh intervals and templates are set in the configuration file.

The running instance (the main process) can be controlled by running kmstatus again with a control command, e.g. to refresh segments immediately or to display a custom text.
Control commands are sent to the main process through a unix socket.

See CONFIGURATION, SEGMENTS and TEMPLATES sections for details.
See EXAMPLES section for a dwm setup.

# QUICK REFERENCE

The following command line options are provided by kmstatus:

    -c, --config PATH       path to the configuration file
    -s, --socketpath PATH   path to the socket (default
                            /tmp/kmstatus.sock)
    -x, --xwindow           set the status as the X root window name
        --print-template    print a configuration template
        --doc               print this documentation
    -v, --version           print the version
    -h, --help              print a short help

The following control commands are provided by kmstatus:

    -r, --refresh           refresh segments and the status now
    -t, --text TEXT         display a custom text (empty clears it)
    -u, --text-unset        clear the custom text

The following options can be set for each segment in the configuration file:

    parsername         string     (required)
    refreshinterval    duration   (default "1s")
    template           string

The following parsers are available as segments:

    processes          icons of running programs
    bluetooth          bluetooth service, controller and device state
    audio              speaker and microphone volume
    cpu                CPU load and frequency
    temperature        temperature sensors
    mem                memory and swap usage
    network            network interfaces' transfer
    clock              date and time

The following template functions are provided by segments:

    bar                cpu        load as a block character
    human              mem        human-readable size (from kB)
    human              network    human-readable size (from bytes)
    ljust              network    pad a value with spaces to a width
    hasPrefix          network    test if a string starts with a prefix
    format             clock      format the time with a layout

# OPTIONS

Options can be given with one or two dashes, e.g. **-c** and **--config** are the same.

## -c, --config *path*

Path to the configuration file.
If not given, kmstatus looks up the following paths and uses the first existing file:

    $XDG_CONFIG_HOME/kmstatus/kmstatusrc.toml
    $HOME/.config/kmstatus/kmstatusrc.toml

If none exists, the built-in default configuration is used.
A given path to a file which does not exist is an error.

## -s, --socketpath *path* (default /tmp/kmstatus.sock)

Path to the unix socket used for control commands.
The main process creates the socket and control commands connect to it, so both need to be given the same path.
Running several instances requires a separate socket for each of them.

## -x, --xwindow

Set the status as the name (WM_NAME) of the X root window instead of printing it.
This is the way of setting the status text of dwm.
Requires kmstatus to be built with X support (**make buildx**).

## --print-template

Print a configuration template and exit.
The template is the default configuration, used when no configuration file is found.
It is a starting point for creating a configuration file (see CONFIGURATION).

## --doc

Print this documentation and exit.

## -v, --version

Print the version and exit.

## -h, --help

Print a short help and exit.

# CONTROL COMMANDS

A control command sends a message to the running main process and exits.
It fails if the main process is not running.
Only one control command can be given at a time.

## -r, --refresh

Parse data of all segments immediately and render the status.
This is useful when a change should be visible right away instead of after the segment's refresh interval, e.g. after changing the volume.
The cpu and network segments are not refreshed, since their values (load, transfer speed) are calculated between parses in equal time intervals.

## -t, --text *text*

Display a custom text at the beginning of the status, before all segments.
The text is surrounded with a space on each side.
An empty text clears the text.

## -u, --text-unset

Clear the custom text.

# CONFIGURATION

The configuration file is written in TOML.
It consists of a list of segments, each defined in a **[[segment]]** table.
Segments are displayed in the order of definition.
The same parser can be used by more than one segment, e.g. to display it with different templates.

    [[segment]]
    parsername = "cpu"
    refreshinterval = "1s"
    template = """{{ .Load | printf "%4.1f" }}%"""

    [[segment]]
    parsername = "clock"
    template = """ {{ . | format "15:04" }}"""

An unknown option (e.g. a misspelled one) is an error.

When no configuration file is found, the built-in default configuration is used.
It is printed by **--print-template** and can be viewed online:

    https://github.com/maicher/kmstatus/blob/master/internal/config/kmstatusrc.example.toml

## Creating a configuration file

The easiest way is to start from the configuration template and adjust it:

    mkdir -p ~/.config/kmstatus
    kmstatus --print-template > ~/.config/kmstatus/kmstatusrc.toml

Edit the file and test it in a terminal, without affecting the running status bar, by using a separate socket:

    cd ~/.config/kmstatus
    kmstatus -c kmstatusrc.toml -s /tmp/kmstatus-test.sock

The configuration is read on start, so restart the running kmstatus to apply it.

## parsername (string)

Name of the parser, one of the parsers listed in the SEGMENTS section.

## refreshinterval (duration) (default "1s")

How often the segment parses its data.
A duration is a number with a unit, e.g. "500ms", "5s", "1m", "1m30s".

The status is rendered with the shortest refresh interval of all segments, but at least once a minute.
For example, a configuration of segments refreshed every "1s" and "1m" renders the status every second, but the data of the latter segment changes once a minute.

## template (string)

Template of the segment's output.
See the TEMPLATES section for the syntax and the SEGMENTS section for the data available in each segment.
The processes segment uses its own format, see its description.

# TEMPLATES

Templates use the syntax of the Go text/template package.
Values are inserted with **{{ }}** actions, everything else is printed as it is:

    template = """CPU: {{ .Load }}%"""

Newline characters are removed from templates.
Thanks to this, a long template can be split into lines in a TOML multi-line string (**"""**), but a space must be written explicitly where it is needed:

    template = """
     {{ if .OutMuted }}muted{{ else }}{{ .OutVolume }}%{{ end }}
    """

## Data

A field or a method of the segment's data is accessed with a dot and its name, e.g. **.Load**.
A dot alone is the whole data, e.g. the time in the clock segment.

## Pipelines

A value can be passed to a function with a pipe **|**.
The passed value becomes the last argument of the function:

    {{ .Load | printf "%4.1f" }}     is    {{ printf "%4.1f" .Load }}
    {{ .Name | hasPrefix "en" }}     is    {{ hasPrefix "en" .Name }}

Hence, for a comparison it is clearer to write the arguments in order:

    {{ if gt .SwapUsed 0 }}          true if .SwapUsed > 0
    {{ if .SwapUsed | gt 0 }}        true if 0 > .SwapUsed (!)

## Built-in functions

The following functions are available in all templates:

    and, or, not              boolean operators
    eq, ne, lt, le, gt, ge    comparisons (==, !=, <, <=, >, >=)
    printf                    format values, e.g. printf "%.1f"
    len, index, slice         length, element, part of a value

## Conditions

    {{ if .IsServiceActive }}on{{ else }}off{{ end }}
    {{ if or (eq .Name "k10temp") (eq .Name "amdgpu") }}
    {{ .Value }}
    {{ end }}

# SEGMENTS

## processes

Displays an icon for each running program from a list.

The template of this segment is not a Go template.
It is a list of lines, each consisting of an icon and a process name separated by whitespace.
The icon can not contain whitespace, the process name may contain spaces.
An icon is displayed when the name of any running process starts with the given process name.
Icons are displayed in the order of lines without separators, so a separator should be a part of the icon.

    template = """
    F firefox
    V vpn
    S Socket Process
    """

Process names are read with **ps -e -o comm=**.
The kernel truncates them to 15 characters.

## bluetooth

Displays the state of the bluetooth service, controller and the connected device.

    .IsServiceActive      bool    bluetooth service is running
    .IsControllerPowered  bool    bluetooth controller is on
    .DeviceType           string  type of the connected device, e.g.
                                  "audio-headset", "audio-card", or
                                  ""

Requires **systemctl** and **bluetoothctl**.

## audio

Displays the volume of the default speaker (sink) and of the microphone.
The microphone is the audio source with "Microphone" in its name or description (see **pamixer --list-sources**).
It is not necessarily the default source.

    .OutAvailable         bool    speaker state could be read
    .OutMuted             bool    speaker is muted
    .OutVolume            int     speaker volume in percents
    .InAvailable          bool    microphone was found
    .InMuted              bool    microphone is muted
    .InVolume             int     microphone volume in percents

Requires **pamixer**.
Since the volume usually changes from a keybinding, a long refresh interval combined with the **--refresh** control command is recommended (see EXAMPLES).

## cpu

Displays the CPU load and frequency, averaged over all cores.

    .Load                 float   load in percents since the
                                  previous parse
    .Freq                 int     frequency in kHz
    .FreqMHz              float   frequency in MHz
    .FreqGHz              float   frequency in GHz

    bar LOAD                      load as a block character:
                                  " " below 2%, "░" below 9%,
                                  "▒" below 30%, "▓" below 90%,
                                  "█" otherwise

Example:

    template = """
    {{ .FreqGHz | printf "%3.1f" }}GHz {{ .Load | bar }}
    """

Data is read from **/proc/stat** and **/sys/devices/system/cpu/cpu\*/cpufreq/scaling_cur_freq**.
Not affected by the **--refresh** control command.

## temperature

Displays temperature sensors.
The template is executed once for each sensor.

    .Name                 string  name of the sensor
    .Value                int     temperature in Celsius
    .Celsius              int     temperature in Celsius
    .Fahrenheit           int     temperature in Fahrenheit

Example, to display a single sensor:

    template = """
    {{ if eq .Name "k10temp" }} {{ .Celsius }}°C{{ end }}
    """

Sensors are read from thermal zones (**/sys/devices/virtual/thermal/thermal_zone\*/temp**), named after their type, e.g. "x86_pkg_temp", "acpitz".
If there are no thermal zones, hwmon sensors are read (**/sys/class/hwmon/hwmon\*/temp1_input**), named after their driver, e.g. "k10temp", "amdgpu", "nvme".
To list the names, use the template:

    template = "{{ .Name }} "

## mem

Displays memory and swap usage.
Sizes are in kB.

    .Total                int     total memory
    .Used                 int     used memory (without buffers and
                                  cache)
    .UsedPercentage       float   used memory in percents
    .SwapTotal            int     total swap
    .SwapUsed             int     used swap
    .SwapUsedPercentage   float   used swap in percents

    human PRECISION SIZE          size in kB with a unit (M, G, T)
                                  and PRECISION decimal places, e.g.
                                  1.5G; sizes below 1M are printed
                                  as a number

Example, to display swap only when used:

    template = """
     {{ .Used | human 1 }}/{{ .Total | human 0 }}
    {{ if gt .SwapUsed 0 }} swap {{ .SwapUsed | human 1 }}{{ end }}
    """

Data is read from **/proc/meminfo**.

## network

Displays transfer of network interfaces.
The template is executed once for each interface, including the loopback (lo).

    .Name                 string  name of the interface, e.g. "eno1"
    .Rx                   int     received bytes per second
    .Tx                   int     transmitted bytes per second
    .RxTotal              int     received bytes in total
    .TxTotal              int     transmitted bytes in total

    human PRECISION SIZE          size in bytes with a unit (k, M,
                                  G, T) and PRECISION decimal
                                  places, e.g. 2k; sizes below 1k
                                  are printed as a number
    ljust WIDTH VALUE             VALUE (a string or an int)
                                  preceded by spaces to WIDTH
                                  characters, for a constant width
    hasPrefix PREFIX STRING       STRING starts with PREFIX

Example, to display wired interfaces only:

    template = """
    {{ if .Name | hasPrefix "en" }}
     {{ .Name }}
     {{ .Rx | human 0 | ljust 4 }}
     {{ .Tx | human 0 | ljust 4 }}
    {{ end }}
    """

Speeds are averaged since the previous parse.
Data is read from **/proc/net/dev**.
Not affected by the **--refresh** control command.

## clock

Displays the current date and time.
The data (dot) is the current time, read on each render.

    format LAYOUT TIME            TIME formatted according to LAYOUT

The layout is the reference time **Mon Jan 2 15:04:05 MST 2006** written in the desired format:

    {{ . | format "2006-01-02 15:04:05" }}      2024-05-01 08:10:19
    {{ . | format "Mon 02 Jan 15:04" }}         Wed 01 May 08:10

Methods of the time can be used as well, e.g. **{{ .Weekday }}**, **{{ .Hour }}**.

The refresh interval of the clock only affects how often the status is rendered.

# EXAMPLES

Start kmstatus as the status bar of dwm, e.g. in **~/.xinitrc**:

    kmstatus -x &
    exec dwm

Change the volume and update the status immediately, e.g. in a dwm keybinding:

    pamixer -i 5 && kmstatus -r

In dwm's config.h (with **#include <X11/XF86keysym.h>**):

    static const char *volup[] = {
        "sh", "-c", "pamixer -i 5 && kmstatus -r", NULL
    };
    { 0, XF86XK_AudioRaiseVolume, spawn, {.v = volup } },

Display a notice while a long task runs:

    kmstatus -t "backup" && rsync -a ~/ /mnt/backup/; kmstatus -u

Run a second instance in a terminal, with another configuration and socket:

    kmstatus -c ~/.config/kmstatus/term.toml -s /tmp/kmstatus-term.sock

Test a configuration, without affecting the running status bar:

    kmstatus -c ./kmstatusrc.toml -s /tmp/kmstatus-test.sock

# FILES

    $XDG_CONFIG_HOME/kmstatus/kmstatusrc.toml     configuration file
    $HOME/.config/kmstatus/kmstatusrc.toml        configuration file
    /tmp/kmstatus.sock                            default socket

The default configuration is built into kmstatus and printed by **--print-template**.
Its source is in the repository:

    https://github.com/maicher/kmstatus/blob/master/internal/config/kmstatusrc.example.toml

# ENVIRONMENT

## XDG_CONFIG_HOME

Directory of the configuration file, see **--config**.

## HOME

Directory of the configuration file, if not found in XDG_CONFIG_HOME.

## DISPLAY

X display to set the status on, with **--xwindow**.

# DIAGNOSTICS

The main process exits with status 1 on an error, e.g. an invalid configuration, or when another instance already uses the socket.
A socket left by an instance that did not exit cleanly is removed on start.

Errors which occur while running, e.g. a sensor that can not be read, are printed to the standard error and do not stop kmstatus.
An external program (pamixer, bluetoothctl, ps) is stopped if it does not finish in 2 seconds.

A control command exits with status 1 if the main process is not running, and with status 2 on invalid options.

# BUGS

kmstatus runs on Linux only.
It was tested on Arch Linux and Ubuntu, with Intel and AMD CPUs.

Report bugs at https://github.com/maicher/kmstatus/issues.

# SEE ALSO

**dwm**(1), **pamixer**(1), **bluetoothctl**(1), **ps**(1), **proc**(5)
