# Changelog

## [Unreleased]

#### Changed
- (#174) Update CI to build version based on tags
- (#173) Use official Docker client
- (#175) Changed the struct to use to start docker service
- (#181) Daemon and Service start and stop functions wait for the docker container to actually run or stop.
- (#183) **BREAKING** Docker image is automatically injected in the `mesg.yml` file for your service. Now `dependencies` attribute is for extra dependencies so for most of service this is not necessary anymore.

#### Added
- (#174) Add CHANGELOG.md file
- (#183) Add a `configuration` attribute in the `mesg.yml` file to accept docker configuration for your service
- (#187) Stop all services when the daemon stops

#### Removed

#### Fixed
- (#185) Fix logs with extra characters when `mesg-core daemon logs`