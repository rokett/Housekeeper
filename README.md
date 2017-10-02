# Housekeeper
Small app to delete specified files older than the specified number days from a specified directory.

## Usage
Command line flags are used to pass in parameters.

| Flag       | Description                                                              | Default Value |
| ---------- | ------------------------------------------------------------------------ | ------------- |
| ext        | Files matching the extension will be deleted.                            | none          |
| older-than | Files which are older than the number of days specified will be deleted. | none          |
| path       | Path to search for files to be deleted.                                  | none          |
| recursive  | Search directory path recursively.                                       | false         |
| test       | Carry out a test run.  No files will be deleted.                         | false         |
| debug      | Enable debug logging.                                                    | false         |

For example; delete files with the **.log** file extension which are older than **30** days and are anywhere within the **c:\logs** directory.

````
housekeeper.exe -ext "log" -older-than 30 -path "c:\logs" -recursive
````

## Building the executable
All dependencies are version controlled, so building the project is really easy.

1. Clone the repository locally.
2. From within the repository directory run `go build`.
3. Hey presto, you have an executable.
