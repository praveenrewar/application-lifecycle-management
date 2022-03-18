# Creating Packages, bundling them into a repo and consuming them

Now that our manifests and images referenced by them are immutable we are ready to package them!

# Creating a package
To start off with we will create a directory for our package.

We can have multiple versions of a package sharing the same metadata. Metadata  might include - a shared display name, a list of maintainers, a short description of what the package aims to do, etc.

Let's quickly make a folder for our package define some metadata for our package.

We will create a file `metadata.yml` and the contents will look something like,
```yaml
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: hello-app.corp.com
spec:
  displayName: "hello-app"
  longDescription: "A simple configurable app to demonstrate Carvel's packaging API"
  shortDescription: "Web Server"
  providerName: Carvel
  maintainers:
    - name: Soumik Majumder
    - name: Praveen Rewar

```
This metadata will be referenced by all versions of our package.

Lets make a directory for the first version of our package,
```bash
mkdir 1.0.0
cd 1.0.0
```

We can now define a package in file `package.yml`.

```yaml
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: hello-app.corp.com.1.0.0
spec:
  refName: hello-app.corp.com
  version: 1.0.0
  template:
    spec:
      fetch:
        - imgpkgBundle:
            image: index.docker.io/100mik/hello-app@sha256:29c02895e51a0157ff844afd97a8ccd42a7ba0dd2e89bf5f9c6a668e17482ccb
      template:
        - ytt:
            paths:
              - config/
        - kbld:
            paths:
              - "-"
              - .imgpkg/images.yml
      deploy:
        - kapp: {}
  valuesSchema:
    openAPIv3:
      title: hello-app.corp.com.1.0.0 values schema
      examples:
        - user_name: stranger
        - svc_port: 8081
      properties:
        user_name:
          type: string
          description: The name the app uses to greet the user
          default: stranger
        svc_port:
          type: integer
          description: The port exposed by the service
          default: 8081
```
This is pretty similar to what we defined to our App CR, with a few minor differences.
- The fetch section refers to the `imgpkg` bundle we created in the previous chapter.
- The template section includes a `ytt` and a `kbld` stage.
    - We ask `ytt` to template and render manifests in the `/config` directory
    - We indicate that kbld should use output provided by `ytt` and resolve images used in it by using the path "-", we also tell `kbld` where it can find the image lock file.
- We deploy using `kapp` as we did for the App CR.
- The values schema defines the values that can be configured by the user, pretty similar to what we defined in our `values-schema.yml` file in cahpter-1.

With that, our package is ready, let's go back to the chapter-4 directory and try deploying it using `kapp`.
```bash
cd ../..
kapp deploy -a test-pkg -f package
```
We can see that `kapp` says that it will create Package and PackageMetadata resources. We can go ahead and say that we do not want to create these resources, as we will add this Package to a repository before consunming it. (If you have confirmed the changes you can delete them by running `kapp delete -a test-pkg`)

We could definitely add individual packages using `kapp`, however we putting it into a repository makes it easier for our users to consume packages.

## Creating a packaage repository

A package repository is an `imgpkg` bundle structured in a certain manner. It looks something like this,
```
repository/
└── .imgpkg/
    └── images.yml
└── packages/
    └── hello-app.corp.com
        └── metadata.yml
        └── 1.0.0.yml
```

Even though we have one version of our package right now, a bundle might include multiple versions.

```
repository/
└── .imgpkg/
    └── images.yml
└── packages/
    └── hello-app.corp.com
        └── metadata.yml
        └── 1.0.0.yml
        └── 1.2.0.yml
        └── 2.0.1.yml
```

Let's quickly create one for our package, do make sure that you are in the `/chapter-4` directory before starting off.

Let's create the folder for it,
```bash
$ mkdir repository
```
Create the folder structure above,
```bash
mkdir repository/.imgpkg 
mkdir repository/packages
mkdir repository/packages/hello-app.corp.com
```
Copy the required files from our package into the correct folder,
```bash
cp package/metadata.yml repository/packages/hello-app.corp.com/metadata.yml
cp package/1.0.0/package.yml repository/packages/hello-app.corp.com/1.0.0.yml
```

