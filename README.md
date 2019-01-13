# ePaxos

**Actor based implementation of the [ePaxos consensus algorithm](http://efficient.github.io/epaxos/).**

## Quick Start

To install this implementation of ePaxos on a system:

```
$ go get github.com/bbengfort/epaxos/...
```

This should install the `epaxos` command on your system. Create a configuration file that defines the peers for the network and other parameters as follows:

```json
{
  "peers": [
    {
      "pid": 1,
      "name": "alpha",
      "ip_address": "127.0.0.1",
      "domain": "localhost",
      "port": 3264
    },
    {
      "pid": 2,
      "name": "bravo",
      "ip_address": "127.0.0.1",
      "domain": "localhost",
      "port": 3265
    },
    {
      "pid": 3,
      "name": "charlie",
      "ip_address": "127.0.0.1",
      "domain": "localhost",
      "port": 3266
    }
  ]
}
```

The configuration file can be stored as a .toml, .json, .yml, .yaml file in the following locations:

- `/etc/epaxos.json`
- `~/.epaxos.json`
- `$(pwd)/epaxos.json`

Or the path to the configuration file can be passed to the command at runtime. To run an ePaxos replica process:

```
$ epaxos serve -n alpha
```

The -n command specifies which peer configures the local replica, by default if no name is specified, then the hostname of the machine is used. To commit a command to the ePaxos log:

```
$ epaxos commit -k "key" -v "value"
```

This commits the command named "key" with the specified "value" to the log. Note that the client is automatically redirected to a leader in a round-robin fashion and requires the same configuration to connect.

