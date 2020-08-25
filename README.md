# Archived Please see https://github.com/davidfischer-ch/terraform-provider-awx

# terraform-provider-awx2

FORKED FROM: https://github.com/mauromedda/terraform-provider-awx

Also Check Out AWX/Tower Provisioner: https://github.com/dahendel/terraform-provisioner-awx

***UNDER DEVELOPMENT ***
terraform-provider-awx is still in developing, and it's roadmap could be found at [here](https://github.com/mauromedda/terraform-provider-awx/blob/master/ROADMAP.md).

### Additions

 - [x] Terraform 0.12.x Supported
 
 - [x] Support to add nested groups for InventoryGroups
 
 - [x] Support for extra credentials in a job template
 
 - [x] Uses go modules  
 
 - [ ] DataSources

Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)


Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.6
-	[Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)
-   [Go Modules](https://blog.golang.org/modules2019)

Building The Provider
---------------------

```bash
go get github.com/dahendel/terraform-provider-awx2
cd $GOPATH/src/github.com/dahendel/terraform-provider-awx2
go build -mod vendor -o terraform-provider-awx
```

Using the provider
----------------------
If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.9+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-awx
...
```

In order to test the provider, you can simply run `make test`.

*Note:* Make sure `AWX_ENDPOINT`, `AWX_USERNAME`, `AWX_PASSWORD` variables are set. The default are:

AWX_ENDPOINT=http://localhost
AWX_USERNAME=admin
AWX_PASSWORD=password

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

Acceptance tests need a fully functional AWX/Tower endpoint.
