# H100 Test Cases

## Assumptions

1. The value of `capacity` is always `[7, 8]`.
2. Weights are always two dimensional.

## Examples

```yaml
gpu:
  number: 2
  capacity: [7, 8]
  initialState:
    0: [3g.40gb]
    1: []
  mappings:
    1g.10gb: [1, 1]
    1g.20gb: [1, 2]
    3g.40gb: [3, 4]
pods:
  3g.40gb: 1
  1g.20gb: 4
```
