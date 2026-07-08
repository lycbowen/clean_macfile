# maclean

`maclean` removes macOS metadata files from folders you use on Windows, Linux, USB drives, shared drives, archives, or anywhere else those extra files get in the way.

It looks for:

- `.DS_Store`
- AppleDouble `._*` files that are 4096 bytes

Matching files are first moved into a `.wait_clean` folder. After that, `maclean` asks whether you want to delete the archived files, so you get one last chance to keep them.

## Install

If you have Go installed:

```sh
go install github.com/hicbowen/clean_macfile/cmd/maclean@latest
```

Make sure your Go binary directory is in `PATH`. It is usually:

- macOS/Linux: `~/go/bin`
- Windows: `%USERPROFILE%\go\bin`

You can also download a release binary from the GitHub Releases page if you do not want to install Go.

## Usage

Clean the current folder:

```sh
maclean
```

Clean a specific folder:

```sh
maclean -t /path/to/folder
```

On Windows:

```powershell
maclean -t D:\SharedFolder
```

Check the installed version:

```sh
maclean -version
```

## What Happens

When matching files are found, `maclean` moves them into `.wait_clean` inside the target folder and then prompts:

```text
Delete archived files? (Y/n):
```

Press `Enter` or type `y` to delete them. Type `n` to keep them in `.wait_clean` for review.

If no matching files are found, `maclean` prints a short message and leaves the folder unchanged.

## Notes

- `.wait_clean` is skipped during scanning.
- Existing files in `.wait_clean` are not overwritten; duplicate names get a numeric suffix.
- `maclean` only targets macOS metadata files. It does not remove general hidden files.
