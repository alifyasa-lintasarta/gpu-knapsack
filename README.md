# GPU Pod Fit Simulation

## Usage

Build and run:

```sh
make
time ./out/gpu-knapsack samples/timestamp.yaml
```

## Input

```yaml
gpu:
  number: 2
  capacity: [7, 8]
  mappings:
    1g.10gb: [1, 1]
    1g.20gb: [1, 2]
    2g.20gb: [2, 2]
    3g.40gb: [3, 4]
    4g.40gb: [4, 4]
    7g.80gb: [7, 8]

pods:
  - type: 3g.40gb
    addTime: 1
    removeTime: 8
  - type: 1g.10gb
    addTime: 5
  - type: 3g.40gb
    addTime: 3
  - type: 1g.20gb
    addTime: 6
    removeTime: 7
```

## Output

```
GPUs: 2
GPU Capacities: [7 8]
Events: 4
  3g.40gb (addTime=1, removeTime=8)
  1g.10gb (addTime=5)
  3g.40gb (addTime=3)
  1g.20gb (addTime=6, removeTime=7)

Simulation Starting...
========================

Time 1: Added 3g.40gb to GPU 0
  GPU Usage: GPU0[3/7,4/8], GPU1[0/7,0/8]
Time 3: Added 3g.40gb to GPU 0
  GPU Usage: GPU0[6/7,8/8], GPU1[0/7,0/8]
Time 5: Added 1g.10gb to GPU 1
  GPU Usage: GPU0[6/7,8/8], GPU1[1/7,1/8]
Time 6: Added 1g.20gb to GPU 1
  GPU Usage: GPU0[6/7,8/8], GPU1[2/7,3/8]
Time 7: Removed 1g.20gb from GPU 1
  GPU Usage: GPU0[6/7,8/8], GPU1[1/7,1/8]
Time 8: Removed 3g.40gb from GPU 0
  GPU Usage: GPU0[3/7,4/8], GPU1[1/7,1/8]

Final GPU Assignment:
GPU 0: 3g.40gb
GPU 1: 1g.10gb

real	0m0.007s
user	0m0.003s
sys   0m0.003s
```
