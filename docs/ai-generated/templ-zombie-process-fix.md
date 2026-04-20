# Fixing Templ Hot-Reloading and Ghost Processes

This document explains a known issue with hot-reloading Go applications using `templ generate --watch` alongside the standard `go run` command, along with the permanent fix applied to this project.

## The Bug: "Zombie/Ghost" Servers
When using `templ generate --watch --cmd="go run ..."`, stopping the development server (e.g., by pressing `Ctrl + C` in the terminal) can sometimes result in the compiled Go binary silently surviving in the background.

This happens because of the process supervision tree:
1. `templ --watch` acts as a proxy supervisor, listening for file changes and relaying commands to the Go runtime.
2. When the command given is `go run`, Go quietly builds a temporary executable deep in a cache folder (`/var/folders/...` or `/tmp/...`) and spins it up as a child process.
3. When `templ` receives a stop signal (`SIGINT`) from the terminal, it relays that signal to `go run`.
4. However, `go run` on macOS/Linux famously mishandles abrupt terminations. It will shut itself down without officially transferring the stop signal into its hidden temporary executable.
5. As a result, the backend API gets abandoned by the terminal but adopted by the OS. It sits running endlessly holding onto port 8080.

If you refresh your browser, the `templ proxy` might connect to this "ghost" server, serving confusing, out-of-date UI elements even though the actual source code was updated or reverted.

## The Fix: Bypassing `go run`
To guarantee a clean terminal shutdown and eliminate zombie processes, we bypass the `go run` wrapper program completely. Instead, we use `go build` to drop a native binary to a tracked location, and then run the raw binary natively.

In `Taskfile.yml`, the command is written as:
```yaml
  templ:
    desc: Templ hot reload, watches all .go files
    cmds:
      - templ generate --watch --proxy="http://localhost:{{.PORT}}" --open-browser=false --cmd="task dev-server"

  dev-server:
    cmds:
      - go build -tags=debug -o {{.BUILD_PATH}}/main-dev {{.SERVER_PATH}}
      - "{{.BUILD_PATH}}/main-dev"
```

## Why it works flawlessly
- **Native Signal Handling:** By dropping the buggy `go run` middleman, the `SIGINT` (Ctrl+C) stop signal coming from `templ` passes straight to `./dist/main-dev`. The binary respects the stop sequence and kills itself cleanly, leaving zero orphans.
- **Isolated Binaries:** Dumping the volatile development build natively into `./dist` (a directory already ignored by `.gitignore`) specifically named `main-dev` keeps the environment completely clean and guarantees separation from the production build (`main`).
- **Retains Hot Reloading:** We purposefully kept the build logic directly inside `templ`'s `--cmd` parameter. This means `templ` stays in the loop: when `.templ` files change, it successfully issues the build sequence, safely destroys the old executable, and smoothly patches the proxy to inject new HTML elements.

