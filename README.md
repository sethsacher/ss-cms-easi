# CMS EASi Application

This repository contains the application code
for the CMS EASi (Easy Access to System Information).

## Codebase Layout

## Go

The server portion of the application is written in Go.
More information can be found in the [pkg documentation](./pkg/README.md).

## Application Setup

### Setup: Developer Setup

There are a number of things you'll need at a minimum
to be able to check out,
develop,
and run this project.

* Install [Homebrew](https://brew.sh)
  * Use the following command
    `/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"`
* We normally use the latest version of Go
  unless there's a known conflict
  (which will be announced by the team)
  or if we're in the time period just after a new version has been released.
  * Install it with Homebrew:
    `brew install go`
  * **Note**:
    If you have previously modified your PATH to point to a specific version of go,
    make sure to remove that.
    This would be either in your `.bash_profile` or `.bashrc`,
    and might look something like
    `PATH=$PATH:/usr/local/opt/go@1.12/bin`.
* Ensure you are using the latest version of bash for this project:
  * Install it with Homebrew:
    `brew install bash`
  * Update list of shells that users can choose from:

    ```bash
        [[ $(cat /etc/shells | grep /usr/local/bin/bash) ]] \
        || echo "/usr/local/bin/bash" | sudo tee -a /etc/shells
    ```

  * If you are using bash as your shell
    (and not zsh, fish, etc)
    and want to use the latest shell as well,
    then change it (optional): `chsh -s /usr/local/bin/bash`
  * Ensure that `/usr/local/bin` comes before `/bin`
    on your `$PATH` by running `echo $PATH`.
    Modify your path by editing `~/.bashrc` or `~/.bash_profile`
    and changing the `PATH`.
    Then source your profile with `source ~/.bashrc` or `~/.bash_profile`
    to ensure that your terminal has it.
* **Note**:
    If you have previously used yarn or Golang, please make sure none
    of them are pinned to an old version by running `brew list --pinned`.
    If they are pinned, please run `brew unpin <formula>`.
    You can upgrade these formulas instead of installing by running
    `brew upgrade <formula`.

### Setup: Git

Use your work email when making commits to our repositories.
The simplest path to correctness is setting global config:

  ```bash
  git config --global user.email "trussel@truss.works"
  git config --global user.name "Trusty Trussel"
  ```

If you drop the `--global` flag,
these settings will only apply to the current repo.
If you ever re-clone that repo or clone another repo,
you will need to remember to set the local config again.
You won't.
Use the global config. :-)

