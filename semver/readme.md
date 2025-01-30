# app semver - use Semantic Versioning in your app

Package semver manages data for reporting semantic versions of your app.

The goal of the package is to provide a structured way to include semantic
version, hash and note information in a CI/CD pipeline:

```sh
# things you can do with the library
app -v          # → 1.2.3-rc
app --version   # → Full App Title 1.2.3-rc+f3c6
app --note      # → Some commit message maybe
app --gitHash   # → d57e36f5f84302c1          - first 16 chars
app --arch      # → arm64                     - GOARCH - expected architecture
```

This should allow you to manage docker images & documentation by using the app
to generate its own version.

Note that idiomatic golang required a prefix of **v** when adding versions to
packages. The underlying [semver][sv] spec does not require a **v** prefix.

This package does not use a **v** prefix and assumes that you'll insert the
letter **v** when needed in your workflow.

## information flow

* `/releases.yaml` is used to track the high level releases of the app
* `semver/semver.go` mixes the yaml and the linker data to display a version
* Linker data is provided with the go build command:
  ```sh
  go build -ldflags "-X gitlab.com/workspace/account/semver.SemVerInfo=\"'x_x_x_x_HASH_x_x_x_x_x_VALUE_x_x_x_x_AS_x_x_x_x_x_TEXT_x_x_x_x_x_x|2024-07-09||myclog|Command_Line_Of_Go'" .
  ```
  The format of the linker string is `commithash|date|suffix|appname|apptitle`:
  * `commithash` - the 40 digit hash from the repo used to build clog
  * `date` - an ISO 8601 date string representing the release/build date
  * `appName` - the command name usually used to run the program e.g clog
  * `app title` - a text string to describe the title. Due to a limitation in
    the construction of linker strings, **all spaces must be represented by
    underscores**. The build will replace underscores with spaces at runtime.
    See the sample above

## usage

Create an embedded yaml (or json) file to track the releases.

```yaml
# Dates must be in YYYY-MM-DD international ISO 8601 format
- {version: "0.2.0", date: 2024-12-02, codename: alpha, note: first rc docker}
- {version: "0.1.0", date: 2024-12-01, codename: alpha, note: first prototype}

```

The package assumes that the first entry is the one that is to be used and that
the rest are some kind of history - even if you're regressing the semantic version
because you're going backwards to a previous branch.

You can add extra fields, however, the ones shown are required.

## linking

A release build script can inject LDLinkerData{} variables:

```shell
  commitHash="$(git rev-list -1 HEAD)"
  printf -v buildDate '%(%Y-%m-%d)T' -1
  buildSuffix="" && [ -z "$(git branch  --show-current|grep main)" ] && buildSuffix="$(git branch  --show-current)"
  buildAppName=myapp
  buildAppTitle="My Awesome App With SemVer"
   # create linker data info:
  ldi="$commitHash|$buildDate|$buildSuffix|$buildAppName|$buildAppTitle"
  # use path to variable in the built project
  # use `go tool objdump -S myExecutable | grep /semver.SemVerInfo` to find the path
  linkerDataSemverPath=github.com/workspace/repo/semver.SemVerInfo
  # build with linker data
  GOOS=$OS GOARCH=$CPU go build -ldflags "-X $linkerDataSemverPath='$ldi'" -o /some/executable
````

to use in your code:

```golang
// Typically this is in the root folder making it easy to find
//go:embed releases.yaml
var embeddedFs embed.FS

```golang
    package main

    // pass an FS and a filePath for the semver to read history
    semver.Initialise(embeddedFs, "releases.yaml")

    fmt.Printf("short version=%s", semver.Info.Short)
    fmt.Printf(" long version=%s", semver.Info.Long)
```

[sv]: https://semver.org/
