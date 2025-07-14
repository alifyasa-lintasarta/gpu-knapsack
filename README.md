# GPU Pod Assignment Example

## Example

For example, we have the following configuration (See [samples/timestamp.yaml](samples/timestamp.yaml)):

```yaml
gpu:
  # There are 3 GPUs
  number: 3
  # Each GPU have 7 part of SM and 8 part of Memory
  # See: https://docs.nvidia.com/datacenter/tesla/mig-user-guide/index.html#id15
  capacity: [7, 8]
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
  - type: 3g.40gb
    addTime: 1
  - type: 1g.10gb
    addTime: 5
  - type: 3g.40gb
    addTime: 3
  - type: 1g.10gb
    addTime: 6
```

First, build the program.

```sh
ubuntu@alifyasa:~/gpu-knapsack$ make
mkdir -p out
go build -o out/gpu-knapsack src/*.go
```

Then run the program with the config.

```sh
ubuntu@ubuntu:~/gpu-knapsack$ time ./out/gpu-knapsack samples/timestamp.yaml 
GPUs: 2
GPU Capacities: [7 8]
Items: 4
  3g.40gb (t=1)
  1g.10gb (t=5)
  3g.40gb (t=3)
  1g.10gb (t=6)

GPU Assignment:
GPU 0: 3g.40gb (t=1), 3g.40gb (t=3)
GPU 1: 1g.10gb (t=5), 1g.10gb (t=6)

real    0m0.010s
user    0m0.007s
sys     0m0.004s
```
