# Changelog

## [Unreleased](https://github.com/mesg-foundation/core/releases/tag/)

#### Added
#### Changed
#### Fixed
#### Removed

## [v0.4.0](https://github.com/mesg-foundation/core/releases/tag/v0.4.0)

#### Added

- ([#534](https://github.com/mesg-foundation/core/pull/534)) Access service dependencies containers based on the name of the dependency
- ([#555](https://github.com/mesg-foundation/core/pull/555)) Add more logs on the CLI

#### Changed

- ([#560](https://github.com/mesg-foundation/core/pull/560)) Store executions in a database - Optimize memory of Core [#542](https://github.com/mesg-foundation/core/pull/542)

#### Fixed

- ([#553](https://github.com/mesg-foundation/core/pull/553)) UI issue with service execute command
- ([#552](https://github.com/mesg-foundation/core/pull/552)) service dev command return with an error when needed
- ([#526](https://github.com/mesg-foundation/core/pull/526)) Improve container deletion when a service is stopped
- ([#524](https://github.com/mesg-foundation/core/pull/524)) fix sync issue on status send chans and sync issue on gRPC deploy stream sends

## [v0.3.0](https://github.com/mesg-foundation/core/releases/tag/v0.3.0)

#### Added

- ([#392](https://github.com/mesg-foundation/core/pull/392)) **BREAKING CHANGE.** Add support for `.dockerignore`. Remove support of `.mesgignore` [#498](https://github.com/mesg-foundation/core/pull/498).
- ([#383](https://github.com/mesg-foundation/core/pull/383)) New API package. [#386](https://github.com/mesg-foundation/core/pull/386). [#444](https://github.com/mesg-foundation/core/pull/444). [#486](https://github.com/mesg-foundation/core/pull/486). [#488](https://github.com/mesg-foundation/core/pull/488).
- ([#409](https://github.com/mesg-foundation/core/pull/409)) Add required validations on service's task, event and output data.
- ([#432](https://github.com/mesg-foundation/core/pull/432)) Configuration of the CLI's output with `--no-color` and `--no-spinner` flags. Colorize JSON. [#453](https://github.com/mesg-foundation/core/pull/453). [#480](https://github.com/mesg-foundation/core/pull/480). [#484](https://github.com/mesg-foundation/core/pull/484).
- ([#435](https://github.com/mesg-foundation/core/pull/435)) Command `service logs` accepts multiple dependency filters with multiple use of `-d` flag.
- ([#478](https://github.com/mesg-foundation/core/pull/478)) Allow multiple core to run on the same computer.
- ([#493](https://github.com/mesg-foundation/core/pull/493)) Support numbers in service task's key, event's key and output's key
- ([#499](https://github.com/mesg-foundation/core/pull/499)) Return service's status from API

#### Changed

- ([#371](https://github.com/mesg-foundation/core/pull/371)) Delegate deployment of Service to Core. [#469](https://github.com/mesg-foundation/core/pull/469).
- ([#404](https://github.com/mesg-foundation/core/pull/404)) Change building tool.
- ([#413](https://github.com/mesg-foundation/core/pull/413)) Improve command `service dev`. [#459](https://github.com/mesg-foundation/core/pull/459).
- ([#417](https://github.com/mesg-foundation/core/pull/417)) Service refactoring. [#402](https://github.com/mesg-foundation/core/pull/402). [#414](https://github.com/mesg-foundation/core/pull/414). [#454](https://github.com/mesg-foundation/core/pull/454). [#458](https://github.com/mesg-foundation/core/pull/458). [#464](https://github.com/mesg-foundation/core/pull/464). [#472](https://github.com/mesg-foundation/core/pull/472). [#490](https://github.com/mesg-foundation/core/pull/490). [#491](https://github.com/mesg-foundation/core/pull/491).
- ([#419](https://github.com/mesg-foundation/core/pull/419)) Use Docker volumes for services. [#477](https://github.com/mesg-foundation/core/pull/477).
- ([#427](https://github.com/mesg-foundation/core/pull/427)) Refactor package Config
- ([#481](https://github.com/mesg-foundation/core/pull/481)) Refactor package Database
- ([#485](https://github.com/mesg-foundation/core/pull/485)) Improve CLI output. [#521](https://github.com/mesg-foundation/core/pull/521).
- Tests improvements. [#381](https://github.com/mesg-foundation/core/pull/381). [#384](https://github.com/mesg-foundation/core/pull/384). [#391](https://github.com/mesg-foundation/core/pull/391). [#446](https://github.com/mesg-foundation/core/pull/446). [#447](https://github.com/mesg-foundation/core/pull/447). [#466](https://github.com/mesg-foundation/core/pull/466). [#489](https://github.com/mesg-foundation/core/pull/489). [#501](https://github.com/mesg-foundation/core/pull/501). [#504](https://github.com/mesg-foundation/core/pull/504). [#506](https://github.com/mesg-foundation/core/pull/506).

#### Fixed

- ([#401](https://github.com/mesg-foundation/core/pull/401)) Gracefully stop gRPC servers.
- ([#429](https://github.com/mesg-foundation/core/pull/429)) Fix issue when stopping services. [#505](https://github.com/mesg-foundation/core/pull/505). [#526](https://github.com/mesg-foundation/core/pull/526).
- ([#476](https://github.com/mesg-foundation/core/pull/476)) Improve database error handling.
- ([#482](https://github.com/mesg-foundation/core/pull/482)) Fix Service hash changed when fetching from git.

#### Removed

- ([#410](https://github.com/mesg-foundation/core/pull/410)) Remove socket server in favor of the TCP server.

#### Documentation

- ([#415](https://github.com/mesg-foundation/core/pull/415)) Added hall-of-fame to README. Thanks [sergey48k](https://github.com/sergey48k).
- ([#423](https://github.com/mesg-foundation/core/pull/423)) Fix documentation issue.
- ([#474](https://github.com/mesg-foundation/core/pull/474)) Documentation/update ux.
- ([#509](https://github.com/mesg-foundation/core/pull/509)) Add forum link. [#513](https://github.com/mesg-foundation/core/pull/513).
- ([#510](https://github.com/mesg-foundation/core/pull/510)) Update ecosystem menu.
- ([#511](https://github.com/mesg-foundation/core/pull/511)) Update tutorial page.
- ([#512](https://github.com/mesg-foundation/core/pull/512)) Add sitemap.

## [v0.2.0](https://github.com/mesg-foundation/core/releases/tag/v0.2.0)

#### Added
- ([#242](https://github.com/mesg-foundation/core/pull/242)) Add more details in command `mesg-core service validate`
- ([#295](https://github.com/mesg-foundation/core/pull/295)) Added more validation on the API for the data of `executeTask`, `submitResult` and `emitEvent`. Now if data doesn't match the service file, the API returns an error
- ([#302](https://github.com/mesg-foundation/core/pull/302)) Possibility to use a config file in ~/.mesg/config.yml
- ([#303](https://github.com/mesg-foundation/core/pull/303)) Add command `service dev` that build and run the service with the logs
- ([#303](https://github.com/mesg-foundation/core/pull/303)) Add command `service execute` that execute a task on a service
- ([#316](https://github.com/mesg-foundation/core/pull/316)) Delete service when stoping the `service dev` command to avoid to keep all the versions of the services.
- ([#317](https://github.com/mesg-foundation/core/pull/317)) Add errors when trying to execute a service that is not running (nothing was happening before)
- ([#344](https://github.com/mesg-foundation/core/pull/344)) Add `service execute --data` flag to pass arguments as key=value.
- ([#362](https://github.com/mesg-foundation/core/pull/362)) Add `tags` list parameter for execution in order to be able to categorize and/or track a specific task execution.
- ([#362](https://github.com/mesg-foundation/core/pull/362)) Add possibility to listen to results with the specific tag(s)

#### Changed
- ([#282](https://github.com/mesg-foundation/core/pull/282)) Branch support added. You can now specify your branches with a `#branch` fragment at the end of your git url. E.g.: https://github.com/mesg-foundation/service-ethereum-erc20#websocket
- ([#299](https://github.com/mesg-foundation/core/pull/299)) Add more user friendly errors when failing to connect to the Core or Docker
- ([#356](https://github.com/mesg-foundation/core/pull/356)) Use github.com/stretchr/testify package
- ([#352](https://github.com/mesg-foundation/core/pull/352)) Use logrus logging package

#### Fixed
- ([#358](https://github.com/mesg-foundation/core/pull/358)) Fix goroutine leaks on api package handlers where gRPC streams are used. Handlers now doesn't block forever by exiting on context cancellation and stream.Send() errors.

#### Removed
- ([#303](https://github.com/mesg-foundation/core/pull/303)) Deprecate command `service test` in favor of `service dev` and `service execute`

## [v0.1.0](https://github.com/mesg-foundation/core/releases/tag/v0.1.0)

#### Added
- ([#267](https://github.com/mesg-foundation/core/pull/267)) Add Command `service gen-doc` that generate a `README.md` in the service based on the informations of the `mesg.yml`
- ([#246](https://github.com/mesg-foundation/core/pull/246)) Add .mesgignore to excluding file from the Docker build

#### Changed
- ([#247](https://github.com/mesg-foundation/core/pull/247)) Update the `service init` command to have initial tasks and events
- ([#257](https://github.com/mesg-foundation/core/pull/257)) Update the `service init` command to fetch for template based on the https://github.com/mesg-foundation/awesome/blob/master/templates.json file but also custom templates by giving the address of the template
- ([#261](https://github.com/mesg-foundation/core/pull/261)) **BREAKING** More consistancy between the APIs, rename `taskData` into `inputData` for the `executeTask` API

#### Fixed
- ([#246](https://github.com/mesg-foundation/core/pull/246)) Ignore files during Docker build

## [v0.1.0-beta3](https://github.com/mesg-foundation/core/releases/tag/v0.1.0-beta3)

#### Added
- ([#246](https://github.com/mesg-foundation/core/pull/246)) Add .mesgignore to excluding file from the Docker build

#### Changed
- ([#247](https://github.com/mesg-foundation/core/pull/247)) Update the `service init` command to have initial tasks and events
- ([#257](https://github.com/mesg-foundation/core/pull/257)) Update the `service init` command to fetch for template based on the https://github.com/mesg-foundation/awesome/blob/master/templates.json file but also custom templates by giving the address of the template
- ([#261](https://github.com/mesg-foundation/core/pull/261)) **BREAKING** More consistancy between the APIs, rename `taskData` into `inputData` for the `executeTask` API

#### Fixed
- ([#246](https://github.com/mesg-foundation/core/pull/246)) Ignore files during Docker build

## [v0.1.0-beta2](https://github.com/mesg-foundation/core/releases/tag/v0.1.0-beta2)

#### Added
- ([#174](https://github.com/mesg-foundation/core/pull/174)) Add CHANGELOG.md file
- ([#179](https://github.com/mesg-foundation/core/pull/179)) Add filters for the core API
  - [API] Add `eventFilter` on `ListenEvent` API to get notification when an event with a specific name occurs
  - [API] Add `taskFilter` on `ListenResult` API to get notification when a result from a specific task occurs
  - [API] Add `outputFilter` on `ListenResult` API to get notification when a result returns a specific output
- ([#183](https://github.com/mesg-foundation/core/pull/183)) Add a `configuration` attribute in the `mesg.yml` file to accept docker configuration for your service
- ([#187](https://github.com/mesg-foundation/core/pull/187)) Stop all services when the MESG Core stops
- ([#190](https://github.com/mesg-foundation/core/pull/190)) Possibility to `test` or `deploy` a service from a git or GitHub url
- ([#233](https://github.com/mesg-foundation/core/pull/233)) Add logs in the `service test` command with service logs by default and all dependencies logs with the `--full-logs` flag
- ([#235](https://github.com/mesg-foundation/core/pull/235)) Add `ListServices` and `GetService` APIs

#### Changed
- ([#174](https://github.com/mesg-foundation/core/pull/174)) Update CI to build version based on tags
- ([#173](https://github.com/mesg-foundation/core/pull/173)) Use official Docker client
- ([#175](https://github.com/mesg-foundation/core/pull/175)) Changed the struct to use to start docker service
- ([#181](https://github.com/mesg-foundation/core/pull/181)) MESG Core and Service start and stop functions wait for the docker container to actually run or stop.
- ([#183](https://github.com/mesg-foundation/core/pull/183)) **BREAKING** Docker image is automatically injected in the `mesg.yml` file for your service. Now `dependencies` attribute is for extra dependencies so for most of service this is not necessary anymore.
- ([#212](https://github.com/mesg-foundation/core/pull/212)) **BREAKING** Communication from services to core is now done through a token provided by the core
- ([#236](https://github.com/mesg-foundation/core/pull/236)) CLI only use the API
- ([#234](https://github.com/mesg-foundation/core/pull/234)) `service list` command now includes the status for every services

#### Fixed
- ([#179](https://github.com/mesg-foundation/core/pull/179)) [Doc] Outdated documentation for the CLI
- ([#185](https://github.com/mesg-foundation/core/pull/185)) Fix logs with extra characters when `mesg-core logs`


#### Removed
- ([#234](https://github.com/mesg-foundation/core/pull/234)) Remove command `service status` in favor of `service list` command that includes status