Sewan's cloud data center go sdk
================================

- Website: https://www.sewan.fr/

<img src="http://entreprises.smallizbeautiful.fr/logo/Sewan-Communications.jpg" width="500px">

Maintainers
-----------

This sdk is maintained by the Sewan's team.

It is consumed by by Sewan's terraform provider plugin (github.com:SewanDevs/terraform-provider-sewan.git) to communicate with [Sewan's cloud data center](https://www.sewan.fr/cloud-data-center/).

Requirements
------------

-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Developing the sdk
---------------------------

SonarQube Requirements : *< coming soon>*

Run unit test :
```sh
$ cd $GOPATH/src/github.com/SewanDevs/sewan_go_sdk
$ go test go test -cover -coverprofile=c.out
$ go tool cover -html=c.out -o coverage.html
```

Acceptance test with real resources test plan and terraform plugin sdk consumption : *< coming soon>*

Doc
--------------------
Available under doc folder, it contains sequence diagrams and a module diagram.
