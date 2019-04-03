# GO MOD
```sh
mkdir go-module-test
cd go-module-test
git init
git remote add origin git@gitlab.com:endevel/go-module-test.git
go mod init gitlab.com/endevel/go-module-test
go build 
./go-module-test
```

## Daily Workflow
Note there was no go get required in the example above.

Your typical day-to-day workflow can be:

- Add import statements to your .go code as needed.
- Standard commands like go build or go test will automatically add new dependencies as needed to satisfy imports (updating go.mod and downloading the new dependencies).
- When needed, more specific versions of dependencies can be chosen with commands such as go get foo@v1.2.3, go get foo@master, go get foo@e3702bed2, or by editing go.mod directly.

A brief tour of other common functionality you might use:

* `go list -m all` — View final versions that will be used in a build for all direct and indirect dependencies ([details](https://github.com/golang/go/wiki/Modules#version-selection))
* `go list -u -m all` — View available minor and patch upgrades for all direct and indirect dependencies ([details](https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies))
* `go get -u` or `go get -u=patch` — Update all direct and indirect dependencies to latest minor or patch upgrades ([details](https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies))
* `go build ./...` or `go test ./...` — Build or test all packages in the module when run from the module root directory ([details](https://github.com/golang/go/wiki/Modules#how-to-define-a-module))
* `go mod tidy` — Prune any no-longer-needed dependencies from `go.mod` and add any dependencies needed for other combinations of OS, architecture, and build tags ([details](https://github.com/golang/go/wiki/Modules#how-to-prepare-for-a-release))
* `replace` directive or `gohack` — Use a fork, local copy or exact version of a dependency ([details](https://github.com/golang/go/wiki/Modules#when-should-i-use-the-replace-directive))
* `go mod vendor` — Optional step to create a `vendor` directory ([details](https://github.com/golang/go/wiki/Modules#how-do-i-use-vendoring-with-modules-is-vendoring-going-away))

## go.mod
A module is defined by a tree of Go source files with a go.mod file in the tree's root directory. Module source code may be located outside of GOPATH. There are four directives: module, require, replace, exclude.

