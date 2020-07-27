# Getting Started with Python

To install:

```txt
protobuf==3.12.2
grpcio==1.30.0
depscloud_api==0.1.5
```

Usage:

```python
import grpc
from depscloud_api.v1alpha.tracker import tracker_pb2_grpc, tracker_pb2

target = "api.deps.cloud:443"
creds = grpc.ssl_channel_credentials()

channel = grpc.secure_channel(target, creds)

sourceService = tracker_pb2_grpc.SourceServiceStub(channel)
moduleService = tracker_pb2_grpc.ModuleServiceStub(channel)
dependencyService = tracker_pb2_grpc.DependencyServiceStub(channel)
```

![Google Analytics](https://www.google-analytics.com/collect?v=1&cid=555&t=event&ec=repo&ea=open&dp=api%2Fpackages%2Fpython&dt=api%2Fpackages%2Fpython&tid=UA-143087272-2)
