# `congest`: very realistic and accessible network tests

Many network tests utilize a single kubernetes cluster to *simulate* a network.
`congest` utilizes cloud api's instead of using kubernetes and artificial
network latency, packet loss, etc to simulate realistic networking conditions.
Besides enabling the collection of hyper realistic network data, this also has
the side benefit of not having to hand roll a kubernetes deployment. Devs that
focus mainly on software all know how to use a unix based command line, so
accessing virtual machines via ssh, while hacky, is incredibly simple and
powerful. Making arbitrary changes is as trivial as writing a bash script.

## Design

The design of `congest` is reletively simple. The genesis and configuration of
the network is done after the IPs of the nodes in the network are known, so
after pulumi communicates with cloud providers to spin up the specified nodes.
After that, a payload is generated and transfered to each node via `scp`. This
payload is then executed on each node. Currently, this payload literally spins
up `celestia-app` and `txsim` in separate tmux sessions. `txsim` should start
running after about 3 minutes after the network begins to bootstrap.

## Forking

This repo could easily be forked to instead run a different chain all together,
that chain just needs to create a genesis and have some mechanism to creating
txs programmatically. It would be great to make things here a bit more general
so that any cosmos chain dev could run their own hyper realistic network tests.

## Usage

### Running a test

1) Setup and install [pulumi](https://www.pulumi.com/docs/install/)
2) Get a digitalocean token and set it as an environment variable `DIGITALOCEAN_TOKEN`
3) Make sure this token has enough permissions to spin up 100 droplets.
4) Add your ssh key to digitalocean, and set the `SSH_KEY_DO_ID` environment
   variable so that we can tell digitalocean to add that public key to all
   droplets we spin up.

   ```sh
   doctl compute ssh-key list
   ```

5) Run `make test <TestName>` to run a test.

After setting pulumi and the `DIGITALOCEAN_TOKEN` env var, you can run the
following commands to deploy the infrastructure (there can be limits set on the
number of droplets one account can deploy, be sure to set those high enough):

```sh
make test Test100Nodes8MB
```

This test will then proceed to configure and spin up 100 geographically
distributed nodes, bootstrap the network by saving all IPs into an addressbook
that each node is initialized with, and then starting `txsim` on each of the 100
nodes.

### Collecting trace data

By default, all nodes have all the message traces enabled. These can be fetched
via the normal mechanisms supported by the tracer (such as pushing to an s3
bucket), however we can also call

```sh
source download_traces.sh validator-1 consensus_block.jsonl
```

### Cleaning up the test instances

The tests *should* destroy themselves after 30 minutes, however its safest to
check up on this or by manually calling

```sh
make destroy
```

which will ask pulumi to destroy all the nodes. If configured properly and if
the cloud provider's api is working then this should work. It's still a good
idea to check the output of this command to ensure that all resources were
properly destroyed.
