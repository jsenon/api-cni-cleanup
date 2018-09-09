# CNI Cleaner API


[![CircleCI](https://circleci.com/gh/jsenon/api-cni-cleanup.svg?style=svg)](https://circleci.com/gh/jsenon/api-cni-cleanup)
[![Go Report Card](https://goreportcard.com/badge/github.com/jsenon/api-cni-cleanup)](https://goreportcard.com/report/github.com/jsenon/api-cni-cleanup)

UNDER DEVELOPMENT

## Func

- Remove unecessary cni file on kubernetes nodes (Running Pod/Orphaned Pod)
- Exports CNI file number in Prometheus endpoint
- Support Opentracing


### Prerequisite

Kubernetes cluster

### Launch Development Env

### Run

Batch Mode cleaning:

```sh
api-cni-cleanup clean --api APITYPE --cnifiles YOURCNIFOLDER
```

API and Prometheus Metrics:

```sh
api-cni-cleanup server --api APITYPE --cnifiles YOURCNIFOLDER
```

Contact API CNI Cleanner
```sh
api-cni-cleanup job --urlserver http://127.0.0.1:9010/cleanup
```


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


When you rn with api and prometheus metrics:

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

- http://localhost:9010/cleanup : Delete unecessary cni files

### TIPS

Vendoring [constraint](https://github.com/kubernetes/minikube/issues/3037#issuecomment-418384405)

### TODO

- [x] Create New Metric on Prometheus = Count files + Folder Size
- [x] Diff IP files and API PodIP
- [ ] Remove files unused
- [x] Api to launch cleanup
- [ ] k8s deployment: DaemonSet + svc + Cronjob
