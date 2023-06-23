klum - Kubernetes Lazy User Manager
========

klum does the following basic tasks:

* Create/Delete/Modify users
* Easily manage roles associated with users
* Issues kubeconfig files for users to use

This is a very simple controller that just create service accounts under the hood. Properly
configured this should work on any Kubernetes cluster.

## Installation

```sh
kubectl apply -f https://raw.githubusercontent.com/ibuildthecloud/klum/master/deploy.yaml
```

## Usage
 
### Create User

```yaml
kind: User
apiVersion: klum.cattle.io/v1alpha1
metadata:
  name: darren
```

### Download Kubeconfig
```shell script
kubectl get kubeconfig darren -o json | jq .spec > kubeconfig
kubectl --kubeconfig=kubeconfig get all
```
The name of the kubeconfig resource will be the same as the user name

### Delete User
```shell script
kubectl delete user darren
```

### Assign Roles
```yaml
kind: User
apiVersion: klum.cattle.io/v1alpha1
metadata:
  name: darren
spec:
  clusterRoles:
  - view
  roles:
  - namespace: default
    # you can assign cluster roles in a namespace
    clusterRole: cluster-admin
  - namespace: other
    # or assign a role specific to that namespace
    role: something-custom
```

If you don't assign a role a default role will be assigned to the user which is
configured on the controller.  The default value is cluster-admin, so change
that if you want a more secure setup.

### Disable user
```yaml
kind: User
apiVersion: klum.cattle.io/v1alpha1
metadata:
  name: darren
spec:
  enabled: false
```

When the user is reenabled a new kubeconfig with new token will be created.

## Configuration
The controller can be configured as follows.  You will need to edit the deployment and change
then environment variables:

```shell script
GLOBAL OPTIONS:
   --namespace value             Namespace to create secrets and SAs in (default: "klum") [$NAMESPACE]
   --context-name value          Context name to put in Kubeconfigs (default: "default") [$CONTEXT_NAME]
   --server value                The external server field to put in the Kubeconfigs (default: "https://localhost:6443") [$SERVER_NAME]
   --ca value                    The value of the CA data to put in the Kubeconfig [$CA]
   --default-cluster-role value  Default cluster-role to assign to users with no roles (default: "cluster-admin") [$DEFAULT_CLUSTER_ROLE]
```

## Building

`go build`

![](https://media.giphy.com/media/3o7TKGMZHi73yzCumQ/giphy.gif)

## Running

`./bin/klum --kubeconfig=${HOME}/.kube/config`

## License
Copyright (c) 2020 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
