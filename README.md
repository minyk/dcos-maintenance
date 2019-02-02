DC/OS CLI Subcommand for mesos maintenance
==========================================

# Subcommands
## maintenance schedule

### Examples

* schedule

```sh
$ dcos maintenance schedule view
Window	Hostname		IP		Start				Duration
0	ap13.dcos.ajway.kr	172.22.30.51	2019-01-11 09:00:00 +0900 KST	168h0m0s
```

* status

```sh
$ dcos maintenance schedule status
Status		Hostname		IP		Offer
Draining	ap13.dcos.ajway.kr	172.22.30.51	none
```

* add 

TBD

* remove

TBD

## maintenance machine

TBD

# Acknowledgement

* Main idea is inspired by https://github.com/dcos-labs/dcos-management
* This repo's client code is heavily adopted from https://github.com/mesosphere/dcos-commons/tree/master/cli


