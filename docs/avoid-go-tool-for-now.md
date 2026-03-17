## In golang 1.24 they added a new feature "go tool"

the way I use it is similar to the idea of "dev dependencies" in npm. go tool is used
for installing tools/packages for development, not the packages that will be used for
building the binaries of the project.

### usage

```sh
go get -tool <tool>
```

### example

```sh
go get -tool sqlc
```

### The reason why we do this instead of just getting the tool through `go install <tool>`

The main advantage is that when setting up the project for development, the user can now
just run `go mod tidy` to install all the necessary tools along with the necessary packages

## There is performance drop in building when it comes to installing `go get -tool` instead of installing through `go install`

- read this well documented blog by howard john: https://blog.howardjohn.info/posts/go-tools-command/#shared-dependency-state

## Action items

- avoid installing packages through `go get -tool <tool>` for now until the performance issues are addressed
- for now we stick to our usual setup steps which is:

1. install task
2. run `task setup` to install our tools
