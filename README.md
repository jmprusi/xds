# xdsctl

Communicate with xDS endpoint.

## Usage

~~~
xdsctl [OPTIONS] VERB [VERB] [ARGS]
~~~

## List

~~~
xdsctl -k -s localhost:18000 list -c cluster [cluster]
~~~

Shows:

~~~
CLUSTER        TYPE
cluster-v0-0   EDS
cluster-v0-1   EDS
cluster-v0-2   EDS
cluster-v0-3   EDS
~~~

~~~
xdsctl -k -s localhost:18000 list -c cluster endpoints

~~~

Will show:

~~~
CLUSTER        ENDPOINTS         STATUSES   WEIGHTS
cluster-v0-0   127.0.0.1:18080   UNKNOWN    0
cluster-v0-1   127.0.0.1:18080   UNKNOWN    0
cluster-v0-2   127.0.0.1:18080   UNKNOWN    0
cluster-v0-3   127.0.0.1:18080   UNKNOWN    0
~~~
## Set

endpoint is identified by address, cluser identified by name

xdsctl set cluster weight|type -c cluster WEIGHT[TYPE]

xdctl set endpoint load|weight|health -c cluster -e endpoint load|weight|health


# rm, remove, delete

xdsctl rm cluster -c cluster

xdctl rm endpoint -c cluster -e endpoint

create a cluster if it doesn't already exist: xdsclt set -c cluster CLUSTER [TYPE]
type defauls to EDS

xdsctl set -c cluster endpoint [ADRESS] [HEALTH]