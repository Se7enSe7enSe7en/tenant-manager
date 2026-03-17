## Development Setup guide:

### Step 1: install task

```sh
go install github.com/go-task/task/v3/cmd/task@latest
```

### Step 2: use task to trigger the "setup" command

```sh
task setup
```

- this runs a command to install all the tools we will need, check Taskfile.yml to see the tools we will use

### Waahoooo your done, to run the project:

```sh
task
```

- this runs the default task command which is `task start` alternatively you can run this to start the project

## WSL2 setup guide (for Windows/WSL2 users):

- additional setup for wsl users

### Step 1: Make sure mirrored network is set to true

- note: before this make sure you are using wsl 2 version 2.0.4 or later

1. Create or edit the .wslconfig file in your user directory (C:\Users\{Username}\.wslconfig).
2. Add the following configuration under the [wsl2] section:

```
[wsl2]
networkingMode=mirrored
```

### Step 2: Install wslview/wslu

```sh
# installation
sudo apt update
sudo apt install wslu

# for checking if you installed successfully
which wslview
```

- this is so that templ can successfully open the proxy port on your default browser in Windows
- normally wsl does not have a browser, so you can use wslview to open the existing browser in Windows
- reference: https://superuser.com/a/1368878
