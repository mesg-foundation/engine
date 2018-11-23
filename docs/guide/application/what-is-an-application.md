# What is an Application?

Applications and business solutions are built on MESG by attaching an event on one [Service](../service/what-is-a-service.md) to a task on another Service. These can be configured in any order and you can easily create chain reactions or synchronicities of any kind.

Future versions of MESG will not require users to code. Instead, you'll send a configuration file to Core, which is like an order slip, listing all of the events and corresponding tasks you'd like the MESG Network to execute for you.   
  
As long as the [Services](../service/what-is-a-service.md) you want to use have already been connected to the MESG Infrastructure, you will be able to list them within the configuration file. If they haven't been connected to MESG yet, you can connect Services yourself with some coding.

Our software architecture is modeled on event-driven architecture \(EDA\). This will be used in future releases of MESG software. EDA is a software architecture pattern promoting the production, detection, consumption of, and reaction to events.

Events are any new occurrences on a technology. \(e.g. receiving an email, a new deposit, a full battery, the first of the month, a delayed flight, etc.\) With the increased use of digital devices, web services and the internet of things, events are happening around us all the time.  
  
We recommend you build applications to react to events in order to create an application that's quite simple to build, easily-maintainable and compatible with future releases of Core. 

Tasks in your Application are reactions to events \(send an email, notify me on my watch, put the car into standby mode, issue a refund, transfer funds, open a new account, turn on the lights, etc.\).

This is how the configuration file \(like an order slip\) in future releases of MESG is laid out, with events and corresponding tasks. So if you want your application to be compatible with future releases of MESG, we recommend you build your Application based on event-driven architecture while we finish completing the infrastructure. 

By creating an Application based in Event-Driven Architecture, you embrace the philosophy of MESG and make an Application that becomes really easy.

### Source of events

::: tip
The event is the **when** for your application
:::

The source of an event can come from two different parts of your service :

* [Events from services](listen-for-events.md)
* [Outputs from the tasks of services](execute-a-task.md)

### Task to execute

::: tip
The task is the **then** for your application
:::

When one event is coming then the only thing to do is to [execute a task](execute-a-task.md) of the service that you want.

You can find some examples on the [use cases](use-cases.md) page.

