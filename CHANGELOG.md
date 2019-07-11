# Changelog

## [Unreleased](https://github.com/mesg-foundation/engine/compare/master...dev)

#### Breaking Changes
#### Added
#### Changed
#### Fixed
#### Removed

## [v0.11.0](https://github.com/mesg-foundation/engine/releases/tag/v0.11.0)

### [Click here to see the release notes](https://forum.mesg.com/t/release-notes-of-engine-v0-11-cli-v1-1-and-js-library-v4/339)

#### Breaking Changes

- ([#1083](https://github.com/mesg-foundation/engine/pull/1083)) Remove old gRPC APIs.
- ([#1170](https://github.com/mesg-foundation/engine/pull/1170)) Change Service and Execution database version. You'll need to redeploy your services.

#### Added

- ([#1023](https://github.com/mesg-foundation/engine/pull/1023)) New gRPC Service APIs and SDK. ([#1066](https://github.com/mesg-foundation/engine/pull/1066)). ([#1067](https://github.com/mesg-foundation/engine/pull/1067)). ([#1068](https://github.com/mesg-foundation/engine/pull/1068)). ([#1071](https://github.com/mesg-foundation/engine/pull/1071)). ([#1077](https://github.com/mesg-foundation/engine/pull/1077)). ([#1097](https://github.com/mesg-foundation/engine/pull/1097)). ([#1099](https://github.com/mesg-foundation/engine/pull/1099)). ([#1105](https://github.com/mesg-foundation/engine/pull/1105)). ([#1107](https://github.com/mesg-foundation/engine/pull/1107)). ([#1110](https://github.com/mesg-foundation/engine/pull/1110)). ([#1112](https://github.com/mesg-foundation/engine/pull/1112)). ([#1113](https://github.com/mesg-foundation/engine/pull/1113)). ([#1128](https://github.com/mesg-foundation/engine/pull/1128)). ([#1153](https://github.com/mesg-foundation/engine/pull/1153)).
- ([#1033](https://github.com/mesg-foundation/engine/pull/1033)) New gRPC Instance APIs and SDK. ([#1034](https://github.com/mesg-foundation/engine/pull/1034)). ([#1035](https://github.com/mesg-foundation/engine/pull/1035)). ([#1036](https://github.com/mesg-foundation/engine/pull/1036)). ([#1037](https://github.com/mesg-foundation/engine/pull/1037)). ([#1060](https://github.com/mesg-foundation/engine/pull/1060)). ([#1074](https://github.com/mesg-foundation/engine/pull/1074)). ([#1075](https://github.com/mesg-foundation/engine/pull/1075)). ([#1076](https://github.com/mesg-foundation/engine/pull/1076)). ([#1078](https://github.com/mesg-foundation/engine/pull/1078)). ([#1079](https://github.com/mesg-foundation/engine/pull/1079)). ([#1087](https://github.com/mesg-foundation/engine/pull/1087)). ([#1102](https://github.com/mesg-foundation/engine/pull/1102)). ([#1106](https://github.com/mesg-foundation/engine/pull/1106)). ([#1109](https://github.com/mesg-foundation/engine/pull/1109)). ([#1117](https://github.com/mesg-foundation/engine/pull/1117)). ([#1122](https://github.com/mesg-foundation/engine/pull/1122)). ([#1137](https://github.com/mesg-foundation/engine/pull/1137)). ([#1138](https://github.com/mesg-foundation/engine/pull/1138)). ([#1156](https://github.com/mesg-foundation/engine/pull/1156)).
- ([#1043](https://github.com/mesg-foundation/engine/pull/1043)) New gRPC Execution APIs and SDK. ([#1052](https://github.com/mesg-foundation/engine/pull/1052)). ([#1064](https://github.com/mesg-foundation/engine/pull/1064)). ([#1065](https://github.com/mesg-foundation/engine/pull/1065)). ([#1070](https://github.com/mesg-foundation/engine/pull/1070)). ([#1124](https://github.com/mesg-foundation/engine/pull/1124)). ([#1132](https://github.com/mesg-foundation/engine/pull/1132)). ([#1135](https://github.com/mesg-foundation/engine/pull/1135)).
- ([#1053](https://github.com/mesg-foundation/engine/pull/1053)) New gRPC Event APIs and SDK. ([#1054](https://github.com/mesg-foundation/engine/pull/1054)). ([#1073](https://github.com/mesg-foundation/engine/pull/1073)). ([#1126](https://github.com/mesg-foundation/engine/pull/1126)). ([#1144](https://github.com/mesg-foundation/engine/pull/1144)). ([#1141](https://github.com/mesg-foundation/engine/pull/1141)).

#### Changed

- ([#1082](https://github.com/mesg-foundation/engine/pull/1082)) Server and SDK package cleanup. ([#1085](https://github.com/mesg-foundation/engine/pull/1085)). ([#1096](https://github.com/mesg-foundation/engine/pull/1096)).
- ([#1087](https://github.com/mesg-foundation/engine/pull/1087)) Update system service deployment system. ([#1156](https://github.com/mesg-foundation/engine/pull/1156)). ([#1160](https://github.com/mesg-foundation/engine/pull/1160)).
- ([#1094](https://github.com/mesg-foundation/engine/pull/1094)) Introduction of a Hash type to manage all encoding and decoding of hashes. ([#1072](https://github.com/mesg-foundation/engine/pull/1072)). ([#1098](https://github.com/mesg-foundation/engine/pull/1098)).
- ([#1045](https://github.com/mesg-foundation/engine/pull/1045)) Update system service Marketplace to use new gRPC APIs. ([#1166](https://github.com/mesg-foundation/engine/pull/1166)).
- ([#1148](https://github.com/mesg-foundation/engine/pull/1148)) Update system service EthWallet to use new gRPC APIs. ([#1167](https://github.com/mesg-foundation/engine/pull/1167)).
- ([#1151](https://github.com/mesg-foundation/engine/pull/1151)) Namespace simplification in package Container.

## [v0.10.1](https://github.com/mesg-foundation/engine/releases/tag/v0.10.1)

### [Click here to see the release notes](https://forum.mesg.com/t/mesg-engine-v0-10-js-cli-and-js-library-v3-0-0-release-notes/317)

#### Fixed

- ([#1050](https://github.com/mesg-foundation/engine/pull/1050)) Fix service Marketplace backward compatibility.

## [v0.10.0](https://github.com/mesg-foundation/engine/releases/tag/v0.10.0)

### [Click here to see the release notes](https://forum.mesg.com/t/mesg-engine-v0-10-js-cli-and-js-library-v3-0-0-release-notes/317)

#### Breaking Changes

- ([#928](https://github.com/mesg-foundation/engine/pull/928)) Rename `volumesfrom` to `volumesFrom` in services' `mesg.yml`.
- ([#974](https://github.com/mesg-foundation/engine/pull/974)) Reduce the number of service outputs to one. An output can still contain multiple parameters. ([#963](https://github.com/mesg-foundation/engine/pull/963)). ([#964](https://github.com/mesg-foundation/engine/pull/964)). ([#965](https://github.com/mesg-foundation/engine/pull/965)). ([#967](https://github.com/mesg-foundation/engine/pull/967)). ([#971](https://github.com/mesg-foundation/engine/pull/971)). ([#972](https://github.com/mesg-foundation/engine/pull/972)). ([#973](https://github.com/mesg-foundation/engine/pull/973)). ([#975](https://github.com/mesg-foundation/engine/pull/975)). ([#976](https://github.com/mesg-foundation/engine/pull/976)). ([#977](https://github.com/mesg-foundation/engine/pull/977)). ([#978](https://github.com/mesg-foundation/engine/pull/978)). ([#979](https://github.com/mesg-foundation/engine/pull/979)). ([#980](https://github.com/mesg-foundation/engine/pull/980)). ([#981](https://github.com/mesg-foundation/engine/pull/981)). ([#982](https://github.com/mesg-foundation/engine/pull/982)). ([#983](https://github.com/mesg-foundation/engine/pull/983)). ([#984](https://github.com/mesg-foundation/engine/pull/984)). ([#987](https://github.com/mesg-foundation/engine/pull/987)). ([#988](https://github.com/mesg-foundation/engine/pull/988)). ([#1028](https://github.com/mesg-foundation/engine/pull/1028)).
- ([#791](https://github.com/mesg-foundation/engine/pull/791)) Remove CLI from the repository. The new cli is available on a [dedicated repository](https://github.com/mesg-foundation/cli). ([#995](https://github.com/mesg-foundation/engine/pull/995)). ([#996](https://github.com/mesg-foundation/engine/pull/996)).
- ([#991](https://github.com/mesg-foundation/engine/pull/991)) Rename Core to Engine. ([#968](https://github.com/mesg-foundation/engine/pull/968)). ([#970](https://github.com/mesg-foundation/engine/pull/970)). ([#1002](https://github.com/mesg-foundation/engine/pull/1002)). ([#1003](https://github.com/mesg-foundation/engine/pull/1003)). ([#1004](https://github.com/mesg-foundation/engine/pull/1004)). ([#1020](https://github.com/mesg-foundation/engine/pull/1020)).
- ([#1032](https://github.com/mesg-foundation/engine/pull/1032)) Simplify Engine configs. ([#1038](https://github.com/mesg-foundation/engine/pull/1038)).

#### Added

- ([#1014](https://github.com/mesg-foundation/engine/pull/1014)) Introduce new deployment api. Only available for development purpose.

#### Changed

- ([#994](https://github.com/mesg-foundation/engine/pull/994)) Update execution database and api. ([#1006](https://github.com/mesg-foundation/engine/pull/1006)). ([#1007](https://github.com/mesg-foundation/engine/pull/1007)). ([#1041](https://github.com/mesg-foundation/engine/pull/1041)).
- ([#997](https://github.com/mesg-foundation/engine/pull/997)) Rename package `api` to `sdk`.
- ([#998](https://github.com/mesg-foundation/engine/pull/998)) Rename package `interface` to `server`.

#### Fixed

- ([#955](https://github.com/mesg-foundation/engine/pull/955)) Catch error when the volume is deleted and it did not exist.

#### Removed

- [#1001](https://github.com/mesg-foundation/engine/pull/1001) Remove web3 http request logs from Marketplace.

## [v0.9.1](https://github.com/mesg-foundation/engine/releases/tag/v0.9.1)

### [Click here to see the release notes](https://forum.mesg.com/t/mesg-core-v0-9-1-release-notes/293)

#### Added

- ([#993](https://github.com/mesg-foundation/engine/pull/993)) Display web3 http request logs in Marketplace.

#### Changed

- ([#949](https://github.com/mesg-foundation/engine/pull/949)) Use MESG's IPFS node in CLI.

#### Fixed

- ([#930](https://github.com/mesg-foundation/engine/pull/930)) Improve generated README when using command `service gen-doc`. ([#948](https://github.com/mesg-foundation/engine/pull/948)). ([#960](https://github.com/mesg-foundation/engine/pull/960)).
- ([#934](https://github.com/mesg-foundation/engine/pull/934)) Return error when an image is passed on a mesg.yml definition.
- ([#929](https://github.com/mesg-foundation/engine/pull/929)) Show more verbose error when deploying service.

#### Documentation

- ([#953](https://github.com/mesg-foundation/engine/pull/953)) Fix links to docs

## [v0.9.0](https://github.com/mesg-foundation/engine/releases/tag/v0.9.0)

### [Click here to see the release notes](https://forum.mesg.com/t/mesg-core-v0-9-release-notes/273)

#### Breaking Changes

- ([#731](https://github.com/mesg-foundation/engine/pull/731)) Deterministic service hash. ([#804](https://github.com/mesg-foundation/engine/pull/804)). ([#877](https://github.com/mesg-foundation/engine/pull/877)).
- ([#801](https://github.com/mesg-foundation/engine/pull/801)) Add Service's hash to reply of Deploy API
- ([#849](https://github.com/mesg-foundation/engine/pull/849)) Use base58 to encode service hash.
- ([#860](https://github.com/mesg-foundation/engine/pull/860)) Separate service's configuration from service's dependencies. ([#880](https://github.com/mesg-foundation/engine/pull/880)).
- ([#866](https://github.com/mesg-foundation/engine/pull/866)) Rename service's `volumesfrom` property.
- ([#905](https://github.com/mesg-foundation/engine/pull/905)) Add version to database path to prevent decoding error and loss of data.

#### Added

- ([#535](https://github.com/mesg-foundation/engine/pull/535)) Run MESG with MESG Services! [#567](https://github.com/mesg-foundation/engine/pull/567).
- ([#755](https://github.com/mesg-foundation/engine/pull/755)) Implementation of the MESG Marketplace. The Marketplace allows the distribution and reutilization of MESG Services. Check `mesg-core marketplace` commands. ([#778](https://github.com/mesg-foundation/engine/pull/778)). ([#788](https://github.com/mesg-foundation/engine/pull/788)). ([#810](https://github.com/mesg-foundation/engine/pull/810)). ([#817](https://github.com/mesg-foundation/engine/pull/817)). ([#826](https://github.com/mesg-foundation/engine/pull/826)). ([#828](https://github.com/mesg-foundation/engine/pull/828)). ([#829](https://github.com/mesg-foundation/engine/pull/829)). ([#837](https://github.com/mesg-foundation/engine/pull/837)). ([#843](https://github.com/mesg-foundation/engine/pull/843)). ([#844](https://github.com/mesg-foundation/engine/pull/844)). ([#845](https://github.com/mesg-foundation/engine/pull/845)). ([#853](https://github.com/mesg-foundation/engine/pull/853)). ([#854](https://github.com/mesg-foundation/engine/pull/854)). ([#863](https://github.com/mesg-foundation/engine/pull/863)). ([#864](https://github.com/mesg-foundation/engine/pull/864)). ([#868](https://github.com/mesg-foundation/engine/pull/868)). ([#874](https://github.com/mesg-foundation/engine/pull/874)). ([#883](https://github.com/mesg-foundation/engine/pull/883)). ([#899](https://github.com/mesg-foundation/engine/pull/899)). ([#898](https://github.com/mesg-foundation/engine/pull/898)). ([#897](https://github.com/mesg-foundation/engine/pull/897)). ([#896](https://github.com/mesg-foundation/engine/pull/896)). ([#902](https://github.com/mesg-foundation/engine/pull/902)). ([#901](https://github.com/mesg-foundation/engine/pull/901)). ([#906](https://github.com/mesg-foundation/engine/pull/906)). ([#907](https://github.com/mesg-foundation/engine/pull/907)). ([#908](https://github.com/mesg-foundation/engine/pull/908)). ([#909](https://github.com/mesg-foundation/engine/pull/909)). ([#924](https://github.com/mesg-foundation/engine/pull/924)). ([#926](https://github.com/mesg-foundation/engine/pull/926)). ([#927](https://github.com/mesg-foundation/engine/pull/927)). ([#936](https://github.com/mesg-foundation/engine/pull/936)). ([#938](https://github.com/mesg-foundation/engine/pull/938)). ([#939](https://github.com/mesg-foundation/engine/pull/939)). ([#942](https://github.com/mesg-foundation/engine/pull/942)). ([#943](https://github.com/mesg-foundation/engine/pull/943)).
- ([#757](https://github.com/mesg-foundation/engine/pull/757)) Read `.dockerignore` in dev and deploy commands.
- ([#779](https://github.com/mesg-foundation/engine/pull/779)) Implementation of the MESG Wallet. Check `mesg-core wallet`. ([#807](https://github.com/mesg-foundation/engine/pull/807)). ([#809](https://github.com/mesg-foundation/engine/pull/809)). ([#812](https://github.com/mesg-foundation/engine/pull/812)). ([#937](https://github.com/mesg-foundation/engine/pull/937)).
- ([#781](https://github.com/mesg-foundation/engine/pull/781)) Improve validation of service definition. ([#869](https://github.com/mesg-foundation/engine/pull/869)).

#### Changed

- ([#823](https://github.com/mesg-foundation/engine/pull/823)) Improve commands `service init` and `service gendoc`.
- ([#875](https://github.com/mesg-foundation/engine/pull/875)) Improve JSON parsing error message.
- ([#790](https://github.com/mesg-foundation/engine/pull/790)) Refactor. ([#792](https://github.com/mesg-foundation/engine/pull/792)). ([#816](https://github.com/mesg-foundation/engine/pull/816)). ([#805](https://github.com/mesg-foundation/engine/pull/805)). ([#813](https://github.com/mesg-foundation/engine/pull/813)). ([#839](https://github.com/mesg-foundation/engine/pull/839)). ([#847](https://github.com/mesg-foundation/engine/pull/847)). ([#850](https://github.com/mesg-foundation/engine/pull/850)). ([#852](https://github.com/mesg-foundation/engine/pull/852)). ([#855](https://github.com/mesg-foundation/engine/pull/855)). ([#858](https://github.com/mesg-foundation/engine/pull/858)). ([#867](https://github.com/mesg-foundation/engine/pull/867)). ([#859](https://github.com/mesg-foundation/engine/pull/859)). ([#870](https://github.com/mesg-foundation/engine/pull/870)). ([#871](https://github.com/mesg-foundation/engine/pull/871)). ([#872](https://github.com/mesg-foundation/engine/pull/872)). ([#873](https://github.com/mesg-foundation/engine/pull/873)). ([#881](https://github.com/mesg-foundation/engine/pull/881)). ([#893](https://github.com/mesg-foundation/engine/pull/893)). ([#892](https://github.com/mesg-foundation/engine/pull/892)). ([#891](https://github.com/mesg-foundation/engine/pull/891)). ([#890](https://github.com/mesg-foundation/engine/pull/890)). ([#889](https://github.com/mesg-foundation/engine/pull/889)). ([#888](https://github.com/mesg-foundation/engine/pull/888)). ([#886](https://github.com/mesg-foundation/engine/pull/886)). ([#885](https://github.com/mesg-foundation/engine/pull/885)). ([#884](https://github.com/mesg-foundation/engine/pull/884)). ([#903](https://github.com/mesg-foundation/engine/pull/903)). ([#919](https://github.com/mesg-foundation/engine/pull/919)).

#### Fixed

- ([#771](https://github.com/mesg-foundation/engine/pull/771)) Fix gRPC stream acknowledgement.
- ([#772](https://github.com/mesg-foundation/engine/pull/772)) Improve command logs errors.
- ([#820](https://github.com/mesg-foundation/engine/pull/820)) Fix container package.

## [v0.8.1](https://github.com/mesg-foundation/engine/releases/tag/v0.8.1)

### [Click here to see the release notes](https://forum.mesg.com/t/mesg-core-v0-8-1-release-notes/249)

#### Fixed

- ([#774](https://github.com/mesg-foundation/engine/pull/774)) Update keep alive of client to 5min to prevent spamming the server.

#### Documentation

- ([#762](https://github.com/mesg-foundation/engine/pull/762)) Fix and improve guide start. ([#763](https://github.com/mesg-foundation/engine/pull/763)).

## [v0.8.0](https://github.com/mesg-foundation/engine/releases/tag/v0.8.0)

### [Click here to see the release notes](https://forum.mesg.com/t/mesg-core-v0-8-release-notes/239/2)

#### Added

- ([#690](https://github.com/mesg-foundation/engine/pull/690)) Support service deployments from tarball urls.
- ([#732](https://github.com/mesg-foundation/engine/pull/732)) Support multiple service id or hash for commands `service start` and `service stop`.
- ([#726](https://github.com/mesg-foundation/engine/pull/726)) Add flag to command `start` to force colors in logs of Core.

#### Changed

- ([#734](https://github.com/mesg-foundation/engine/pull/734)) Returns service sid in commands instead of hash.
- ([#724](https://github.com/mesg-foundation/engine/pull/724)) Changed system services deployment system. ([#727](https://github.com/mesg-foundation/engine/pull/727)). ([#725](https://github.com/mesg-foundation/engine/pull/725)). ([#743](https://github.com/mesg-foundation/engine/pull/743)).

#### Fixed

- ([#738](https://github.com/mesg-foundation/engine/pull/738)) Fix stream disconnection because of more than 15min of inactivity. ([#739](https://github.com/mesg-foundation/engine/pull/739)). ([#742](https://github.com/mesg-foundation/engine/pull/742)). ([#744](https://github.com/mesg-foundation/engine/pull/744)).

#### Documentation

- ([#721](https://github.com/mesg-foundation/engine/pull/721)) Move documentation to [dedicated repository](https://github.com/mesg-foundation/docs).

## [v0.7.0](https://github.com/mesg-foundation/engine/releases/tag/v0.7.0)

### [Click here to see the release notes](https://forum.mesg.com/t/mesg-core-v0-7-release-notes/197)

#### Added

- ([#677](https://github.com/mesg-foundation/engine/pull/677)) Stream acknowledgement system. The core notifies client when streams are ready.
- ([#679](https://github.com/mesg-foundation/engine/pull/679)) Add support of repeated parameters to service definition. ([#680](https://github.com/mesg-foundation/engine/pull/680)). ([#684](https://github.com/mesg-foundation/engine/pull/684)).
- ([#682](https://github.com/mesg-foundation/engine/pull/682)) Add support of type Any to service definition. ([#689](https://github.com/mesg-foundation/engine/pull/689)).
- ([#691](https://github.com/mesg-foundation/engine/pull/691)) Add database transaction mechanism to database execution.
- ([#696](https://github.com/mesg-foundation/engine/pull/696)) Add support of nested type definition for type Object.
- ([#704](https://github.com/mesg-foundation/engine/pull/704]) Move go-service to package client/service.

#### Changed

- ([#688](https://github.com/mesg-foundation/engine/pull/688)) Change sid auto-generated prefix.
- ([#699](https://github.com/mesg-foundation/engine/pull/699)) Updated to golang v1.11.4.

#### Fixed

- ([#687](https://github.com/mesg-foundation/engine/pull/687)) Fix execution generated id.
- ([#703](https://github.com/mesg-foundation/engine/pull/703)) Return error when core is not running in command dev and deploy.

#### Removed

- ([#675](https://github.com/mesg-foundation/engine/pull/675)) Remove workflow grpc client.
- ([#693](https://github.com/mesg-foundation/engine/pull/693)) Remove vendor folder.

## [v0.6.0](https://github.com/mesg-foundation/engine/releases/tag/v0.6.0)

### [Click here to see the release notes](https://forum.mesg.com/t/mesg-core-v0-6-release-notes/166)

#### Added

- ([#641](https://github.com/mesg-foundation/engine/pull/641)) Services definition accept env variables. Users can override them on deploy. [#660](https://github.com/mesg-foundation/engine/pull/660). [#666](https://github.com/mesg-foundation/engine/pull/666).
- ([#651](https://github.com/mesg-foundation/engine/pull/651)) Error added in task execution result.

#### Changed

- ([#611](https://github.com/mesg-foundation/engine/pull/611)) Switch to go1.11.
- ([#648](https://github.com/mesg-foundation/engine/pull/672)) Print all service definition in command `service detail`.
- ([#649](https://github.com/mesg-foundation/engine/pull/649)) Lowercase sid.

#### Documentation

- ([#638](https://github.com/mesg-foundation/engine/pull/638)) Fix marketplace link
- ([#643](https://github.com/mesg-foundation/engine/pull/643)) Add instruction to start the core without CLI
- ([#656](https://github.com/mesg-foundation/engine/pull/656)) Show instruction to create manually system services folder
- ([#665](https://github.com/mesg-foundation/engine/pull/665)) Add favicon

## [v0.5.0](https://github.com/mesg-foundation/engine/releases/tag/v0.5.0)

### [Click here to see the release notes](https://forum.mesg.com/t/mesg-core-v0-5-release-notes/136)

#### Breaking Changes

- ([#608](https://github.com/mesg-foundation/engine/pull/608)) Rename "command" property and add "args" property in service definition.

#### Added

- ([#583](https://github.com/mesg-foundation/engine/pull/583)) Add property Sid (Service ID) in service definition file. Allow a service to reuse the same volumes after stopping. [#627](https://github.com/mesg-foundation/engine/pull/627). [#613](https://github.com/mesg-foundation/engine/pull/613). [#619](https://github.com/mesg-foundation/engine/pull/619).

#### Changed

- ([#580](https://github.com/mesg-foundation/engine/pull/580)) Refactor package Daemon.
- ([#588](https://github.com/mesg-foundation/engine/pull/588)) Refactor tests of package container.
- ([#604](https://github.com/mesg-foundation/engine/pull/604)) Improve hash function.
- ([#609](https://github.com/mesg-foundation/engine/pull/609)) Delete all service in parallel in commands.
- ([#615](https://github.com/mesg-foundation/engine/pull/615)) Remove initialization of swarm but display useful error.
- ([#617](https://github.com/mesg-foundation/engine/pull/617)) Improve template of command service gen doc.
- ([#630](https://github.com/mesg-foundation/engine/pull/630)) Rename service id to hash.

#### Fixed

- ([#585](https://github.com/mesg-foundation/engine/pull/585)) Handle gracefully task executions without inputs.
- ([#598](https://github.com/mesg-foundation/engine/pull/598)) Start service dependencies one by one. Solve issue when dependencies request access to same resource.

#### Documentation

- ([#568](https://github.com/mesg-foundation/engine/pull/568)) Update what-is-mesg.md.
- ([#569](https://github.com/mesg-foundation/engine/pull/569)) Update README.md.
- ([#620](https://github.com/mesg-foundation/engine/pull/620)) Add docker swarm init steps to doc.

## [v0.4.0](https://github.com/mesg-foundation/engine/releases/tag/v0.4.0)

### [Click here to see the release notes](https://forum.mesg.com/t/mesg-core-v0-4-0-release-notes/116)

#### Added

- ([#534](https://github.com/mesg-foundation/engine/pull/534)) Access service dependencies based on the name of the dependency through the network.
- ([#555](https://github.com/mesg-foundation/engine/pull/555)) Add more logs on the CLI.

#### Changed

- ([#560](https://github.com/mesg-foundation/engine/pull/560)) Store executions in a database - Fix memory leak [#542](https://github.com/mesg-foundation/engine/pull/542)

#### Fixed

- ([#553](https://github.com/mesg-foundation/engine/pull/553)) UI issue with service execute command.
- ([#552](https://github.com/mesg-foundation/engine/pull/552)) Service dev command returns with an error when needed.
- ([#526](https://github.com/mesg-foundation/engine/pull/526)) Improve container deletion when a service is stopped.
- ([#524](https://github.com/mesg-foundation/engine/pull/524)) Fix sync issue on status send chans and sync issue on gRPC deploy stream sends.

## [v0.3.0](https://github.com/mesg-foundation/engine/releases/tag/v0.3.0)

### [Click here to see the release notes](https://forum.mesg.com/t/mesg-core-v0-3-0-release-notes/88)

#### Added

- ([#392](https://github.com/mesg-foundation/engine/pull/392)) **BREAKING CHANGE.** Add support for `.dockerignore`. Remove support of `.mesgignore` [#498](https://github.com/mesg-foundation/engine/pull/498).
- ([#383](https://github.com/mesg-foundation/engine/pull/383)) New API package. [#386](https://github.com/mesg-foundation/engine/pull/386). [#444](https://github.com/mesg-foundation/engine/pull/444). [#486](https://github.com/mesg-foundation/engine/pull/486). [#488](https://github.com/mesg-foundation/engine/pull/488).
- ([#409](https://github.com/mesg-foundation/engine/pull/409)) Add required validations on service's task, event and output data.
- ([#432](https://github.com/mesg-foundation/engine/pull/432)) Configuration of the CLI's output with `--no-color` and `--no-spinner` flags. Colorize JSON. [#453](https://github.com/mesg-foundation/engine/pull/453). [#480](https://github.com/mesg-foundation/engine/pull/480). [#484](https://github.com/mesg-foundation/engine/pull/484).
- ([#435](https://github.com/mesg-foundation/engine/pull/435)) Command `service logs` accepts multiple dependency filters with multiple use of `-d` flag.
- ([#478](https://github.com/mesg-foundation/engine/pull/478)) Allow multiple core to run on the same computer.
- ([#493](https://github.com/mesg-foundation/engine/pull/493)) Support numbers in service task's key, event's key and output's key
- ([#499](https://github.com/mesg-foundation/engine/pull/499)) Return service's status from API

#### Changed

- ([#371](https://github.com/mesg-foundation/engine/pull/371)) Delegate deployment of Service to Core. [#469](https://github.com/mesg-foundation/engine/pull/469).
- ([#404](https://github.com/mesg-foundation/engine/pull/404)) Change building tool.
- ([#413](https://github.com/mesg-foundation/engine/pull/413)) Improve command `service dev`. [#459](https://github.com/mesg-foundation/engine/pull/459).
- ([#417](https://github.com/mesg-foundation/engine/pull/417)) Service refactoring. [#402](https://github.com/mesg-foundation/engine/pull/402). [#414](https://github.com/mesg-foundation/engine/pull/414). [#454](https://github.com/mesg-foundation/engine/pull/454). [#458](https://github.com/mesg-foundation/engine/pull/458). [#464](https://github.com/mesg-foundation/engine/pull/464). [#472](https://github.com/mesg-foundation/engine/pull/472). [#490](https://github.com/mesg-foundation/engine/pull/490). [#491](https://github.com/mesg-foundation/engine/pull/491).
- ([#419](https://github.com/mesg-foundation/engine/pull/419)) Use Docker volumes for services. [#477](https://github.com/mesg-foundation/engine/pull/477).
- ([#427](https://github.com/mesg-foundation/engine/pull/427)) Refactor package Config
- ([#481](https://github.com/mesg-foundation/engine/pull/481)) Refactor package Database
- ([#485](https://github.com/mesg-foundation/engine/pull/485)) Improve CLI output. [#521](https://github.com/mesg-foundation/engine/pull/521).
- Tests improvements. [#381](https://github.com/mesg-foundation/engine/pull/381). [#384](https://github.com/mesg-foundation/engine/pull/384). [#391](https://github.com/mesg-foundation/engine/pull/391). [#446](https://github.com/mesg-foundation/engine/pull/446). [#447](https://github.com/mesg-foundation/engine/pull/447). [#466](https://github.com/mesg-foundation/engine/pull/466). [#489](https://github.com/mesg-foundation/engine/pull/489). [#501](https://github.com/mesg-foundation/engine/pull/501). [#504](https://github.com/mesg-foundation/engine/pull/504). [#506](https://github.com/mesg-foundation/engine/pull/506).

#### Fixed

- ([#401](https://github.com/mesg-foundation/engine/pull/401)) Gracefully stop gRPC servers.
- ([#429](https://github.com/mesg-foundation/engine/pull/429)) Fix issue when stopping services. [#505](https://github.com/mesg-foundation/engine/pull/505). [#526](https://github.com/mesg-foundation/engine/pull/526).
- ([#476](https://github.com/mesg-foundation/engine/pull/476)) Improve database error handling.
- ([#482](https://github.com/mesg-foundation/engine/pull/482)) Fix Service hash changed when fetching from git.

#### Removed

- ([#410](https://github.com/mesg-foundation/engine/pull/410)) Remove socket server in favor of the TCP server.

#### Documentation

- ([#415](https://github.com/mesg-foundation/engine/pull/415)) Added hall-of-fame to README. Thanks [sergey48k](https://github.com/sergey48k).
- ([#423](https://github.com/mesg-foundation/engine/pull/423)) Fix documentation issue.
- ([#474](https://github.com/mesg-foundation/engine/pull/474)) Documentation/update ux.
- ([#509](https://github.com/mesg-foundation/engine/pull/509)) Add forum link. [#513](https://github.com/mesg-foundation/engine/pull/513).
- ([#510](https://github.com/mesg-foundation/engine/pull/510)) Update ecosystem menu.
- ([#511](https://github.com/mesg-foundation/engine/pull/511)) Update tutorial page.
- ([#512](https://github.com/mesg-foundation/engine/pull/512)) Add sitemap.

## [v0.2.0](https://github.com/mesg-foundation/engine/releases/tag/v0.2.0)

#### Added
- ([#242](https://github.com/mesg-foundation/engine/pull/242)) Add more details in command `mesg-core service validate`
- ([#295](https://github.com/mesg-foundation/engine/pull/295)) Added more validation on the API for the data of `executeTask`, `submitResult` and `emitEvent`. Now if data doesn't match the service file, the API returns an error
- ([#302](https://github.com/mesg-foundation/engine/pull/302)) Possibility to use a config file in ~/.mesg/config.yml
- ([#303](https://github.com/mesg-foundation/engine/pull/303)) Add command `service dev` that build and run the service with the logs
- ([#303](https://github.com/mesg-foundation/engine/pull/303)) Add command `service execute` that execute a task on a service
- ([#316](https://github.com/mesg-foundation/engine/pull/316)) Delete service when stoping the `service dev` command to avoid to keep all the versions of the services.
- ([#317](https://github.com/mesg-foundation/engine/pull/317)) Add errors when trying to execute a service that is not running (nothing was happening before)
- ([#344](https://github.com/mesg-foundation/engine/pull/344)) Add `service execute --data` flag to pass arguments as key=value.
- ([#362](https://github.com/mesg-foundation/engine/pull/362)) Add `tags` list parameter for execution in order to be able to categorize and/or track a specific task execution.
- ([#362](https://github.com/mesg-foundation/engine/pull/362)) Add possibility to listen to results with the specific tag(s)

#### Changed
- ([#282](https://github.com/mesg-foundation/engine/pull/282)) Branch support added. You can now specify your branches with a `#branch` fragment at the end of your git url. E.g.: https://github.com/mesg-foundation/service-ethereum-erc20#websocket
- ([#299](https://github.com/mesg-foundation/engine/pull/299)) Add more user friendly errors when failing to connect to the Core or Docker
- ([#356](https://github.com/mesg-foundation/engine/pull/356)) Use github.com/stretchr/testify package
- ([#352](https://github.com/mesg-foundation/engine/pull/352)) Use logrus logging package

#### Fixed
- ([#358](https://github.com/mesg-foundation/engine/pull/358)) Fix goroutine leaks on api package handlers where gRPC streams are used. Handlers now doesn't block forever by exiting on context cancellation and stream.Send() errors.

#### Removed
- ([#303](https://github.com/mesg-foundation/engine/pull/303)) Deprecate command `service test` in favor of `service dev` and `service execute`

## [v0.1.0](https://github.com/mesg-foundation/engine/releases/tag/v0.1.0)

#### Added
- ([#267](https://github.com/mesg-foundation/engine/pull/267)) Add Command `service gen-doc` that generate a `README.md` in the service based on the informations of the `mesg.yml`
- ([#246](https://github.com/mesg-foundation/engine/pull/246)) Add .mesgignore to excluding file from the Docker build

#### Changed
- ([#247](https://github.com/mesg-foundation/engine/pull/247)) Update the `service init` command to have initial tasks and events
- ([#257](https://github.com/mesg-foundation/engine/pull/257)) Update the `service init` command to fetch for template based on the https://github.com/mesg-foundation/awesome/blob/master/templates.json file but also custom templates by giving the address of the template
- ([#261](https://github.com/mesg-foundation/engine/pull/261)) **BREAKING** More consistancy between the APIs, rename `taskData` into `inputData` for the `executeTask` API

#### Fixed
- ([#246](https://github.com/mesg-foundation/engine/pull/246)) Ignore files during Docker build

## [v0.1.0-beta3](https://github.com/mesg-foundation/engine/releases/tag/v0.1.0-beta3)

#### Added
- ([#246](https://github.com/mesg-foundation/engine/pull/246)) Add .mesgignore to excluding file from the Docker build

#### Changed
- ([#247](https://github.com/mesg-foundation/engine/pull/247)) Update the `service init` command to have initial tasks and events
- ([#257](https://github.com/mesg-foundation/engine/pull/257)) Update the `service init` command to fetch for template based on the https://github.com/mesg-foundation/awesome/blob/master/templates.json file but also custom templates by giving the address of the template
- ([#261](https://github.com/mesg-foundation/engine/pull/261)) **BREAKING** More consistancy between the APIs, rename `taskData` into `inputData` for the `executeTask` API

#### Fixed
- ([#246](https://github.com/mesg-foundation/engine/pull/246)) Ignore files during Docker build

## [v0.1.0-beta2](https://github.com/mesg-foundation/engine/releases/tag/v0.1.0-beta2)

#### Added
- ([#174](https://github.com/mesg-foundation/engine/pull/174)) Add CHANGELOG.md file
- ([#179](https://github.com/mesg-foundation/engine/pull/179)) Add filters for the core API
  - [API] Add `eventFilter` on `ListenEvent` API to get notification when an event with a specific name occurs
  - [API] Add `taskFilter` on `ListenResult` API to get notification when a result from a specific task occurs
  - [API] Add `outputFilter` on `ListenResult` API to get notification when a result returns a specific output
- ([#183](https://github.com/mesg-foundation/engine/pull/183)) Add a `configuration` attribute in the `mesg.yml` file to accept docker configuration for your service
- ([#187](https://github.com/mesg-foundation/engine/pull/187)) Stop all services when the MESG Core stops
- ([#190](https://github.com/mesg-foundation/engine/pull/190)) Possibility to `test` or `deploy` a service from a git or GitHub url
- ([#233](https://github.com/mesg-foundation/engine/pull/233)) Add logs in the `service test` command with service logs by default and all dependencies logs with the `--full-logs` flag
- ([#235](https://github.com/mesg-foundation/engine/pull/235)) Add `ListServices` and `GetService` APIs

#### Changed
- ([#174](https://github.com/mesg-foundation/engine/pull/174)) Update CI to build version based on tags
- ([#173](https://github.com/mesg-foundation/engine/pull/173)) Use official Docker client
- ([#175](https://github.com/mesg-foundation/engine/pull/175)) Changed the struct to use to start docker service
- ([#181](https://github.com/mesg-foundation/engine/pull/181)) MESG Core and Service start and stop functions wait for the docker container to actually run or stop.
- ([#183](https://github.com/mesg-foundation/engine/pull/183)) **BREAKING** Docker image is automatically injected in the `mesg.yml` file for your service. Now `dependencies` attribute is for extra dependencies so for most of service this is not necessary anymore.
- ([#212](https://github.com/mesg-foundation/engine/pull/212)) **BREAKING** Communication from services to core is now done through a token provided by the core
- ([#236](https://github.com/mesg-foundation/engine/pull/236)) CLI only use the API
- ([#234](https://github.com/mesg-foundation/engine/pull/234)) `service list` command now includes the status for every services

#### Fixed
- ([#179](https://github.com/mesg-foundation/engine/pull/179)) [Doc] Outdated documentation for the CLI
- ([#185](https://github.com/mesg-foundation/engine/pull/185)) Fix logs with extra characters when `mesg-core logs`


#### Removed
- ([#234](https://github.com/mesg-foundation/engine/pull/234)) Remove command `service status` in favor of `service list` command that includes status
