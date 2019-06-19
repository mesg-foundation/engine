# MESG Engine

[Website](https://mesg.com/) - [Docs](https://docs.mesg.com/) - [Forum](https://forum.mesg.com/) - [Chat](https://discordapp.com/invite/SaZ5HcE) - [Blog](https://medium.com/mesg)


[![GoDoc](https://godoc.org/github.com/mesg-foundation/core?status.svg)](https://godoc.org/github.com/mesg-foundation/core)
[![CircleCI](https://img.shields.io/circleci/project/github/mesg-foundation/core.svg)](https://github.com/mesg-foundation/core)
[![Docker Pulls](https://img.shields.io/docker/pulls/mesg/core.svg)](https://hub.docker.com/r/mesg/core/)
[![Maintainability](https://api.codeclimate.com/v1/badges/86ad77f7c13cde40807e/maintainability)](https://codeclimate.com/github/mesg-foundation/core/maintainability)
[![codecov](https://codecov.io/gh/mesg-foundation/core/branch/dev/graph/badge.svg)](https://codecov.io/gh/mesg-foundation/core)


MESG is a platform for the creation of efficient and easy-to-maintain applications that connect any and all technologies. 

MESG Engine is a communication and connection layer which manages the interaction of all connected services and applications so they can remain lightweight, yet feature packed.

To build an application, follow the [Quick Start Guide](https://docs.mesg.com/guide/quick-start-guide.html)

If you'd like to build Services and share them with the community, go to the [Services](#services) section.

To help us build and maintain MESG Engine, refer to the [Contribute](#contribute) section below.

# Contents

- [Quick Start Guide](#quick-start-guide)
- [Services](#services)
- [Architecture](#architecture)
- [Marketplace](#marketplace)
- [Roadmap](#roadmap)
- [Community](#community)
- [Contribute](#contribute)


# Quick Start Guide

This step-by-step guide will show you how to create an application that gets the ERC20 token balance of an Ethereum account every 10 seconds and send it to a Webhook.

[Check out the Quick Start](https://docs.mesg.com/guide/quick-start-guide.html)

# Services

Services are build and shared on the [Marketplace](https://marketplace.mesg.com/). They are small and reusable pieces of code that, when grouped together, allow developers to build incredible applications with ease.

You can develop a service for absolutely anything you want, as long as it can run inside Docker. Check the [documentation to create your own services](https://docs.mesg.com/guide/service/what-is-a-service.html).

Services implement two types of communication: executing tasks and submitting events.

### Executing Tasks

Tasks have input parameters and outputs with varying data. A task is like a function with inputs and outputs.

Let's take an example of a task that takes 2 number and add them (a sum):

The task accepts as inputs: `a` and `b`.

The task will return the following output: `{ result: xx }`.

Where `result = a + b`

Check out the documentation for more information on [how to create tasks](https://docs.mesg.com/guide/service/listen-for-tasks.html).

### Submitting Events

Services can also submit events to MESG Engine. They allow two-way communication with MESG Engine and Applications.

Let's say the service is an HTTP web server. An event could be submitted when the web server receives a request with the request's payload as the event's data. The service could also submit a specific event for every route of your HTTP API.

For more info on how to create your events, visit the [Emit an Event](https://docs.mesg.com/guide/service/emit-an-event.html) page.


# Architecture

[![MESG Architecture](https://cdn.rawgit.com/mesg-foundation/core/dev/schema1.svg)](https://docs.mesg.com)

# Marketplace

We have a common place to post all community-developed Services and Applications. Check out the [Marketplace](https://marketplace.mesg.com).

# Community

You can find us and other MESG users on the [forum](https://forum.mesg.com). Feel free to check existing posts and help other users of MESG.

Also, be sure to check out the [blog](https://medium.com/mesg) to stay up-to-date with our articles.

# Contribute

Contributions are more than welcome. For more details on how to contribute, please check out the [contribution guide](/CONTRIBUTING.md).

If you have any questions, please reach out to us directly on [Discord](https://discordapp.com/invite/5tVTHJC).

[![0](https://sourcerer.io/fame/antho1404/mesg-foundation/core/images/0)](https://sourcerer.io/fame/antho1404/mesg-foundation/core/links/0)
[![1](https://sourcerer.io/fame/antho1404/mesg-foundation/core/images/1)](https://sourcerer.io/fame/antho1404/mesg-foundation/core/links/1)
[![2](https://sourcerer.io/fame/antho1404/mesg-foundation/core/images/2)](https://sourcerer.io/fame/antho1404/mesg-foundation/core/links/2)
[![3](https://sourcerer.io/fame/antho1404/mesg-foundation/core/images/3)](https://sourcerer.io/fame/antho1404/mesg-foundation/core/links/3)
[![4](https://sourcerer.io/fame/antho1404/mesg-foundation/core/images/4)](https://sourcerer.io/fame/antho1404/mesg-foundation/core/links/4)
[![5](https://sourcerer.io/fame/antho1404/mesg-foundation/core/images/5)](https://sourcerer.io/fame/antho1404/mesg-foundation/core/links/5)
[![6](https://sourcerer.io/fame/antho1404/mesg-foundation/core/images/6)](https://sourcerer.io/fame/antho1404/mesg-foundation/core/links/6)
[![7](https://sourcerer.io/fame/antho1404/mesg-foundation/core/images/7)](https://sourcerer.io/fame/antho1404/mesg-foundation/core/links/7)