We will generate the lock file for the images used using `kbld`.
```bash
kbld -f repository/packages --imgpkg-lock-output repository/.imgpkg/images.yml
```

Our PackageRepository is ready to be bundled and pushed!
If you havve made it so far the `/repository` directory should have a file structure similar to the one we discussed above. You can also sneak a peek at `/final/repository` to see what it looks like.

Let us now bundle and push our PackageRepository using `imgpkg` 🚀
```bash
imgpkg push -b 100mik/hello-app-repo -f .
dir: .
dir: .imgpkg
file: .imgpkg/images.yml
dir: packages
dir: packages/hello-app.corp.com
file: packages/hello-app.corp.com/1.0.0.yml
file: packages/hello-app.corp.com/metadata.yml
Pushed 'index.docker.io/100mik/hello-app-repo@sha256:912c02a668cb871134bf1e90997fe24b15de0e9a02769d24b12e5fbf0c256bf1'
Succeeded
```

Package consumers can now consume the Packages in the repo easily using `kctrl` by pointing it to the bundle URL.

## Consuming packages
Let us start by adding the packages bundled into the repo to the cluster.
```bash
kctrl package repo add -r hello-repo --url index.docker.io/100mik/hello-app-repo@sha256:912c02a668cb871134bf1e90997fe24b15de0e9a02769d24b12e5fbf0c256bf1
```
`kctrl` waits for the repository to add the bundld packages to the cluster. We can now list the packages on the cluster.
```bash
kctrl package available list
```
We can get more information about the package and view the versions available on the cluster.
```bash
kctrl package available get -p hello-app.corp.com
```
Lets install the package on the cluster,
```bash
kctrl package install -i hello-app -p hello-app.corp.com --version 1.0.0
```
`kctrl` waits for the app to finish reconciling, that is, till the resources bundled into the package are created on the cluster.

We can now check if any `kapp` apps have been created since we asked the package to deploy the resources using kapp.
```bash
kapp list
```
We can see that an app `hello-app-ctrl` has been created on the cluster. We can inspect it to ensure that the resources we were trying to install are up and running.
```bash
kapp inspect -a hello-app-ctrl -t
```
Lets see if our app works by using port-forward.
```bash
kubectl port-forward svc/simple-server-app 8081:8081
```
We can open up `https://localhost:8081` in a browser to verify that the app is up and running.

Now that we have installed the package, lets try and configure it.
We can list configurable values by,
```
kctrl package available get -p hello-app.corp.com/1.0.0 --values-schema
```
Note that we specify a version here as configurable values might change over versions.

Let's configure the username used by the app. Lets create a file `values.yml` with the following content.
```yaml
---
user_name: 100mik
```
(Feel free to add your own name)
```bash
kctrl package intalled update -i hello-app --values-file values.yml
```
Once the installation has reconciled, we can port forward again and verify that the app has upgraded.
```bash
kubectl port-forward svc/simple-server-app 8081:8081
```
Now if we open up the app in our browser windows, it should greet you with the configured username!

## Congratulations!
You now know how to:
- Manage your workloads on you cluster effectively
- Reproduce your dev workflows on the cluster reliably
- Package such workflows, version them and distribute them
- Get up and running with packaged workflows on your cluster!

## Join the Carvel community
We would love to hear how Carvel tools help with you out in your day to day workflows. We are also happy to answer more questions you might have. Stay in touch!

* Join Carvel's slack channel, [#carvel in Kubernetes]({https://kubernetes.slack.com/archives/CH8KCCKA5) workspace, and connect with over 1000+ Carvel users.
* Find us on [GitHub](https://github.com/vmware-tanzu/carvel). Suggest how we can improve the project, the docs, or share any other feedback.
* Attend our Community Meetings, happening every Thursday at 10:30 am PT / 1:30 pm ET. Check out the [Community page](/community/) for full details on how to attend.