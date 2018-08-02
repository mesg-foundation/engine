# Service file

## Service file

To define a [Service](what-is-a-service.md), you will need to create a specific folder with a `mesg.yml` file that describes its functionalities. This file can contain the following information in a `YAML`syntax:

You can create a default file using the CLI by entering the command:

```bash
mesg-core service init
```

This will create a `mesg.yml` file in your current directory with the following attributes:

| **Attribute** | **Default value** | **Type** | **Description** |
| --- | --- | --- | --- | --- | --- | --- | --- |
| **name** | `""` | `String` | Each Service has a name chosen by the developer. This name is used to identify the service in a nice humanlike way. |
| **description** | `""` | `String` | A description that will be useful to explain the features of your service. |
| **events** | `{}` | `map<id,`[`Event`](emit-an-event.md)`>` | Services must declare a list of events they can emit. Events are actions on a technology the Service is connected to. |
| **tasks** | `{}` | `map<id,`[`Task`](listen-for-tasks.md)`>` | Services declare a list of tasks they can execute. A task is an action that accepts parameters as inputs, executes something on the connected technology, and returns one output to Core, with data. |
| **repository** | `""` | `String` | The url of the repository eg: `https://github.com/org/repo` |
| **configuration** | `{}` | [`Dependency`](dockerize-the-service.md#add-dependencies) | Service can specify one configuration that will be use for the main docker container of the service. |
| **dependencies** | `{}` | `map<id,`[`Dependency`](dockerize-the-service.md#add-dependencies)`>` | Services can specify internal dependencies such as a database, cache or blockchain client. |

You can find an example of `mesg.yml` file [here](https://github.com/mesg-foundation/service-ethereum/blob/master/mesg.yml).



