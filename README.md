DC/OS CLI Subcommand for mesos maintenance
==========================================

**CAUTION: This is not battle-hardened yet. USE AT YOUR OWN RISK.** 

# Subcommands

## maintenance machine

### down

### up

### Examples

* down

```sh
$ dcos maintenance machine down --hostname="ap1" --ip="10.10.10.1"
```

```sh
$ dcos maintenance machine down --list="list.csv"
```


* up

```sh
$ dcos maintenance machine up --hostname="ap1" --ip="10.10.10.1"
```

```sh
$ dcos maintenance machine up --list="list.csv"
```

## maintenance schedule

### add

Add machines to schedule with `--start-at` and `--duration` flag. Target machines are supplied with `--file` flag, `CSV` file format. This command does not add machines to current schedules. It creates new `Window` with start time and duration, and add it to schedule. 

 * Thanks to [araddon/dateparse](https://github.com/araddon/dateparse) library, start time can handle very wide range of format. See more on their [examples](https://github.com/araddon/dateparse#extended-example).
 * Duration uses golang `time.Duration` which can handle a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are `ns`, `us`, `ms`, `s`, `m`, `h`.
 * Host/IP list file is handled by golang `encoding/csv`. First column is `hostname` and second is `IP address`. The CSV reader can handle comment line starting with `#` also.
 ```
 # Host, IP
 ap1,10.10.10.1
 ap2,10.10.10.2
 ```

### remove
  
Remove machines from all current schedule windows.  Target machines are supplied with `--file` flag, `CSV` file format.

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
$ dcos maintenance schedule add --start-at="2019-01-01" --duration="200s" --file="list.csv"
Maintenance schedule updated successfully.
```

* remove

```sh
$ dcos maintenance schedule remove --file="list.csv"
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

# Acknowledgement

* Main idea is inspired by https://github.com/dcos-labs/dcos-management
* The client code is heavily adopted from https://github.com/mesosphere/dcos-commons/tree/master/cli
* Thanks to @araddon for great datetime parser for go


