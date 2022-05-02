# Housekeeper

Small app to delete specified files older than the specified number days from a specified directory.

## Usage

Command line flags are used to pass in parameters.

| Flag               | Description                                                                             | Default Value | Required? |
| ------------------ | --------------------------------------------------------------------------------------- | ------------- | --------- |
| ext                | Files matching the extension will be deleted. Use * to match all files.                 | none          | yes       |
| older-than         | Files which are older than the number of <older-than-units> specified will be deleted.  | none          | yes       |
| older-than-units   | Specifies the time units to use; d(ays), (h)ours, or (m)inutes.                         | d             | no        |
| path               | Path to search for files to be deleted. DO NOT use trailing slashes.                    | none          | yes       |
| recursive          | Search directory path recursively.                                                      | false         | no        |
| test               | Carry out a test run.  No files will be deleted.                                        | false         | no        |
| debug              | Enable debug logging.                                                                   | false         | no        |
| version            | Display version and build number                                                        | false         | no        |
| case-insensitive   | Match file extensions regardless of case                                                | false         | no        |
| remove-directories | Remove **empty** subdirectories, those without files in, when doing a recursive search, **if** the directory is older than the `--older-than` flag. | false         | no        |

For example; delete files with the **.log** file extension, ignoring case, which are older than **30 hours** and are anywhere within the **c:\logs** directory.  Additionally remove any empty directories, those not containing files, within the **c:\logs** directory if they are older than **30 hours**.

````Batchfile
housekeeper.exe --ext "log" --older-than 30 --older-than-units h --path "c:\logs" --recursive --case-insensitive --remove-directories
````

**NOTE:** Do not use trailing slashes in the path.  On Windows this causes the final `"` to be escaped resulting in Housekeeper thinking that the rest of the command is all part of the `path` flag.  This isn't a Housekeeper limitation technically, it's how Windows works.

Logs are printed to the Windows Event Log (if on Windows) and stdout which means you can do with them as you wish.  A good option is to pipe the output to another application which sends the logs where you want.

Logging to the Windows Event Log will require permissions to create an event source.

## Downloading a release

<https://github.com/rokett/Housekeeper/releases>

## Building the executable

All dependencies are version controlled, so building the project is really easy.

1. `go get github.com/rokett/housekeeper`.
2. From within the repository directory run `make.bat`.
3. Hey presto, you have an executable.
