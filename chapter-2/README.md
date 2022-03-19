# Replicating dev workflows in production clusters reliably

Now that we have our dev workflow set up. We will explore how we can define it declaratively using a custom resource App. This custom resource is defined by kapp-controller which we installed before we started off.

An App CR can be used to define 3 stages:
- Where can it fetch the required resources from?
- How it should template it?
- How it should deploy it?

We will be authoring an App CR which fetches the templated manifests from `/chapter-1/config`, templates it using ytt and then deploys it using `kapp`.

To start off we can first delete the resources we created in chapter-1 as our App CR will be creating them for us!
```bash
kapp delete -a simple-server
```

The App CR will require a service account which it can use to craete resources on the cluster. We can add the required resources and create the service account needed - `default-ns-sa`.
```bash
 kapp deploy -a default-ns-rbac -f https://raw.githubusercontent.com/vmware-tanzu/carvel-kapp-controller/develop/examples/rbac/default-ns.yml
```

We can create a file `app.yaml` which defines this workflow.

```yaml
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: hello-app
  namespace: default
spec:
  serviceAccountName: default-ns-sa
  fetch:
    - git:
        url: https://github.com/praveenrewar/application-lifecycle-management
        ref: origin/main
  template:
    - ytt:
        paths:
        - chapter-1/config
  deploy:
    - kapp: {}
```
We specify the link to a github repository and the branch it should fetch, we tell `ytt` which directory it can find the manifests in and then ask the App CR to deploy it using `kapp`.

We can now create this custom resource on the cluster.
```bash
kapp deploy -a hello-app -f app.yaml
```
After confirming the changes and waiting for the app to reconcile, lets list the kapp apps on the cluster.
```bash
$ kapp list
Target cluster 'https://192.168.64.7:8443' (nodes: minikube)

Apps in namespace 'default'

Name               Namespaces                             Lcs   Lca  
default-ns-rbac    default                                true  8m  
hello-app-ctrl     default                                true  7s  
kc                 (cluster),kapp-controller,kube-system  true  1h  
hello-app          default                                true  9s  

Lcs: Last Change Successful
Lca: Last Change Age

4 apps

Succeeded
```
Here, `hello-app` consists of our App CR where as `hello-app-ctrl` is created by `kapp-controller` using instructions we defined in our custom resource.

We can find the resources deployed by running,
```bash
$ kapp inspect -a hello-app-ctrl -t
Target cluster 'https://192.168.64.9:8443' (nodes: minikube)

Resources in app 'hello-app-ctrl'

Namespace  Name                                     Kind           Owner    Conds.  Rs  Ri  Age  
default    simple-server-app                        Deployment     kapp     2/2 t   ok  -   1m  
default     L simple-server-app-6cb797b95c          ReplicaSet     cluster  -       ok  -   1m  
default     L.. simple-server-app-6cb797b95c-8xttn  Pod            cluster  4/4 t   ok  -   1m  
default     L.. simple-server-app-6cb797b95c-gm685  Pod            cluster  4/4 t   ok  -   1m  
default     L.. simple-server-app-6cb797b95c-rrvdf  Pod            cluster  4/4 t   ok  -   1m  
default    simple-app-config-ver-1                  ConfigMap      kapp     -       ok  -   1m  
default    simple-server-app                        Service        kapp     -       ok  -   1m  
default     L simple-server-app                     Endpoints      cluster  -       ok  -   1m  
default     L simple-server-app-x2clv               EndpointSlice  cluster  -       ok  -   1m  

Rs: Reconcile state
Ri: Reconcile information

9 resources

Succeeded
```

You can use App CRs to replicate your dev workflows reliably in production environments. The App CR syncs up with the source defined in the specification and maakes sure that the cluster is updated accordingly.

The next challenge would be to distribute different versions of such workflows while allowing the consumer to upgrade easily. But before that we need to ensure that the resources being deployed by the App CR are immutable.

Lets dig into more whys and hows in chapter-3!