DC/OS CLI Subcommand for mesos maintenance
==========================================

**CAUTION: This is not battle-hardened yet. USE AT YOUR OWN RISK.** 

# Subcommands

```sh
$ dcos maintenance

usage: dcos maintenance [<flags>] <command>


Flags:
  -h, --help                Show context-sensitive help.
  -v, --verbose             Enable extra logging of requests/responses
      --name="maintenance"  Name of the service instance to query

Commands:
  help [<command> ...]
    Show help.


  machine up [<flags>]
    Stop maintenance for machines

    --list=LIST  Name of a specific file which contains hostname, ip lists for machine up.


  machine down [<flags>]
    Start maintenance for machines

    --list=LIST  Name of a specific file which contains hostname, ip lists for machine down.


  schedule view [<flags>]
    Display current maintenance schedule

    --json  Show raw JSON response instead of user-friendly tree


  schedule add --start-at=START-AT --duration=DURATION --list=LIST
    Add maintenance schedule

    --start-at=START-AT  Start time of this maintenance schedule.
    --duration=DURATION  Duration of maintenance schedule. Can use unit h for hours, m for minutes, s for seconds. e.g: 1h.
    --list=LIST          Name of the list file to update


  schedule remove --list=LIST
    Remove maintenance schedule

    --list=LIST  Name of the list file to update


  status [<flags>]
    Display current maintenance status

    --json  Show raw JSON response instead of user-friendly tree


  loglevel set --duration=DURATION --level=LEVEL [<flags>]
    set log level of agent.

    --agent-id=AGENT-ID  Agent id of the target mesos slave node. If missing, all slave nodes are targeted.
    --duration=DURATION  Duration of modified log level. Can use unit h for hours, m for minutes, s for seconds. e.g: 1h.
    --level=LEVEL        Level of log. 0 to 6.


  loglevel get [<flags>]
    Get logging level of an agent.

    --agent-id=AGENT-ID  Agent ID of the target mesos slave node. If missing, all slave nodes are targeted.


  exhibitor delete --[no-]confirm [<flags>]
    Delete znode on DC/OS Master Exhibitor

    --znode=ZNODE  ZNode path for deletion.
    --confirm      Confirm this deletion.

```

## maintenance machine

### down

### up

### Examples

* down

```sh
$ dcos maintenance machine down --list="list.csv"
```


* up

```sh
$ dcos maintenance machine up --list="list.csv"
```

## maintenance schedule

### add

Add machines to schedule with `--start-at` and `--duration` flag. Target machines are supplied with `--list` flag, `CSV` file format. This command does not add machines to current schedules. It creates new `Window` with start time and duration, and add it to schedule. 

 * Thanks to [araddon/dateparse](https://github.com/araddon/dateparse) library, start time can handle very wide range of format. See more on their [examples](https://github.com/araddon/dateparse#extended-example).
 * Duration uses golang `time.Duration` which can handle a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are `ns`, `us`, `ms`, `s`, `m`, `h`.
 * Host/IP list file is handled by golang `encoding/csv`. First column is `hostname` and second is `IP address`. The CSV reader can handle comment line starting with `#` also.
 ```
 # Host, IP
 ap1,10.10.10.1
 ap2,10.10.10.2
 ```

### remove
  
Remove machines from all current schedule windows.  Target machines are supplied with `--list` flag, `CSV` file format.

### view 

Simply view current schedule.


### Examples

* view

```sh
$ dcos maintenance schedule view
Window	Hostname		IP		Start				Duration
0	ap13.dcos.ajway.kr	172.22.30.51	2019-01-11 09:00:00 +0900 KST	168h0m0s
```

* add 

```sh
$ dcos maintenance schedule add --start-at="2019-01-01" --duration="200s" --list="list.csv"
Maintenance schedule updated successfully.
```

* remove

```sh
$ dcos maintenance schedule remove --list="list.csv"
Maintenance schedule updated successfully.
```

* list.csv

```
# Host,IP
ap1,10.10.10.1
ap2,10.10.10.2
```

## maintenance status

Display status of current schedule.

### Examples

* status

```sh
$ dcos maintenance status
Status		Hostname		IP		Offer
Draining	ap13.dcos.ajway.kr	172.22.30.51	none
```

## maintenance loglevel

### Examples

* Set log level of agent ef71ac72-3f3e-4bd8-9i4a-4db098706e06-S10 to 6 during 10 seconds: 
```sh
$ dcos maintenance loglevel set --agent-id=ef71ac72-3f3e-4bd8-9i4a-4db098706e06-S10 --level=6 --duration=10s

```

# How to

## Build

* Install docker >= 1.13
* Run `make build`(or manually run command in `Makefile`).

## Install

* First of all, see this document for dc/os cli subcommand plugin: https://github.com/dcos/dcos-cli/blob/master/design/plugin.md

* Make direcdories like:
```
$ mkdir -p ~/.dcos/clusters/<cluster-hash>/subcommands/maintenance
$ mkdir -p ~/.dcos/clusters/<cluster-hash>/subcommands/maintenance/env/bin
```

* Copy binary into `env/bin`

```
$ cp build/dcos-maintenance-<OS> ~/.dcos/clusters/<cluster-hash>/subcommands/maintenance/env/bin/dcos-maintenance
```

* Copy `plugin.toml`

```
$ cp plugin.toml ~/.dcos/clusters/<cluster-hash>/subcommands/maintenance/env/
```

* Touch `package.json`

```
$ touch ~/.dcos/clusters/<cluster-hash>/subcommands/maintenance/package.json
```

From this moment, your dcos-cli should recognize maintenance subcommand plugin.

```
$ dcos
Command line utility for the Mesosphere Datacenter Operating
System (DC/OS). The Mesosphere DC/OS is a distributed operating
system built around Apache Mesos. This utility provides tools
for easy management of a DC/OS installation.

Available DC/OS commands:

	auth           	Authenticate to DC/OS cluster
	cluster        	Manage your DC/OS clusters
	config         	Manage the DC/OS configuration file
	experimental   	Manage commands that are under development
	help           	Display help information about DC/OS
	job            	Deploy and manage jobs in DC/OS
	maintenance    	Maintenance DC/OS CLI Module
	marathon       	Deploy and manage applications to DC/OS
	node           	View DC/OS node information
	package        	Install and manage DC/OS software packages
	service        	Manage DC/OS services
	task           	Manage DC/OS tasks

Get detailed command description with 'dcos <command> --help'.
```

# Acknowledgement

* Main idea is inspired by https://github.com/dcos-labs/dcos-management
* The client code is heavily adopted from https://github.com/mesosphere/dcos-commons/tree/master/cli
* Thanks to @araddon for great datetime parser for go


