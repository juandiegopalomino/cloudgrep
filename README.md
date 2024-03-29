<!---
Warning: This markdown file is generated, do not edit manually.
Instead edit `README.md.tmpl` and run `make markdown`
-->

<p align="center"><img src="https://user-images.githubusercontent.com/855699/172711194-43330c43-c13e-4b04-9e4a-11eabe8cf850.png" width="250"><br/>
Cloud Asset Explorer</p>

<p align="center">
  <a href="https://github.com/juandiegopalomino/cloudgrep/releases/latest">
    <img src="https://img.shields.io/github/release/run-x/cloudgrep.svg" alt="Current Release" />
  </a>
  <a href="https://github.com/juandiegopalomino/cloudgrep/actions/workflows/checks.yml">
    <img src="https://github.com/juandiegopalomino/cloudgrep/actions/workflows/checks.yml/badge.svg" alt="Tests" />
  </a>

  <a href="http://www.apache.org/licenses/LICENSE-2.0.html">
    <img src="https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg" alt="License" />
  </a>

  <img src="https://img.shields.io/github/commit-activity/w/run-x/cloudgrep.svg?style=plastic" alt="Commit Activity" />

</p>
<p align="center">
<a href="https://join.slack.com/t/cloudgrep/shared_invite/zt-1bl7fewv4-_iwch50U_pP8S3YKBvyyQQ">
    Slack Community
  </a>
  </p>

# What is Cloudgrep?
Cloudgrep is an asset explorer for cloud resources. It shows everything that's being run in the cloud and enables the user to slice and dice these based on tags and properties. It is a UI tool built on open source technologies and runs completely client side (so no data leaves user's machine).
<p align="center">
<img width="820" alt="Screenshot" src="https://user-images.githubusercontent.com/855699/175440360-d6e759d0-ecd6-4a36-889c-b329563979db.png">
</p>



### Why use Cloudgrep?
Cloudgrep's goal is to help engineering teams ensure every resource follows consistent tagging schema. It helps identify missing tags, misspellings and unowned resources. Consistent tagging leads to better cost attribution and faster incident resolution.

Additionally, Cloudgrep is a great tool to visualize all cloud resources in a single place - across regions, accounts and providers.

