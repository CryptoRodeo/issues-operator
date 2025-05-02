# Issue Operator
:construction: Work In Progress :construction:

This tool offers a Kubernetes-native solution for tracking issues within your cluster.

```bash
kubectl get issues -A
NAMESPACE    NAME                                AGE
team-alpha   database-connection-timeout         6m44s
team-alpha   dependency-update-needed-frontend   6m44s
team-alpha   failed-build-frontend               6m44s
team-alpha   failed-test-api                     6m44s
team-alpha   test-flaky-mobile                   6m44s
team-beta    outdated-dependency-database        6m44s
team-beta    release-failed-production           6m44s
team-beta    resource-quota-exceeded             6m44s
team-delta   failed-pipeline-run                 6m44s
team-delta   permission-config-incorrect         6m44s
team-gamma   build-warning-logging               6m44s
team-gamma   pipeline-outdated                   6m44s

```

## What This Tool Doesn't Do 
This tool won't:
- Monitor cluster resources for errors
- Generate issues automatically based on predefined rules
- Close issues automatically for you

You'll need to manage those tasks yourself in whatever way makes sense for your workflow.

## What This Tool Will Do
- Create a Kubernetes Issue Custom Resource Definition (CRD) that you can manage through the standard Kubernetes API and `kubectl`
- Enable you to automate the creation, updating, and resolution of cluster issue records
- Filter and group issues using labels derived from Issue spec attributes

## Common Use Cases
The Issue Operator can generate and update issues on many common Kubernetes problems, including:

- **Resource Quota Management**: Get warned when namespaces approach CPU, memory, or storage limits
- **Certificate Expiration**: Receive alerts before TLS/SSL certificates expire (7, 30, 60 days)
- **Health Probe Failures**: Know when application health checks are failing
- **Network Issues**: Detect connectivity problems between services
- **Configuration Drift**: Get alerted when live cluster state diverges from Git source (ArgoCD/Flux)
- **Pod Scheduling Issues**: Identify when pods can't schedule due to resource constraints or taints/tolerations
- **Persistent Volume Problems**: Discover storage issues before they affect applications
- **API Deprecations**: Get warned about using deprecated Kubernetes APIs
- **Node Health Issues**: Monitor node conditions like disk pressure or memory pressure
- **Security Policy Violations**: Detect when workloads violate cluster security policies


## Example Usage
### Basic commands
```bash
# List all issues
kubectl get issues -A

# List all issues with detailed output
kubectl get issues -A -o wide

# Get details of a specific issue
kubectl describe issue failed-build-frontend -n team-alpha

# Get the issue in YAML format
kubectl get issue failed-build-frontend -n team-alpha -o yaml
```
### Filtering
#### By Severity
```bash
# Get all critical issues across all namespaces
kubectl get issues -A -l severity=critical

# Get all major issues
kubectl get issues -A -l severity=major

# Get all minor and info issues
kubectl get issues -A -l 'severity in (minor,info)'
```

#### Issue Type
```bash
# Get all build-related issues
kubectl get issues -A -l issueType=build

# Get all dependency and release issues
kubectl get issues -A -l 'issueType in (dependency,release)'

```
#### Resource Type
```bash
# Get all component-level issues
kubectl get issues -A -l resourceType=component

# Get all application-level issues
kubectl get issues -A -l resourceType=application

# Get all workspace-level issues
kubectl get issues -A -l resourceType=workspace

```

#### State
```bash
# Get all active issues
kubectl get issues -A -l state=active

# Get all resolved issues
kubectl get issues -A -l state=resolved

# Filter by Resource Name

# Get all issues related to the frontend-ui component
kubectl get issues -A -l resourceName=frontend-ui

# Get all issues related to the e-commerce-app
kubectl get issues -A -l resourceName=e-commerce-app
```
#### Namespace
```bash
# Get all issues in team-alpha namespace
kubectl get issues -n team-alpha
```

#### Combining filters
```bash
# Get all critical active issues
kubectl get issues -A -l 'severity=critical,state=active'

# Get all build and test issues that are still active
kubectl get issues -A -l 'issueType in (build,test),state=active'

# Get all component-level issues in team-alpha namespace
kubectl get issues -n team-alpha -l resourceType=component

# Get all critical issues related to applications
kubectl get issues -A -l 'severity=critical,resourceType=application'
```

#### Advanced filters
```bash
# Get issues NOT of a certain type (using set-based requirements)
kubectl get issues -A -l 'issueType notin (build,test)'

# Get all non-critical issues that are still active
kubectl get issues -A -l 'severity!=critical,state=active'
```

### Grouping issues
```bash
# Group issues by severity and count them
kubectl get issues -A --sort-by=.metadata.labels.severity -o custom-columns="SEVERITY:.metadata.labels.severity,NAMESPACE:.metadata.namespace,NAME:.metadata.name,STATE:.metadata.labels.state"

# Group issues by type and show related details
kubectl get issues -A --sort-by=.metadata.labels.issueType -o custom-columns="TYPE:.metadata.labels.issueType,SEVERITY:.metadata.labels.severity,NAME:.metadata.name,RESOURCE:.metadata.labels.resourceName,STATE:.metadata.labels.state"

# Group by resource type with count
kubectl get issues -A --sort-by=.metadata.labels.resourceType -o custom-columns="RESOURCE_TYPE:.metadata.labels.resourceType,NAME:.metadata.name,TYPE:.metadata.labels.issueType,SEVERITY:.metadata.labels.severity"
```

### Use field selectors for advanced filtering
```bash
# Find issues with specific text in the title (requires JSONPath)
kubectl get issues -A -o jsonpath='{range .items[?(@.spec.title contains "failed")]}{.metadata.name}{"\n"}{end}'

# Find issues detected after a specific date (requires JSONPath)
kubectl get issues -A -o jsonpath='{range .items[?(@.status.detectedTimestamp > "2025-04-30")]}{.metadata.name}{"\t"}{.status.detectedTimestamp}{"\n"}{end}'
```

### Export filtered results
```bash
# Export all critical issues to a YAML file
kubectl get issues -A -l severity=critical -o yaml > critical-issues.yaml

# Export all active issues for a specific team to JSON
kubectl get issues -n team-alpha -l state=active -o json > team-alpha-active-issues.json
```

### Watch issues in real time
```bash
# Watch for changes to critical issues
kubectl get issues -A -l severity=critical --watch

# Watch for issues being resolved
kubectl get issues -A -l state=resolved --watch
```

### Automation examples
```bash
# Count issues by severity
echo "Critical: $(kubectl get issues -A -l severity=critical --no-headers | wc -l)"
echo "Major: $(kubectl get issues -A -l severity=major --no-headers | wc -l)"
echo "Minor: $(kubectl get issues -A -l severity=minor --no-headers | wc -l)"
echo "Info: $(kubectl get issues -A -l severity=info --no-headers | wc -l)"

# Find related issues using jsonpath
kubectl get issue database-connection-timeout -n team-alpha -o jsonpath='{.spec.relatedIssues}'
```

## Getting Started

### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/issues-operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/issues-operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/issues-operator:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/issues-operator/<tag or branch>/dist/install.yaml
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

