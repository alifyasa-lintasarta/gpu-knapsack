# GPU Pod Assignment Example

## Example

For example, we have the following configuration (See [`config/example-2-7-12.yaml`](config/example-2-7-12.yaml)):

```yaml
gpu:
  # There are 5 GPUs
  number: 5
  # Each GPU have 7 part of SM and 8 part of Memory
  # See: https://docs.nvidia.com/datacenter/tesla/mig-user-guide/index.html#id15
  capacity: [7, 8]
  # MIG Weights
  mappings:
    1g.10gb: [1, 1]  # 1/7 the SM and 1/8 the Memory
    1g.20gb: [1, 2]  # 2/7 the SM and 2/8 the Memory
    3g.40gb: [3, 4]  # 3/7 the SM and 4/8 the Memory
# Can these pods be assigned to the GPUs above?
pods:
  1g.10gb: 12  
  1g.20gb: 7
  3g.40gb: 2
  ```

First, build the program.

```sh
ubuntu@alifyasa:~/gpu-knapsack$ make
mkdir -p out
go build -o out/app src/*.go
```

Then run the program with the config.

```sh
ubuntu@alifyasa:~/gpu-knapsack$ time ./out/app config/example-2-7-12.yaml 
GPUs: 5
GPU Capacities: [7 8]
Requested Pods:
  3g.40gb: 2
  1g.10gb: 12
  1g.20gb: 7

GPU Assignment:
GPU 0: 3g.40gb, 3g.40gb
GPU 1: 1g.20gb, 1g.20gb, 1g.20gb, 1g.20gb
GPU 2: 1g.10gb, 1g.10gb, 1g.20gb, 1g.20gb, 1g.20gb
GPU 3: 1g.10gb, 1g.10gb, 1g.10gb, 1g.10gb, 1g.10gb, 1g.10gb, 1g.10gb
GPU 4: 1g.10gb, 1g.10gb, 1g.10gb

Additional Pod Combinations:
1. 1g.20gb: 2, 1g.10gb: 2
2. 3g.40gb: 1, 1g.10gb: 2
3. 1g.10gb: 4, 1g.20gb: 1
4. 1g.10gb: 6
5. 1g.20gb: 1, 3g.40gb: 1
6. 1g.20gb: 3

real    0m0.010s
user    0m0.004s
sys     0m0.006s
```