Try it out by downloading the latest [release](https://github.com/juandiegopalomino/cloudgrep/releases)! For any questions, feel free to join our [Slack workspace](https://join.slack.com/t/cloudgrep/shared_invite/zt-1bl7fewv4-_iwch50U_pP8S3YKBvyyQQ).

<p align="center">
  <a href="https://www.youtube.com/watch?v=Ip-lY9x7bh4"><img width="478" alt="Group 2" src="https://user-images.githubusercontent.com/855699/175441143-5834e9f2-8e23-471a-95be-7b3388d1f455.png"></a>
  </br>
  <span><i>
Demo video</i></span>

</p>

# Features
* Cross-platform support OSX/Linux/Windows 32/64-bit
* Simple installation (distributed as a single binary)
* Zero dependencies
* Supports AWS (If you'd like GCP/Azure support, do let us know by filing an issue!)
* Supports for major AWS resources (like EC2, RDS, S3, and many others - please file an issue if something is missing!)

# Installation

- Mac:
  ```bash
  brew install run-x/tap/cloudgrep
  ```
- Or use [Precompiled binaries](https://github.com/juandiegopalomino/cloudgrep/releases) for supported
operating systems are available.

# Basic Usage

Cloudgrep uses the cloud cloud provider credentials that are available on the user's machine. Make sure to properly set these up (see
[here](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) for AWS).

**NOTE: Cloudgrep only needs ReadOnly credentials -- it creates nothing, it modifies nothing. Moreover, it will
do a best effort scan based on available permissions, so the user does not need to have read access to all resources.**

Once downloaded, just execute the binary to run:
```bash
cloudgrep
```

Cloudgrep will then:

1. Scan the cloud account for global resources and resources on your currently configured AWS region
2. Launch the webapp

## Arguments
You can easily pass cli arguments to cloudgrep for customized behavior, such as multiple/different regions to scan,
what port to serve the webapp on, etc... The cli arguments are all fully documented under the cli's `help` option.
To view documentation for them, simply add the `--help` flag like so:

```bash
cloudgrep --help
```


To specify the AWS regions, you can use the `--regions` flag.
```bash
# will scan using the default AWS credentials for regions us-east-1,us-west-2
cloudgrep --regions us-east-1,us-west-2
```

To use multiple AWS accounts, you can use the `--profiles` flag.
```bash
# will scan 2 accounts using the AWS profiles called "dev" and "prod", using the default region for each profile
cloudgrep --profiles dev,prod

# combine regions and profiles
cloudgrep --profiles dev,prod --regions us-east-1,us-west-2
```

# Advanced Usage
Cloudgrep's behavior can further be configured via a user-inputted config yaml. Configs are then resolved at runtime by
considering the cli arguments, the user-passed config yaml, and the defaults in that order of precedence.

The config yaml can be passed in by using the `-c` or `--config` flag as follows:

```bash
cloudgrep -c my_config.yaml
```

Here is the annotated config file, you only need to define this file if you wish to change the default values.

```yaml
# This config represents all the user-configurable settings for cloudgrep and their default values
# https://github.com/juandiegopalomino/cloudgrep/blob/main/pkg/config/config.yaml

# web represents the specs cloudgrep uses for creating the webapp server
web:
  # host is the host the server is running as
  host: localhost
  # port is the port the server is running in
  port: 8080
  # prefix is the url prefix the server uses
  prefix: "/"
  # skipOpen determines whether to automatically open the webui on startup
  skipOpen: false

# datastore represents the specs cloudgrep uses for creating and/or connecting to the datastore/database used.
datastore:
  # type is the kind of datastore to be used by cloudgrep (currently only supports SQLite)
  type: sqlite
  #  skipRefresh determines whether to refresh the data (i.e. scan the cloud) on startup.
  skipRefresh: false
  # dataSourceName is the Type-specific data source name or uri for connecting to the desired data source
  # default: use memory DB - no data stored locally
  dataSourceName: "file::memory:?cache=shared"
  # use a file DB - the data is persisted on your disk
  # dataSourceName: "~/cloudgrep_data.db"

# providers represents the cloud providers cloudgrep will scan w/ the current credentials
providers:
  - cloud: aws # cloud is the type of the cloud provider (currently only AWS is supported)
  
    # regions is the list of different regions within the cloud provider to scan
    # default: use the default AWS region set in your terminal
    # ex: use one region
    # regions: [us-east-1]
    # ex: use one region + global some resources like S3 Bucket are not region specific
    # regions: [us-east-1, global]
    # ex: use "all" region to scan all available regions
    # regions: [all]

    # use a specific AWS profile
    # profile: dev-AKIAXXXXXXXXXXXXXX

```

# Supported resources

- autoscaling.AutoScalingGroup
- cloudfront.Distribution
- ec2.Address
- ec2.CapacityReservation *(untested)*
- ec2.ClientVpnEndpoint *(untested)*
- ec2.Fleet *(untested)*
- ec2.FlowLogs *(untested)*
- ec2.Image
- ec2.Instance
- ec2.KeyPair
- ec2.LaunchTemplate
- ec2.NatGateway
- ec2.NetworkAcl
- ec2.NetworkInterface
- ec2.ReservedInstance *(untested)*
- ec2.RouteTable
- ec2.SecurityGroup
- ec2.Snapshot
- ec2.SpotInstanceRequest *(untested)*
- ec2.Subnet
- ec2.Volume
- ec2.Vpc
- eks.Cluster
- eks.Nodegroup
- elasticache.CacheCluster
- elb.LoadBalancer
- iam.InstanceProfile
- iam.OpenIDConnectProvider
- iam.Policy
- iam.Role
- iam.SAMLProvider *(untested)*
- iam.User
- iam.VirtualMFADevice
- lambda.Function
- rds.DBCluster
- rds.DBClusterSnapshot
- rds.DBInstance
- rds.DBSnapshot
- route53.HealthCheck
- route53.HostedZone
- s3.Bucket
- sns.Topic
- sqs.Queue

# Development

We love user contributions! Check out our [Dev guide](https://github.com/juandiegopalomino/cloudgrep/blob/main/DEVELOP.md) to get started.

# Important Resources
* [The Cloudgrep Team](https://www.runx.dev/about)
* [Check Out our Blog](https://blog.runx.dev/)
