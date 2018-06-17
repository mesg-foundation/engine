# Changelog

## [Unreleased]

#### Changed
- (#174) Update CI to build version based on tags
- (#173) Use official Docker client
- (#175) Changed the struct to use to start docker service
- (#181) MESG Core and Service start and stop functions wait for the docker container to actually run or stop.
- (#183) **BREAKING** Docker image is automatically injected in the `mesg.yml` file for your service. Now `dependencies` attribute is for extra dependencies so for most of service this is not necessary anymore.
- (#212) **BREAKING** Communication from services to core is now done through a token provided by the core

#### Added
- (#174) Add CHANGELOG.md file
- (#179) Add filters for the core API
  - [API] Add `eventFilter` on `ListenEvent` API to get notification when an event with a specific name occurs
  - [API] Add `taskFilter` on `ListenResult` API to get notification when a result from a specific task occurs
  - [API] Add `outputFilter` on `ListenResult` API to get notification when a result returns a specific output
- (#183) Add a `configuration` attribute in the `mesg.yml` file to accept docker configuration for your service
- (#187) Stop all services when the MESG Core stops
- (#190) Possibility to `test` or `deploy` a service from a git or GitHub url
- (#235) Add `ListServices` and `GetService` APIs

#### Removed

#### Fixed
- (#179) [Doc] Outdated documentation for the CLI
- (#185) Fix logs with extra characters when `mesg-core logs`
