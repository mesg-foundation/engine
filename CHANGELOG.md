# Changelog

## [Unreleased]

#### Changed
- (#174) Update CI to build version based on tags
- (#173) Use official Docker client
- (#175) Changed the struct to use to start docker service

#### Added
- (#174) Add CHANGELOG.md file
- (#179) [API] Add `eventFilter` on `ListenEvent` API to get notification when an event with a specific name occurs
         [API] Add `taskFilter` on `ListenResult` API to get notification when a result from a specific task occurs
         [API] Add `outputFilter` on `ListenResult` API to get notification when a result returns a specific output

#### Removed

#### Fixed
- (#179) [Doc] Outdated documentation for the CLI