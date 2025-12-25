# mcs-addon helm chart

## Installation type

### global

Only `controller` will be deployed. Other components aren't needed. In other words. The `controller` configuration in `values.yaml` is required and the `dataPlane` block is needed.

The default values.yaml is to install chart with `global` installtationType.

### cluster

All components are needed. Following options in values are required.

- `dataPlane` - stores the broker configurations and will be consumed by `lighthouse-agent` and `controller`. The configuration is from the `global` installation.
- `clusterID` - should be unique across all clusters in clusterset.
- `controller` - the mcs controller configuration.
- `lighthouseAgent` - the `lighthouse-agent` deployment configuration.
- `lighthouseCoredns` - the `lighthouse-coredns` deployment configuration.

example additional `values.yaml`:

```yaml
installationType: cluster
dataPlane:
  namespace: mcs-addon-system
  caCert: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJ2RENDQ\
    VdPZ0F3SUJBZ0lCQURBS0JnZ3Foa2pPUFFRREFqQkdNUnd3R2dZRFZRUUtFeE5rZVc1aGJXbGoKY\
    kdsemRHVnVaWEl0YjNKbk1TWXdKQVlEVlFRRERCMWtlVzVoYldsamJHbHpkR1Z1WlhJdFkyRkFNV\
    GMwTVRBMgpPVGN5TWpBZUZ3MHlOVEF6TURRd05qSTROREphRncwek5UQXpNREl3TmpJNE5ESmFNR\
    Vl4SERBYUJnTlZCQW9UCkUyUjVibUZ0YVdOc2FYTjBaVzVsY2kxdmNtY3hKakFrQmdOVkJBTU1IV\
    1I1Ym1GdGFXTnNhWE4wWlc1bGNpMWoKWVVBeE56UXhNRFk1TnpJeU1Ga3dFd1lIS29aSXpqMENBU\
    VlJS29aSXpqMERBUWNEUWdBRXI2bjVUQ2hDejRNegprM2hSUTNXVXR3akRYN01Eblprd3dQZGpac\
    ExucDlNVnFqRkxBdERFL3gxTG8rK0x4U2YyMXVQY3BjTzdzTjBPCkJ5bSt1SXpqdXFOQ01FQXdEZ\
    1lEVlIwUEFRSC9CQVFEQWdLa01BOEdBMVVkRXdFQi93UUZNQU1CQWY4d0hRWUQKVlIwT0JCWUVGR\
    EMrSHc1RTZiNFkxRi9SeEc4Lzk1UzdJWENBTUFvR0NDcUdTTTQ5QkFNQ0EwY0FNRVFDSUhhRwpHd\
    mEvdzUyRFlFcTRxT3dXdVBLcUptSThsNGpZV0xVVlFqUlhRQnVuQWlCOEhkc0cwUnlxMUYzK01Se\
    nl1REZVClcvWFFwbWQ5OGZkYVZpZ01wQXpndVE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t"
  apiServer: "192.168.31.40:9443/k8s/clusters/local"
  token: "xxxxx"
  insecure: false
clusterID: clusterX
```

### submariner

Only `controller` will be deployed. and Following options are required.

- `dataPlane`
- `clusterID`
- `controller`
