# Changelog

## [Unreleased](https://github.com/mesg-foundation/core/releases/tag/)

#### Breaking Changes
#### Added
#### Changed
#### Fixed
#### Removed
#### Documentation

## [v0.9.1](https://github.com/mesg-foundation/core/releases/tag/v0.9.1)

#### Added

- ([#993](https://github.com/mesg-foundation/engine/pull/993)) Display web3 http request logs in Marketplace

#### Changed

- ([#949](https://github.com/mesg-foundation/engine/pull/949)) Use MESG's IPFS node in CLI

#### Fixed

- ([#930](https://github.com/mesg-foundation/engine/pull/930)) Improve generated README when using command `service gen-doc`. ([#948](https://github.com/mesg-foundation/engine/pull/948)). ([#960](https://github.com/mesg-foundation/engine/pull/960)).
- ([#934](https://github.com/mesg-foundation/engine/pull/934)) Return error when image is passed in configuration in mesg.yml
- ([#929](https://github.com/mesg-foundation/engine/pull/929)) Show better error when deploying service

#### Documentation

- ([#953](https://github.com/mesg-foundation/engine/pull/953)) Fix links to docs

## [v0.9.0](https://github.com/mesg-foundation/core/releases/tag/v0.9.0)

#### Breaking Changes

- ([#731](https://github.com/mesg-foundation/core/pull/731)) Deterministic service hash. ([#804](https://github.com/mesg-foundation/core/pull/804)). ([#877](https://github.com/mesg-foundation/core/pull/877)).
- ([#801](https://github.com/mesg-foundation/core/pull/801)) Add Service's hash to reply of Deploy API
- ([#849](https://github.com/mesg-foundation/core/pull/849)) Use base58 to encode service hash.
- ([#860](https://github.com/mesg-foundation/core/pull/860)) Separate service's configuration from service's dependencies. ([#880](https://github.com/mesg-foundation/core/pull/880)).
- ([#866](https://github.com/mesg-foundation/core/pull/866)) Rename service's `volumesfrom` property.
- ([#905](https://github.com/mesg-foundation/core/pull/905)) Add version to database path to prevent decoding error and loss of data.

#### Added

- ([#535](https://github.com/mesg-foundation/core/pull/535)) Run MESG with MESG Services! [\#567](https://github.com/mesg-foundation/core/pull/567).
- ([#755](https://github.com/mesg-foundation/core/pull/755)) Implementation of the MESG Marketplace. The Marketplace allows the distribution and reutilization of MESG Services. Check `mesg-core marketplace` commands. ([#778](https://github.com/mesg-foundation/core/pull/778)). ([#788](https://github.com/mesg-foundation/core/pull/788)). ([#810](https://github.com/mesg-foundation/core/pull/810)). ([#817](https://github.com/mesg-foundation/core/pull/817)). ([#826](https://github.com/mesg-foundation/core/pull/826)). ([#828](https://github.com/mesg-foundation/core/pull/828)). ([#829](https://github.com/mesg-foundation/core/pull/829)). ([#837](https://github.com/mesg-foundation/core/pull/837)). ([#843](https://github.com/mesg-foundation/core/pull/843)). ([#844](https://github.com/mesg-foundation/core/pull/844)). ([#845](https://github.com/mesg-foundation/core/pull/845)). ([#853](https://github.com/mesg-foundation/core/pull/853)). ([#854](https://github.com/mesg-foundation/core/pull/854)). ([#863](https://github.com/mesg-foundation/core/pull/863)). ([#864](https://github.com/mesg-foundation/core/pull/864)). ([#868](https://github.com/mesg-foundation/core/pull/868)). ([#874](https://github.com/mesg-foundation/core/pull/874)). ([#883](https://github.com/mesg-foundation/core/pull/883)). ([#899](https://github.com/mesg-foundation/core/pull/899)). ([#898](https://github.com/mesg-foundation/core/pull/898)). ([#897](https://github.com/mesg-foundation/core/pull/897)). ([#896](https://github.com/mesg-foundation/core/pull/896)). ([#902](https://github.com/mesg-foundation/core/pull/902)). ([#901](https://github.com/mesg-foundation/core/pull/901)). ([#906](https://github.com/mesg-foundation/core/pull/906)). ([#907](https://github.com/mesg-foundation/core/pull/907)). ([#908](https://github.com/mesg-foundation/core/pull/908)). ([#909](https://github.com/mesg-foundation/core/pull/909)). ([#924](https://github.com/mesg-foundation/core/pull/924)). ([#926](https://github.com/mesg-foundation/core/pull/926)). ([#927](https://github.com/mesg-foundation/core/pull/927)). ([#936](https://github.com/mesg-foundation/core/pull/936)). ([#938](https://github.com/mesg-foundation/core/pull/938)). ([#939](https://github.com/mesg-foundation/core/pull/939)). ([#942](https://github.com/mesg-foundation/core/pull/942)). ([#943](https://github.com/mesg-foundation/core/pull/943)).
- ([#757](https://github.com/mesg-foundation/core/pull/757)) Read `.dockerignore` in dev and deploy commands.
- ([#779](https://github.com/mesg-foundation/core/pull/779)) Implementation of the MESG Wallet. Check `mesg-core wallet`. ([#807](https://github.com/mesg-foundation/core/pull/807)). ([#809](https://github.com/mesg-foundation/core/pull/809)). ([#812](https://github.com/mesg-foundation/core/pull/812)). ([#937](https://github.com/mesg-foundation/core/pull/937)).
- ([#781](https://github.com/mesg-foundation/core/pull/781)) Improve validation of service definition. ([#869](https://github.com/mesg-foundation/core/pull/869)).

#### Changed

- ([#823](https://github.com/mesg-foundation/core/pull/823)) Improve commands `service init` and `service gendoc`.
- ([#875](https://github.com/mesg-foundation/core/pull/875)) Improve JSON parsing error message.
- ([#790](https://github.com/mesg-foundation/core/pull/790)) Refactor. ([#792](https://github.com/mesg-foundation/core/pull/792)). ([#816](https://github.com/mesg-foundation/core/pull/816)). ([#805](https://github.com/mesg-foundation/core/pull/805)). ([#813](https://github.com/mesg-foundation/core/pull/813)). ([#839](https://github.com/mesg-foundation/core/pull/839)). ([#847](https://github.com/mesg-foundation/core/pull/847)). ([#850](https://github.com/mesg-foundation/core/pull/850)). ([#852](https://github.com/mesg-foundation/core/pull/852)). ([#855](https://github.com/mesg-foundation/core/pull/855)). ([#858](https://github.com/mesg-foundation/core/pull/858)). ([#867](https://github.com/mesg-foundation/core/pull/867)). ([#859](https://github.com/mesg-foundation/core/pull/859)). ([#870](https://github.com/mesg-foundation/core/pull/870)). ([#871](https://github.com/mesg-foundation/core/pull/871)). ([#872](https://github.com/mesg-foundation/core/pull/872)). ([#873](https://github.com/mesg-foundation/core/pull/873)). ([#881](https://github.com/mesg-foundation/core/pull/881)). ([#893](https://github.com/mesg-foundation/core/pull/893)). ([#892](https://github.com/mesg-foundation/core/pull/892)). ([#891](https://github.com/mesg-foundation/core/pull/891)). ([#890](https://github.com/mesg-foundation/core/pull/890)). ([#889](https://github.com/mesg-foundation/core/pull/889)). ([#888](https://github.com/mesg-foundation/core/pull/888)). ([#886](https://github.com/mesg-foundation/core/pull/886)). ([#885](https://github.com/mesg-foundation/core/pull/885)). ([#884](https://github.com/mesg-foundation/core/pull/884)). ([#903](https://github.com/mesg-foundation/core/pull/903)). ([#919](https://github.com/mesg-foundation/core/pull/919)).

#### Fixed

- ([#771](https://github.com/mesg-foundation/core/pull/771)) Fix gRPC stream acknowledgement.
- ([#772](https://github.com/mesg-foundation/core/pull/772)) Improve command logs errors.
- ([#820](https://github.com/mesg-foundation/core/pull/820)) Fix container package.

#### Removed
#### Documentation

## [v0.8.1](https://github.com/mesg-foundation/core/releases/tag/v0.8.1)

#### Fixed

- ([#774](https://github.com/mesg-foundation/core/pull/774)) Update keep alive of client to 5min to prevent spamming the server.

#### Documentation

- ([#762](https://github.com/mesg-foundation/core/pull/762)) Fix and improve guide start. ([#763](https://github.com/mesg-foundation/core/pull/763)).

## [v0.8.0](https://github.com/mesg-foundation/core/releases/tag/v0.8.0)

#### Added

- ([#690](https://github.com/mesg-foundation/core/pull/690)) Support service deployments from tarball urls.
- ([#732](https://github.com/mesg-foundation/core/pull/732)) Support multiple service id or hash for commands `service start` and `service stop`.
- ([#726](https://github.com/mesg-foundation/core/pull/726)) Add flag to command `start` to force colors in logs of Core.

#### Changed

- ([#734](https://github.com/mesg-foundation/core/pull/734)) Returns service sid in commands instead of hash.
- ([#724](https://github.com/mesg-foundation/core/pull/724)) Changed system services deployment system. ([#727](https://github.com/mesg-foundation/core/pull/727)). ([#725](https://github.com/mesg-foundation/core/pull/725)). ([#743](https://github.com/mesg-foundation/core/pull/743)).

#### Fixed

- ([#738](https://github.com/mesg-foundation/core/pull/738)) Fix stream disconnection because of more than 15min of inactivity. ([#739](https://github.com/mesg-foundation/core/pull/739)). ([#742](https://github.com/mesg-foundation/core/pull/742)). ([#744](https://github.com/mesg-foundation/core/pull/744)).

#### Documentation

- ([#721](https://github.com/mesg-foundation/core/pull/721)) Move documentation to [dedicated repository](https://github.com/mesg-foundation/docs).

## [v0.7.0](https://github.com/mesg-foundation/core/releases/tag/v0.7.0)

#### Added

- ([#677](https://github.com/mesg-foundation/core/pull/677)) Stream acknowledgement system. The core notifies client when streams are ready.
- ([#679](https://github.com/mesg-foundation/core/pull/679)) Add support of repeated parameters to service definition. ([#680](https://github.com/mesg-foundation/core/pull/680)). ([#684](https://github.com/mesg-foundation/core/pull/684)).
- ([#682](https://github.com/mesg-foundation/core/pull/682)) Add support of type Any to service definition. ([#689](https://github.com/mesg-foundation/core/pull/689)).
- ([#691](https://github.com/mesg-foundation/core/pull/691)) Add database transaction mechanism to database execution.
- ([#696](https://github.com/mesg-foundation/core/pull/696)) Add support of nested type definition for type Object.
- ([#704](https://github.com/mesg-foundation/core/pull/704]) Move go-service to package client/service.

#### Changed

- ([#688](https://github.com/mesg-foundation/core/pull/688)) Change sid auto-generated prefix.
- ([#699](https://github.com/mesg-foundation/core/pull/699)) Updated to golang v1.11.4.

#### Fixed

- ([#687](https://github.com/mesg-foundation/core/pull/687)) Fix execution generated id.
- ([#703](https://github.com/mesg-foundation/core/pull/703)) Return error when core is not running in command dev and deploy.

#### Removed

- ([#675](https://github.com/mesg-foundation/core/pull/675)) Remove workflow grpc client.
- ([#693](https://github.com/mesg-foundation/core/pull/693)) Remove vendor folder.

## [v0.6.0](https://github.com/mesg-foundation/core/releases/tag/v0.6.0)

#### Added

- ([\#641](https://github.com/mesg-foundation/core/pull/641)) Services definition accept env variables. Users can override them on deploy. [\#660](https://github.com/mesg-foundation/core/pull/660). [\#666](https://github.com/mesg-foundation/core/pull/666).
- ([\#651](https://github.com/mesg-foundation/core/pull/651)) Error added in task execution result.

#### Changed

- ([\#611](https://github.com/mesg-foundation/core/pull/611)) Switch to go1.11.
- ([\#648](https://github.com/mesg-foundation/core/pull/672)) Print all service definition in command `service detail`.
- ([\#649](https://github.com/mesg-foundation/core/pull/649)) Lowercase sid.

#### Documentation

- ([\#638](https://github.com/mesg-foundation/core/pull/638)) Fix marketplace link
- ([\#643](https://github.com/mesg-foundation/core/pull/643)) Add instruction to start the core without CLI
- ([\#656](https://github.com/mesg-foundation/core/pull/656)) Show instruction to create manually system services folder
- ([\#665](https://github.com/mesg-foundation/core/pull/665)) Add favicon

## [v0.5.0](https://github.com/mesg-foundation/core/releases/tag/v0.5.0)

#### Breaking Changes

- ([\#608](https://github.com/mesg-foundation/core/pull/608)) Rename "command" property and add "args" property in service definition.

#### Added

- ([\#583](https://github.com/mesg-foundation/core/pull/583)) Add property Sid (Service ID) in service definition file. Allow a service to reuse the same volumes after stopping. [\#627](https://github.com/mesg-foundation/core/pull/627). [\#613](https://github.com/mesg-foundation/core/pull/613). [\#619](https://github.com/mesg-foundation/core/pull/619).

#### Changed

- ([\#580](https://github.com/mesg-foundation/core/pull/580)) Refactor package Daemon.
- ([\#588](https://github.com/mesg-foundation/core/pull/588)) Refactor tests of package container.
- ([\#604](https://github.com/mesg-foundation/core/pull/604)) Improve hash function.
- ([\#609](https://github.com/mesg-foundation/core/pull/609)) Delete all service in parallel in commands.
- ([\#615](https://github.com/mesg-foundation/core/pull/615)) Remove initialization of swarm but display useful error.
- ([\#617](https://github.com/mesg-foundation/core/pull/617)) Improve template of command service gen doc.
- ([\#630](https://github.com/mesg-foundation/core/pull/630)) Rename service id to hash.

#### Fixed

- ([\#585](https://github.com/mesg-foundation/core/pull/585)) Handle gracefully task executions without inputs.
- ([\#598](https://github.com/mesg-foundation/core/pull/598)) Start service dependencies one by one. Solve issue when dependencies request access to same resource.

#### Documentation

- ([\#568](https://github.com/mesg-foundation/core/pull/568)) Update what-is-mesg.md.
- ([\#569](https://github.com/mesg-foundation/core/pull/569)) Update README.md.
- ([\#620](https://github.com/mesg-foundation/core/pull/620)) Add docker swarm init steps to doc.

## [v0.4.0](https://github.com/mesg-foundation/core/releases/tag/v0.4.0)

#### Added

- ([#534](https://github.com/mesg-foundation/core/pull/534)) Access service dependencies based on the name of the dependency through the network.
- ([#555](https://github.com/mesg-foundation/core/pull/555)) Add more logs on the CLI.

#### Changed

- ([#560](https://github.com/mesg-foundation/core/pull/560)) Store executions in a database - Fix memory leak [#542](https://github.com/mesg-foundation/core/pull/542)

#### Fixed

- ([#553](https://github.com/mesg-foundation/core/pull/553)) UI issue with service execute command.
- ([#552](https://github.com/mesg-foundation/core/pull/552)) Service dev command returns with an error when needed.
- ([#526](https://github.com/mesg-foundation/core/pull/526)) Improve container deletion when a service is stopped.
- ([#524](https://github.com/mesg-foundation/core/pull/524)) Fix sync issue on status send chans and sync issue on gRPC deploy stream sends.

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
