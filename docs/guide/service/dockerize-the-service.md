# Dockerize the Service

## Why do I need Docker ?

Services run in Docker to provide a sandbox and a normalized environment to remove any side effects that may occur when running on many different machines. See more information on the [Docker website](https://www.docker.com/).

## Steps to be compatible with Docker

* [Create the Dockerfile](#create-the-dockerfile)
* [Add a config](#add-a-configuration-and-dependencies) in your [`mesg.yml`](service-file.md) file, if needed
* [Add dependencies](#add-a-configuration-and-dependencies) in your [`mesg.yml`](service-file.md) file, if needed

## Create the Dockerfile

In order to be compatible with [Docker](https://www.docker.com/), a `Dockerfile` needs to be created in the folder of the service. See the [Dockerfile reference](https://docs.docker.com/engine/reference/builder/).

### Examples

<tabs>
  <tab title="Node" vp-markdown>
    
```bash
FROM node:carbon
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 8080
CMD [ "npm", "start" ]
```

[source](https://nodejs.org/en/docs/guides/nodejs-docker-webapp/)

  </tab>
  <tab title="Go" vp-markdown>

```bash
FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]
```

[source](https://blog.codeship.com/building-minimal-docker-containers-for-go-applications/)

  </tab>
</tabs>

## Add a configuration and dependencies

::: tip
Configuration and dependencies are an advanced feature and your service might not need this. This is totally optional and really depends on your service needs.
:::

Once the Service can run on Docker, [Core](../start-here/core.md) should be able to start it automatically. Update the [`mesg.yml`](service-file.md) file with the config. and optional dependencies the service needs.

The `configuration` key is a Dependency object that will be use to configure the main Docker container of the service. All Dependency attributes are available except image. The attribute `image` will be set automatically when the service is deployed.

If the service requires dependencies to other Docker container, specify them in the `dependencies` map.

### Definitions

| **Attribute** | **Type** | **Description** |
| --- | --- | --- | --- | --- | --- |
| **image** | `String` | The docker image of the Service. Only available for dependencies. |
| **volumes** | `array[string]` | A list of [volumes](https://docs.docker.com/storage/volumes/) that will be mounted in the Service. |
| **ports** | `array[string]` | A list of ports that the Service needs to expose. |
| **command** | `String` | The command to run when the Service starts if not defined in your [Dockerfile](#create-the-dockerfile). |
| **volumeFrom** | `array[string]` | List of dependencies' names to mount a volume from. |

### Example

```yaml
name: serviceX
tasks: {}
events: {}
configuration:
  command: "node start"
dependencies:
  serviceToConnectWith:
    image: "..."
    volumes:
      - "/tmp"
    ports:
      - "1234"
```

::: tip Get Help
You need help ? Check out the <a href="https://forum.mesg.com" target="_blank">MESG Forum</a>.