# Housekeeper

Small app to delete specified files older than the specified number days from a specified directory.

## Usage

Command line flags are used to pass in parameters.

| Flag       | Description                                                              | Default Value |
| ---------- | ------------------------------------------------------------------------ | ------------- |
| ext        | Files matching the extension will be deleted. Use * to match all files.  | none          |
| older-than | Files which are older than the number of days specified will be deleted. | none          |
| path       | Path to search for files to be deleted.                                  | none          |
| recursive  | Search directory path recursively.                                       | false         |
| test       | Carry out a test run.  No files will be deleted.                         | false         |
| debug      | Enable debug logging.                                                    | false         |
| version    | Display version and build number                                         | false         |

For example; delete files with the **.log** file extension which are older than **30** days and are anywhere within the **c:\logs** directory.

````Batchfile
housekeeper.exe -ext "log" -older-than 30 -path "c:\logs" -recursive
````

Logs are printed to stdout which means you can do with them as you wish.  A good option is to pipe the output to another application which sends the logs where you want.

## Downloading a release

<https://github.com/rokett/Housekeeper/releases>

## Building the executable

All dependencies are version controlled, so building the project is really easy.

1. Clone the repository locally.
2. From within the repository directory run `make.bat build`.
3. Hey presto, you have an executable.
