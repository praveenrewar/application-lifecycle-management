# Application Lifecycle Management

This tutorial goes over how the Carvel toolset addresses some of the problems frequently faced while managing and distributing applications buit on top of Kubernetes.

## Pre-requisites
- An active kubernetes cluster and `kubectl` to interact with cluster resources.
    - For this tutorial, a `minikube` cluster on your system works just fine. We can find installation instructions for the same [here](https://minikube.sigs.k8s.io/docs/start/). (Alternatives such as `kind` work just as well)
    - Installation instructions for `kubectl` can be found [here](https://kubernetes.io/docs/tasks/tools/#kubectl)

Once set up, running `kubectl api-resources` should list all available api-resources on the cluster. And we are good to go 🚀

*Note:* If the user is on Windows we recommend using WSL.

## Setting up the Carvel tools
We will first install the Carvel CLI tools.
We can do this by running,
```bash
wget -O- https://carvel.dev/install.sh | bash
```
Or with curl...
```bash
curl -L https://carvel.dev/install.sh | bash
```

Alternatively if you prefer `brew`,
```bash
brew tap vmware-tanzu/carvel
brew install ytt kbld kapp imgpkg kwt vendir kctrl
```
should get us up and runnning with the tools.

We can now install kapp-controller on the cluster by running,
```bash
kapp deploy -a kc -f https://github.com/vmware-tanzu/carvel-kapp-controller/releases/latest/download/release.yml -y
```
Confirm the changes to the cluster when you are prompted.

Once successful, the output of kapp list should look something like this:
```bash
$ kapp list

Target cluster 'https://192.168.64.7:8443' (nodes: minikube)

Apps in namespace 'default'

Name  Namespaces                             Lcs   Lca  
kc    (cluster),kapp-controller,kube-system  true  39s  

Lcs: Last Change Successful
Lca: Last Change Age

1 apps

Succeeded
```

## Other requirements
- Do take some time to create a Docker Hub account, it will be useful chapter-3 onwards

If you have made it this far, we are good to get started with chapter-1!
