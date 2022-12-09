# kmstatus

## How to run

	kmstatus 2> error.log

Booting phase is when parsers are getting initialized.
After this, parsers periodically parse system files to read and display the data.

If any errors are raised during booting phase, the program will crash displaying
the errors and will not go into the next phase.
If any errors are raised during the run phase, they are logged into the stdout
and parsing is continued.
