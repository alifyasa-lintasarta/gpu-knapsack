# GPU Pod Assignment Example

## Example

For example, we have the following configuration (See [test/h100/initial-state-mixed.yaml](test/h100/initial-state-mixed.yaml)):

```yaml
gpu:
  # There are 3 GPUs
  number: 3
  # Each GPU have 7 part of SM and 8 part of Memory
  # See: https://docs.nvidia.com/datacenter/tesla/mig-user-guide/index.html#id15
  capacity: [7, 8]
  # Initial State of each GPUs
  initialState:
    0: [1g.10gb, 1g.10gb]
    1: [2g.20gb]
    2: []
  # MIG Weights
  mappings:
    1g.10gb: [1, 1] # 1/7 the SM and 1/8 the Memory
    1g.20gb: [1, 2]
    2g.20gb: [2, 2]
    3g.40gb: [3, 4]
    4g.40gb: [4, 4]
    7g.80gb: [7, 8] # 7/7 the SM and 8/8 the Memory
# Can these pods be assigned to the GPUs above?
pods:
  3g.40gb: 1
  1g.20gb: 2
```

First, build the program.

```sh
ubuntu@alifyasa:~/gpu-knapsack$ make
mkdir -p out
go build -o out/app src/*.go
```

Then run the program with the config.

```sh
ubuntu@alifyasa:~/gpu-knapsack$ time ./out/app test/h100/initial-state-mixed.yaml 
GPUs: 3
GPU Capacities: [7 8]
Requested Pods:
  3g.40gb: 1
  1g.20gb: 2

GPU Assignment:
GPU 0: 1g.10gb (existing), 1g.10gb (existing), 1g.20gb (new)
GPU 1: 2g.20gb (existing), 1g.20gb (new)
GPU 2: 3g.40gb (new)

real    0m0.007s
user    0m0.002s
sys     0m0.005s
```
