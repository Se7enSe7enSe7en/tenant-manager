# Tenant Manager
app to manage tenants

## Setup guide:
### 1: install task 
```
go get -tool github.com/go-task/task/v3/cmd/task@latest
```
### 2: use task to trigger the "setup" command 
```
go tool task setup
```

## WSL setup guide:
### (optional) Install wslview/wslu
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
