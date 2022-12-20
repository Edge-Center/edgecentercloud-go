EdgeCenter cloud API client
====================================

Command line client to EdgeCenter cloud API.

Installation
------------------------------------

- Clone repo
- Add $GOPATH/bin into $PATH
- Run `make install`.

Getting started
------------------------------------

You will need to set the following env:
```bash
export EC_CLOUD_USERNAME=username
export EC_CLOUD_PASSWORD=secret
export EC_CLOUD_PROJECT=1
export EC_CLOUD_REGION=1
export EC_CLOUD_AUTH_URL=https://api.edgecenter.ru
export EC_CLOUD_API_URL=https://api.edgecenter.ru/cloud
export EC_CLOUD_CLIENT_TYPE=platform
```

* **EC_CLOUD_USERNAME** - username
* **EC_CLOUD_PASSWORD** - user's password
* **EC_CLOUD_PROJECT** - project id
* **EC_CLOUD_REGION** - region id
* **EC_CLOUD_AUTH_URL** - authentication url, you could use the same as in example above
* **EC_CLOUD_API_URL** - api url, you could use the same as in example above
* **EC_CLOUD_CLIENT_TYPE** - client type, you could use the same as in example above

After setting the env, use `-h` key to retrieve all available commands:
```bash
./ec_client -h

   NAME:
   ec_client - EdgeCloud API client

   Environment variables example:

   EC_CLOUD_AUTH_URL=
   EC_CLOUD_API_URL=
   EC_CLOUD_API_VERSION=v1
   EC_CLOUD_USERNAME=
   EC_CLOUD_PASSWORD=
   EC_CLOUD_REGION=
   EC_CLOUD_PROJECT=

USAGE:
   ec_client [global options] command [command options] [arguments...]

VERSION:
   v0.3.00

COMMANDS:
   network        EdgeCloud networks API
   task           EdgeCloud tasks API
   keypair        EdgeCloud keypairs V2 API
   volume         EdgeCloud volumes API
   subnet         EdgeCloud subnets API
   flavor         EdgeCloud flavors API
   loadbalancer   EdgeCloud loadbalancers API
   instance       EdgeCloud instances API
   heat           EdgeCloud Heat API
   securitygroup  EdgeCloud security groups API
   floatingip     EdgeCloud floating ips API
   port           EdgeCloud ports API
   snapshot       EdgeCloud snapshots API
   image          EdgeCloud images API
   region         EdgeCloud regions API
   project        EdgeCloud projects API
   keystone       EdgeCloud keystones API
   quota          EdgeCloud quotas API
   limit          EdgeCloud limits API
   cluster        EdgeCloud k8s cluster commands
   pool           EdgeCloud K8s pool commands
   l7policy       EdgeCloud l7policy API
   router         EdgeCloud router API
   fixed_ip       EdgeCloud reserved fixed ip API
   help, h        Shows a list of commands or help for one command

```
