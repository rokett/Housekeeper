# Housekeeper

Small app to delete specified files older than the specified number days from a specified directory.

## Usage

Command line flags are used to pass in parameters.

| Flag               | Description                                                              | Default Value |
| ------------------ | ------------------------------------------------------------------------ | ------------- |
| ext                | Files matching the extension will be deleted. Use * to match all files.  | none          |
| older-than         | Files which are older than the number of days specified will be deleted. | none          |
| path               | Path to search for files to be deleted. DO NOT use trailing slashes.     | none          |
| recursive          | Search directory path recursively.                                       | false         |
| test               | Carry out a test run.  No files will be deleted.                         | false         |
| debug              | Enable debug logging.                                                    | false         |
| version            | Display version and build number                                         | false         |
| case-insensitive   | Match file extensions regardless of case                                 | false         |
| remove-directories | Remove **empty** subdirectories when doing a recursive search            | false         |

For example; delete files with the **.log** file extension, ignoring case, which are older than **30** days and are anywhere within the **c:\logs** directory.  Additionally remove any empty directories within the **c:\logs** directory.

````Batchfile
housekeeper.exe --ext "log" --older-than 30 --path "c:\logs" --recursive --case-insensitive --remove-directories
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
