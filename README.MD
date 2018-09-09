# CNI Cleaner API

[![CircleCI](https://circleci.com/gh/jsenon/api-cni-cleanup.svg?style=svg)](https://circleci.com/gh/jsenon/api-cni-cleanup)
[![Go Report Card](https://goreportcard.com/badge/github.com/jsenon/api-cni-cleanup)](https://goreportcard.com/report/github.com/jsenon/api-cni-cleanup)

## Func

- Remove unecessary cni file on kubernetes nodes (Running Pod/Orphaned Pod)
- Exports CNI file number in Prometheus endpoint
- Use Opencensus for tracing and metrics support

### Prerequisite

This application must run inside kubernetes cluster.

It's recommended to run it must on Daemonset in order to access to all nodes where the cni files are located.

### Deployment

Coming soon

### Usage

```sh
CNI File Cleanner and Monitoring

Usage:
  api-cni-cleanup [command]

Available Commands:
  clean       Launch CNI Cleanner
  help        Help about any command
  job         Job Cleanner
  server      Launch CNI Cleanner Server

Flags:
      --api string        External or Internal K8S cluster (default "internal")
      --cnifiles string   Set CNI Folder (default "/var/lib/cni")
      --config string     config file (default is $HOME/.api-cni-cleanup.yaml)
      --debug             Set log level to Debug
  -h, --help              help for api-cni-cleanup
      --jaeger string     Set jaegger collector endpoint (default "http://localhost:14268")

Use "api-cni-cleanup [command] --help" for more information about a command.
```

### Run Server Mode

When you launch application in server mode you will have prometheus metrics:

```sh
api-cni-cleanup server --api APITYPE --cnifiles YOURCNIFOLDER
```

Prometheus metrics are available on:

- http://localhost:9010/metrics

Stats are:

```
# HELP apicnicleanup_number_count number of files
# TYPE apicnicleanup_number_count gauge
apicnicleanup_number_count XX
# HELP apicnicleanup_size_bytes Size of the folder
# TYPE apicnicleanup_size_bytes gauge
apicnicleanup_size_bytes XX
```

Api entrypoint to launch the cleaning is:

- http://localhost:9010/cleanup : Delete unecessary cni files

### Run Job Mode

In order to schedule the cleaning with a cronjob you can launch application in job mode:

Contact API CNI Cleanner
```sh
api-cni-cleanup job --urlserver http://127.0.0.1:9010/cleanup,http://192.168.0.2:9010/cleanup
```

Just specify all you `urlserver` that's represent all you api-cni-cleanup servers 

### (OPTIONAL) Run Batch Mode

This application could be launch on demand. But you need to know where your pod is deployed in order to clean cni file on a targeted node.

Batch Mode cleaning:

```sh
api-cni-cleanup clean --api APITYPE --cnifiles YOURCNIFOLDER
```

### TIPS

Vendoring [constraint](https://github.com/kubernetes/minikube/issues/3037#issuecomment-418384405)

### TODO

- [x] Create New Metric on Prometheus = Count files + Folder Size
- [x] Diff IP files and API PodIP
- [ ] Remove files unused
- [x] Api to launch cleanup
- [ ] k8s deployment: DaemonSet + svc + Cronjob
