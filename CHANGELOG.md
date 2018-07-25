# Changelog

## [Unreleased]

#### Changed
- (#282) Branch support added. You can now specify your branches with a `#branch` fragment at the end of your git url. E.g.: https://github.com/mesg-foundation/service-ethereum-erc20#websocket

#### Added
- (#242) Add more details in command `mesg-core service validate`
- (#302) Possibility to use a config file in ~/.mesg/config.yml

#### Removed
#### Fixed

## [v1.0.0]

#### Changed
- (#247) Update the `service init` command to have initial tasks and events
- (#257) Update the `service init` command to fetch for template based on the https://github.com/mesg-foundation/awesome/blob/master/templates.json file but also custom templates by giving the address of the template
- (#261) **BREAKING** More consistancy between the APIs, rename `taskData` into `inputData` for the `executeTask` API

#### Added
- (#267) Add Command `service gen-doc` that generate a `README.md` in the service based on the informations of the `mesg.yml`
- (#246) Add .mesgignore to excluding file from the Docker build

#### Removed

#### Fixed
- (#246) Ignore files during Docker build

## [v1.0.0-beta2]

#### Changed
- (#247) Update the `service init` command to have initial tasks and events
- (#257) Update the `service init` command to fetch for template based on the https://github.com/mesg-foundation/awesome/blob/master/templates.json file but also custom templates by giving the address of the template
- (#261) **BREAKING** More consistancy between the APIs, rename `taskData` into `inputData` for the `executeTask` API

#### Added
- (#246) Add .mesgignore to excluding file from the Docker build

#### Removed

#### Fixed
- (#246) Ignore files during Docker build

## [v1.0.0-beta]

#### Changed
- (#174) Update CI to build version based on tags
- (#173) Use official Docker client
- (#175) Changed the struct to use to start docker service
- (#181) MESG Core and Service start and stop functions wait for the docker container to actually run or stop.
- (#183) **BREAKING** Docker image is automatically injected in the `mesg.yml` file for your service. Now `dependencies` attribute is for extra dependencies so for most of service this is not necessary anymore.
- (#212) **BREAKING** Communication from services to core is now done through a token provided by the core
- (#236) CLI only use the API
- (#234) `service list` command now includes the status for every services

#### Added
- (#174) Add CHANGELOG.md file
- (#179) Add filters for the core API
  - [API] Add `eventFilter` on `ListenEvent` API to get notification when an event with a specific name occurs
  - [API] Add `taskFilter` on `ListenResult` API to get notification when a result from a specific task occurs
  - [API] Add `outputFilter` on `ListenResult` API to get notification when a result returns a specific output
- (#183) Add a `configuration` attribute in the `mesg.yml` file to accept docker configuration for your service
- (#187) Stop all services when the MESG Core stops
- (#190) Possibility to `test` or `deploy` a service from a git or GitHub url
- (#233) Add logs in the `service test` command with service logs by default and all dependencies logs with the `--full-logs` flag
- (#235) Add `ListServices` and `GetService` APIs

#### Removed
- (#234) Remove command `service status` in favor of `service list` command that includes status

#### Fixed
- (#179) [Doc] Outdated documentation for the CLI
- (#185) Fix logs with extra characters when `mesg-core logs`
