# clean_macfile

`clean_macfile` is a small command-line tool for cleaning macOS metadata files from a directory tree.

It is useful when files have been copied from macOS to Windows, Linux, USB drives, archives, or shared folders and leave behind files such as `.DS_Store` and AppleDouble `._*` metadata files.

## What It Cleans

The tool recursively scans the target directory and moves the following files into a temporary archive directory named `.wait_clean`:

- `.DS_Store`
- `._*` files with a size of `4096` bytes

The `.wait_clean` directory is skipped during scanning, so previously archived files are not processed again.

## Usage

Clean the current directory:

```sh
clean_macfile
```

Clean a specific directory:

```sh
clean_macfile -t /path/to/target
```

After matching files are moved into `.wait_clean`, the tool asks whether to delete the archived files:

```text
Delete archived files? (Y/n):
```

Press `Enter` or type `y` to delete them. Type `n` to keep them in `.wait_clean`.

## Build From Source

Requires Go 1.23 or later.

```sh
go build -o clean_macfile
```

On Windows:

```sh
go build -o clean_macfile.exe
```

## Examples

Clean the current folder and delete archived files:

```sh
clean_macfile
# Press Enter when prompted
```

Clean a mounted drive and keep the archive for review:

```sh
clean_macfile -t /Volumes/USB
# Type n when prompted
```

Clean a Windows folder:

```powershell
.\clean_macfile.exe -t D:\SharedFolder
```

## Release Builds

The GitHub Actions workflow builds release binaries for:

- Windows amd64
- Linux amd64
- macOS arm64

Release builds are generated when pushing a tag that starts with `v`, such as:

```sh
git tag v1.0.0
git push origin v1.0.0
```

## Notes

- The tool moves files before deleting them, giving you a chance to review the archive.
- If a file with the same name already exists in `.wait_clean`, a numeric suffix is added to avoid overwriting it.
- The tool is intended for macOS metadata cleanup only; it does not remove general hidden files.