For web-based Git operations,
GitHub will use your primary email unless you choose
"Keep my email address private".
If you don't want to set your work address as primary,
please [turn on the privacy setting](https://github.com/settings/emails).

Note that with 2-factor-authentication enabled,
in order to push local code to GitHub through HTTPS,
you need to [create a personal access token](https://gist.github.com/ateucher/4634038875263d10fb4817e5ad3d332f)
and use that as your password.

### Setup: Golang

All of Go's tooling expects Go code to be checked out in a specific location.
Please read about [Go workspaces](https://golang.org/doc/code.html#Workspaces)
for a full explanation.
If you just want to get started,
then decide where you want all your go code to live
and configure the GOPATH environment variable accordingly.
For example,
if you want your go code to live at `~/code/go`,
you should add the following like to your `.bash_profile`:

  ```bash
  export GOPATH=~/code/go
  ```

Golang expects the `GOPATH` environment variable to be defined.
If you'd like to use the default location,
then add the following to your `.bash_profile`
or hard code the default value.
This line will set the GOPATH environment variable
to the value of `go env GOPATH`
if it is not already set.

  ```bash
  export GOPATH=${GOPATH:-$(go env GOPATH)}
  ```

**Regardless of where your go code is located**,
you need to add `$GOPATH/bin` to your `PATH`
so that executables installed with the go tooling can be found.
Add the following to your `.bash_profile`:

  ```bash
  export PATH=$(go env GOPATH)/bin:$PATH
  ```

Finally to have these changes applied to your shell,
you must `source` your profile:

  ```bash
  source ~/.bash_profile
  ```

You can confirm that the values exist with:

  ```bash
  env | grep GOPATH
  # Verify the GOPATH is correct
  env | grep PATH
  # Verify the PATH includes your GOPATH bin directory
  ```

### Setup: Project Checkout

You can checkout this repository by running
`git clone git@github.com:cmsgov/easi-app.git`.
Please check out the code in a directory like
`~/Projects/easi-app` and NOT in your `$GOPATH`. As an example:

  ```bash
  mkdir -p ~/Projects
  git clone git@github.com:cmsgov/easi-app.git
  cd easi-app
  ```

You will then find the code at `~/Projects/easi-app`.
You can check the code out anywhere EXCEPT inside your `$GOPATH`.
So this is customization that is up to you.

### Setup: direnv

TODO: Need direnv setup

### Setup: Yarn (temporary)

Run `brew install yarn` to install yarn.
Run `yarn install` to install dependencies.

This is temporary for setting up pre-commit package,
until we decide how to structure our CLI tools.

### Setup: Pre-Commit

Run `pre-commit install`
to install a pre-commit hook
into `./git/hooks/pre-commit`.
This is different than `brew install pre-commit`
and must be done
so that the hook will check files
you are about to commit to the repository.
Next install the pre-commit hook libraries
with `pre-commit install-hooks`.

### Setup: CircleCI (optional)

If you want to make changes to the CircleCI configuration, you will need to
install the `circleci` cli tool so that the changes can be validated by
pre-commit: `brew install circleci`

### Setup: Docker

To set up docker on your local machine:

```console
brew cask install docker
brew install docker-completion
```

Now you will need to start the Docker service: run Spotlight and type in
'docker'.

#### docker-compose

To run the app in Docker locally, we are using
[`docker-compose`](https://docs.docker.com/compose/):

```console
brew install docker-compose
brew install docker-compose-completion  # optional
```

Now, you can stand up the app & database locally in containers.  The app server
will stand up and migrations from your local filesystem will be run against the
db automatically:

```console
$ docker-compose up --detach
Starting app_db_1 ... done
Starting app_db_clean_1 ... done
Starting app_db_migrate_1 ... done
Starting app_easi_1       ... done
$ curl http://localhost:8080/api/v1/healthcheck
{"status":"pass","datetime":"APPLICATION_DATETIME","version":"APPLICATION_VERSION","timestamp":"APPLICATION_TS"}
```

Two options to take it down:

```console
docker-compose kill  # stops the running containers
docker-compose down  # stops and also removes the containers
```

You can also force rebuilding the images (e.g. after using `kill`) with
`docker-compose build`.

To inspect the database from your shell, `pgcli` is recommended:

```sh
brew install pgcli
```

Thanks to variables set in the `.envrc`, connecting is simple:

```console
$ pgcli
Server: PostgreSQL 11.6 (Debian 11.6-1.pgdg90+1)
Version: 2.2.0
Chat: https://gitter.im/dbcli/pgcli
Home: http://pgcli.com
postgres@localhost:postgres> SHOW server_version;
+-------------------------------+
| server_version                |
|-------------------------------|
| 11.6 (Debian 11.6-1.pgdg90+1) |
+-------------------------------+
SHOW
Time: 0.016s
postgres@localhost:postgres>
```

### Setup: Cloud Services

You may need to access cloud service
to develop the application.
This allows access to AWS resources (ex. SES Email).

Follow the instructions in the infra repo
[here](https://github.com/CMSgov/easi-infra#ctkey-wrapper).
You'll need to add the infra account environment variables
to your `.envrc.local`.
You can then run the `ctkey` command
to get/set AWS environment variables.

```bash
https_proxy=localhost:8888 \\
ctkey --username=$CTKEY_USERNAME \\
--password=$CTKEY_PASSWORD \\
--account=$AWS_ACCOUNT_ID \\
--url=$CTKEY_URL \\
--idms=$CT_IDMS \\
--iam-role=$CT_AWS_ROLE setenv
```

Eventually, we will move this over to wrapper
so developers do not need to manually run these commands.

### Live Reload Go with Air (Optional)

If you want to reload the Go application on changes locally,
you can use [Air](github.com/cosmtrek/air)

Install it:

```bash
go get -u github.com/cosmtrek/air
```

Run it:

```bash
air
```

It's not currently set up to run with docker.
You can edit the config [here](.air.conf)

## Build

### Swagger Generation

The EASi server uses Swagger generation
to access the CEDAR (data source) API.
Swagger specs can be downloaded from webMethods:

* [IMPL](webmethods-apigw.cedarimpl.cms.gov)

Put this file in ️`$CEDAR_DIRECTORY`
and name it `swagger-<env>.yaml` respectively.

If you haven't run go-swagger before,
you'll need to install it.
Run:

```go
go build -o bin/swagger github.com/go-swagger/go-swagger/cmd/swagger
```

Then, to generate the client run:

```go
swagger generate client -f $CEDAR_SWAGGER_FILE -c $CEDAR_DIRECTORY/gen/client -m $CEDAR_DIRECTORY/gen/models
```

### Golang cli app

To build the cli application in your local filesystem:

```sh
go build -a -o bin/easi ./cmd/easi
```

You can then access the tool with the `easi` command.

### Migrating the Database

To add a new migration, add a new file to the `migrations` directory
following the standard
`V_${last_migration_version + 1}_your_migration_name_here.sql`

## Testing

### Server tests

Once your developer environment is setup,
you can run tests with the `easi test` command.

If you run into into various `(typecheck)` errors when
running `easi test` follow the directions for [installing
golangci-lint](https://github.com/golangci/golangci-lint#install)
to upgrade golangci-lint.

### Cypress tests (End-to-end integration tests)

Run `npx cypress run` to run the tests in the CLI. To have a slightly more interactive
experience, you can instead run `npx cypress open`.

## Development and Debugging

### APIs

The APIs reside at `localhost:8080` when running.
To run a test request,
you can send a GET to the health check endpoint:
`curl localhost:8080/api/v1/healthcheck`

#### Authorization

Setting this `APP_ENV` environment variable to "local"
will turn off API authorization.
This allows you to debug quickly without retrieving Okta access tokens.

If you need to test the authorization,
unset that environment variable.
You can then retrieve access tokens
by logging into [the development app](dev.easi.cms.gov)
and copying `okta-token-storage/accessToken` from the browser's local storage.
Place this in the `Authorization` header
as `Bearer ${accessToken}`.
