<img src="docs/images/starboard-logo.png" width="200" alt="Starboard logo">

> Kubernetes-native security tool kit.

[![GitHub Release][release-img]][release]
[![Build Actions][build-action-img]][build-action]
[![License][license-img]][license]

## Table of Contents

- [Abstract](#abstract)
- [Rationale](#rationale)
- [Use Cases](#use-cases)
  - [Starboard for DevOps](#starboard-for-devops)
  - [Starboard for Enterprises](#starboard-for-enterprises)
  - [Starboard for Security Vendors](#starboard-for-security-vendors)
- [Installing](#installing)
  - [From the Binary Releases](#from-the-binary-releases)
  - [From Krew](#from-krew)
  - [From Source (Linux, macOS)](#from-source-linux-macos)
- [Getting Started](#getting-started)
- [Custom Security Resources Definitions](#custom-security-resources-definitions)
  - [Specification](#specification)
  - [Generating Code](#generating-code)
  - [Using Generated Code](#using-generated-code)
- [Starboard CLI](#starboard-cli)
- [kubectl starboard](#kubectl-starboard)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Abstract

Starboard is a Kubernetes-native security tool kit for finding risks in your Kubernetes workloads and environments.
It provides [custom security resources definitions][starboard-crds] and the [Go module][starboard-go-module] to work
alongside a range of existing security tools, allowing for use cases such as these:

- Develop Kubernetes-native security applications such as:
  - admission webhook servers (like [Anchore Image Validator][anchore-image-validator]
    and [Starboard Admission Webhook][starboard-admission-webhook])
  - container security operators (like [Container Security Operator][container-security-operator]
    and [Starboard Security Operator][starboard-security-operator])
  - vulnerability adapters and exporters (like [KubeTrivyExporter][kube-trivy-exporter])
  - Kubernetes audit tools (like [kubeaudit][kubeaudit] and [Polaris][polaris])
  - Kubernetes resources sanitizers and linters (like [Popeye][popeye])
  - [kubectl plugins][kubectl-plugins] to scan workloads early on in the development stage
    (like [kubectl starboard][kubectl-starboard] plugin)
  - webhook servers for integrating with enterprise cloud native artifacts registries
    (like [Starboard Harbor Webhook][starboard-harbor-webhook])
  - webhook servers for integrating with commercial cloud native security platforms
    (like [Starboard Aqua CSP Webhook][starboard-aqua-csp-webhook])
- Extend existing Kubernetes dashboards, such as [Octant][octant] or [OpenShift Console][openshift-console], with
  vulnerability assessment reports
- Implement scoring, health check, and metrics systems for the whole Kubernetes cluster or a single namespace,
  aggregating results from different tools to simplify overall security assessments
- Implement custom security dashboards from scratch

## Rationale

By looking at existing Kubernetes security tools you can quickly realize two things. On one hand they differ in many
ways, i.e. have different capabilities, data models, output sinks, license, maturity level and credibility.
On the other hand, they usually have the same or very similar modus operandi, i.e.:

1. Discover Kubernetes workloads via Kubernetes API or by parsing descriptor YAML files
2. Invoke some type of scanner which finds risks, e.g. execute a [Trivy][trivy] binary executable to find container
   image vulnerabilities, invoke a Go function to check SecurityContext of a given Pod, or evaluate a Pod spec against
   some [Rego][opa-rego] rules.
3. Save risk assessment report somewhere, typically to the standard output or a file. JSON/YAML with a free-style schema
   seems to be an "industry" standard.

It's not easy to deal with the results from these different, standalone Kubernetes security tools. 
With all these heterogeneous data models it's very hard to take advantage of all the features provided by a given tool.
Especially when you want to use a few or all of them.

What if all the Kubernetes security tools spoke the same language that everyone knows and understands?
Similarly to the standardized and well known Pod spec, we could come up with the schema for a *vulnerability*,
a *risk assessment check*, a *black-* or *white-listed vulnerability*, or maybe even a *scanner config*. What if you
could combine the results from different tools to give an easy-to-understand overview of current security status? 
This would allow security vendors to focus on what they do best, whereas others could consume the data in the
homogeneous format.

Project Starboard illustrates how the outputs from different security tools can be stored and combined using native
Kubernetes approaches: 

* Storing results in Kubernetes CRDs that can be queried using the Kubernetes API
* Using Kubernetes Operators to efficiently manage security assessments of different resources within the cluster
* Using Kubernetes Operators to aggregate results, using flexible policies, into Kubernetes-native CRDs 

## Use Cases

### Starboard for DevOps

One idea behind Starboard is to help development and DevOps teams deliver secure and compliant applications from the
get-go. As shown in the figure below, Dave Loper is using [`kubectl`][kubectl] to deploy and test his applications.
Without learning the whole new tool, he can now use a familiar [`kubectl starboard`][kubectl-starboard] plugin interface
to scan container images, which comprise his applications, for potentially dangerous and exploitable vulnerabilities. He
can also look for configuration issues that might affect stability, reliability, and scalability of his deployment. This
makes Dave Loper a new security guard of his organization. What's more, by doing that his organization effectively
implemented the shift left security principle in SDLC.

Sometimes, to better understand the complexity of his applications, Dave is using [Octant][octant], a Kubernetes
introspective and object management platform. With [Starboard Octant plugin][starboard-octant-plugin] we extended the
Octant's capabilities to present vulnerability and configuration audits in user-friendly manner.

![](docs/images/use-cases/starboard-for-devops.png)

### Starboard for Enterprises

Manual scanning through the [`kubectl starboard`][kubectl-starboard] plugin is useful, but it has its limitations: it
doesn't scale well with a huge number of Kubernetes workloads and / or multi-tenant clusters as is the case for
enterprises.

In such cases a more suitable option is to deploy the [Starboard Security Operator][starboard-security-operator], which
constantly monitors Kubernetes-native resources, such as Deployments, and runs appropriate scanners against the
underlying deployment descriptors. The scan reports can be saved as custom resources in the same instance of
[etcd][etcd] used by the Kubernetes cluster running the workloads, or an etcd instance external to the cluster.

Because they are accessible over the Kubernetes API, the vulnerability reports or any other security audits can be used
to build or integrate with dashboards tailored for SRE and Security teams.

The same data can be used by the [Starboard Admission Webhook][starboard-admission-webhook] to accept or reject new
deployments based on security policies put in place, e.g. number of critical vulnerabilities found in a container
image.

Another interesting scenario would be to take advantage of [Starboard Harbor Webhook][starboard-harbor-webhook] or
[Starboard Aqua CSP Webhook][starboard-aqua-csp-webhook] components to import existing vulnerability reports generated
by Harbor or Aqua CSP respectively via Webhook integrations.

![](docs/images/use-cases/starboard-for-enterprises.png)

### Starboard for Security Vendors

Starboard provides a framework for security tool developers and vendors, making it easy to integrate their tooling into Kubernetes in a way that's familiar for users. As a tool developer you can re-use one of the existing CRD definitions, or create a new one. The Starboard code generator saves time and effort on creating the tools to integrate with Kubernetes tooling, for example creating web hooks and plugins for Starboard. 

## Installing

This guide shows how to install the [Starboard CLI][starboard-cli]. Starboard CLI can be installed either from source,
or from pre-built binary releases.

> The Starboard CLI is compatible with [kubectl][kubectl] and is intended as [kubectl plugin][kubectl-plugins],
> but it's perfectly fine to run it as a stand-alone executable. If you rename the `starboard` executable to `kubectl-starboard` and if it's in your path, you can invoke it using `kubectl starboard` (as shown in the examples below). Alternatively, if you want to run as a stand-alone executable, simply call `starboard` in place of `kubectl starboard` in the examples below. 

### From the Binary Releases

Every [release][release] of Starboard provides binary releases for a variety of operating systems. These
binary versions can be manually downloaded and installed.

1. Download your [desired version][release]
2. Unpack it (`tar -zxvf starboard_darwin_x86_64.tar.gz`)
3. Find the `starboard` binary in the unpacked directory, and move it to its desired destination
   (`mv starboard_darwin_x86_64/starboard /usr/local/bin/starboard`)

From there, you should be able to run Starboard CLI commands: `starboard help`

### From Krew

Once we resolve [#8][issue-8] our intention is to submit Starboard to [krew-index][krew-index] so that if accepted, you'll be able to install starboard with
[Krew][krew] plugins manager:

```
$ kubectl krew install starboard
$ kubectl starboard help
```

### From Source (Linux, macOS)

Building from source is slightly more work, but is the best way to go if you want to test the latest (pre-release)
version of Starboard.

You must have a working Go environment.

```
$ git clone git@github.com:aquasecurity/starboard.git
$ cd starboard
$ make
```

If required, it will fetch the dependencies and cache them. It will then compile `starboard` and place it in
`bin/starboard`.

## Getting Started

For those who're familiar with [`kubectl`][kubectl], the easiest way to get started with Starboard is to use the
[Starboard CLI][starboard-cli] or [`kubectl starboard`][kubectl-starboard] plugin, which allow you to scan any
Kubernetes workload deployed in the cluster.

To begin with, initialize Starboard by sending custom security resources definitions to the Kubernetes API:

```
$ kubectl starboard init
$ kubectl api-resources
NAME                              SHORTNAMES   APIGROUP                       NAMESPACED   KIND
ciskubebenchreports               kubebench    aquasecurity.github.io         false        CISKubeBenchReport
configauditreports                             aquasecurity.github.io         true         ConfigAuditReport
kubehunterreports                 kubehunter   aquasecurity.github.io         false        KubeHunterReport
vulnerabilities                   vulns,vuln   aquasecurity.github.io         true         Vulnerability
...
```

The `init` subcommand also creates the `starboard` namespace, in which Starboard executes Kubernetes Jobs to perform
scans.

Next, let's assume there's the `nginx` Deployment in the `dev` namespace:

```
$ kubectl create deployment nginx --image=nginx:1.16 --namespace=dev
```

You can easily find its containers' images vulnerabilities with the following command:

```
$ kubectl starboard find vulnerabilities -n dev deploy/nginx --namespace=dev
```

To retrieve the latest vulnerabilities reports use the standard `kubectl get` command:

```
$ kubectl get vulnerabilities -n dev -o yaml \
    -l starboard.resource.kind=Deployment,starboard.resource.name=nginx
```

Starboard relies on labels and label selectors to associate vulnerability reports with the specified Deployment.
For Deployments with *N* containers Starboard creates *N* instances of `vulnerabilities.aquasecurity.github.io`
resources. Each instance has the `starboard.container.name` label to associate it with a particular container's image.

In this example, the `nginx` deployment has a single container called `nginx`, hence only one instance of the
`vulnerabilities.aquasecurity.github.io` resource is created with the label `starboard.container.name=nginx`

For those for whom the CLI interface is not enough, we've implemented the
[Octant Starboard plugin][starboard-octant-plugin] to display the same vulnerability reports in the Octant's interface.

<p align="center">
  <img src="docs/images/getting-started/octant-plugin-vulnerabilities.png">
</p>

Check the plugin's repository for the installation instructions.

To learn more about the available Starboard commands, use `kubectl starboard help` or type a command followed by the
`-h` flag:

```
$ kubectl starboard kube-bench -h
```

## Custom Security Resources Definitions

[Custom resources][k8s-custom-resources] (CR) is one of the central extension mechanisms used throughout the Kubernetes
ecosystem. Custom resources can be used for small, in-house configuration or data objects without any corresponding
controller logic. But they may also play a central role in projects built on top of Kubernetes that want to offer
a Kubernetes-native API experience.

This project houses CustomResourceDefinitions (CRDs) related to security and compliance checks along with the code
generated by Kubernetes [code generators][k8s-code-generator] to write such custom resources in a natural way.

| NAME                                           | SHORTNAMES   | APIGROUP               | NAMESPACED |  KIND              |
| ---------------------------------------------- | ------------ | ---------------------- | ---------- | ------------------ |
| [vulnerabilities][vulnerabilities-crd]         | vulns,vuln   | aquasecurity.github.io | true       | Vulnerability      |
| [ciskubebenchreports][ciskubebenchreports-crd] | kubebench    | aquasecurity.github.io | false      | CISKubeBenchReport |
| [kubehunterreports][kubehunterreports-crd]     | kubehunter   | aquasecurity.github.io | false      | KubeHunterReport   |
| [configauditreports][configauditreports-crd]   |              | aquasecurity.github.io | true       | ConfigAuditReport  |

### Specification

> **TODO:** Link to Custom Security Resources Specification

### Generating Code

Code generators are used a lot in the implementation of native Kubernetes resources, and we're using the very same
generators here. This project follows the patterns of [k8s.io/sample-controller][k8s-sample-controller], which is
a blueprint for many controllers built in Kubernetes itself.

The code generation starts with:

```
$ go mod vendor
$ export GOPATH="$(go env GOPATH)"
$ ./hack/update-codegen.sh
```

In addition, there is a second script called `./hack/verify-codegen.sh`. This script calls the
`./hack/update-codegen.sh` script and checks whether anything changed, and then it terminates with a nonzero return
code if any of the generated files is not up-to-date. We're running it as a step in the CI/CD pipeline.

### Using Generated Code

An instance of a client set can be created with the `NewForConfig` helper function. This is analogous to the client sets
for core Kubernetes resources. The following listings shows how to create an instance of the
`vulnerabilities.aquasecurity.github.io` resource and send it to the Kubernetes API.

```go
package main

import (
	"log"
	"os"
	"time"

	"k8s.io/client-go/tools/clientcmd"

	starboard "github.com/aquasecurity/starboard/pkg/apis/aquasecurity/v1alpha1"
	starboardapi "github.com/aquasecurity/starboard/pkg/generated/clientset/versioned"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run(_ []string) (err error) {
	config, err := clientcmd.BuildConfigFromFlags("", "~/.kube/config")
	if err != nil {
		return
	}
	client, err := starboardapi.NewForConfig(config)
	if err != nil {
		return
	}

	vulnerability := &starboard.Vulnerability{
		ObjectMeta: meta.ObjectMeta{
			Name:      "a2a6b603-97b4-4e5d-bbcd-404723c4177a",
			Namespace: "dev",
			Labels: map[string]string{
				"starboard.resource.kind":  "Deployment",
				"starboard.resource.name":  "nginx",
				"starboard.container.name": "nginx",
			},
			Annotations: map[string]string{
				"starboard.history.limit": "10",
				"starboard.image.digest":  "sha256:72c42ed48c3a2db31b7dafe17d275b634664a708d901ec9fd57b1529280f01fb",
			},
		},
		Report: starboard.VulnerabilityReport{
			Scanner: starboard.Scanner{
				Name:    "Trivy",
				Vendor:  "Aqua Security",
				Version: "0.5.2",
			},
			Artifact: starboard.Artifact{
				Repository: "library/nginx",
				Digest:     "sha256:72c42ed48c3a2db31b7dafe17d275b634664a708d901ec9fd57b1529280f01fb",
				Tag:        "1.16",
				MimeType:   "application/vnd.docker.distribution.manifest.v2+json",
			},
			Summary: starboard.VulnerabilitySummary{
				CriticalCount: 0,
				HighCount:     0,
				MediumCount:   1,
				LowCount:      0,
				UnknownCount:  0,
			},
			Vulnerabilities: []starboard.VulnerabilityItem{
				{
					VulnerabilityID:  "CVE-2019-1549",
					Resource:         "openssl",
					Severity:         starboard.SeverityMedium,
					InstalledVersion: "1.1.1c-r0",
					FixedVersion:     "1.1.1d-r0",
					Title:            "openssl: information disclosure in fork()",
				},
			},
		},
	}

	_, err = client.AquasecurityV1alpha1().
		Vulnerabilities("dev").
		Create(vulnerability)
	return
}
```

Note that higher-level tools like informers and listers are also generated and available.

## Starboard CLI

Starbord CLI is a single executable binary which can be used to find risks, such as vulnerabilities or insecure Pod
specs, in Kubernetes workloads. By default, the risk assessment reports are stored as
[custom security resources][starboard-crds].

To learn more about the available Starboard CLI commands, run `starboard help` or type a command followed by the
`-h` flag:

```
$ starboard kube-hunter -h
```

## kubectl starboard

The `kubectl starboard` plugin is the same binary as [Starboard CLI](#starboard-cli) but has a different name to follow
the `kubectl` plugins convention. It is compatible with the `kubectl` interface to integrate with `kubectl` plugins
managers such as [krew][krew].

## Troubleshooting

### "starboard" cannot be opened because the developer cannot be verified. (macOS)

Since Starboard CLI is not registered with Apple by an identified developer, if you try to run it for the first time
you might get a warning dialog. This doesn't mean that something is wrong with the release binary, rather macOS can't
check whether the binary has been modified or broken since it was released.

<p align="center">
  <img src="docs/images/troubleshooting/developer-not-verified.png">
</p>

To override your security settings and use the Starboard CLI anyway, follow these steps:

1. In the Finder on your Mac, locate the `starboard` binary.
2. Control-click the binary icon, then choose Open from the shortcut menu.
3. Click Open.

   <p align="center">
     <img src="docs/images/troubleshooting/control-click-open.png">
   </p>

   The `starboard` is saved as an exception to your security settings, and you can use it just as you can any registered
   app.

You can also grant an exception for a blocked Starboard release binary by clicking the Allow Anyway button in the
General pane of Security & Privacy preferences. This button is available for about an hour after you try to run the
Starboard CLI command.

To open this pane on your Mac, choose Apple menu > System Preferences, click Security & Privacy, then click General.

<p align="center">
  <img src="docs/images/troubleshooting/developer-not-verified-remediation.png">
</p>

## Contributing

At this early stage we would love your feedback on the overall concept of Starboard. Over time we'd love to see
contributions integrating different security tools so that users can access security information in standard,
Kubernetes-native ways.

## License

This repository is available under the [Apache License 2.0][license].

[release-img]: https://img.shields.io/github/release/aquasecurity/starboard.svg
[release]: https://github.com/aquasecurity/starboard/releases
[build-action-img]: https://github.com/aquasecurity/starboard/workflows/build/badge.svg
[build-action]: https://github.com/aquasecurity/starboard/actions
[license-img]: https://img.shields.io/github/license/aquasecurity/starboard.svg
[license]: https://github.com/aquasecurity/starboard/blob/master/LICENSE

[starboard-crds]: #custom-security-resources-definitions
[vulnerabilities-crd]: ./kube/crd/vulnerabilities-crd.yaml
[ciskubebenchreports-crd]: ./kube/crd/ciskubebenchreports-crd.yaml
[kubehunterreports-crd]: ./kube/crd/kubehunterreports-crd.yaml
[configauditreports-crd]: ./kube/crd/configauditreports-crd.yaml
[starboard-go-module]: ./pkg
[starboard-cli]: #starboard-cli
[kubectl-starboard]: #kubectl-starboard
[starboard-octant-plugin]: https://github.com/aquasecurity/octant-starboard-plugin
[starboard-security-operator]: https://github.com/aquasecurity/starboard-security-operator
[starboard-admission-webhook]: https://github.com/aquasecurity/starboard-admission-webhook
[starboard-aqua-csp-webhook]: https://github.com/aquasecurity/starboard-aqua-csp-webhook
[starboard-harbor-webhook]: https://github.com/aquasecurity/starboard-harbor-webhook

[k8s-custom-resources]: https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources
[k8s-code-generator]: https://github.com/kubernetes/code-generator
[k8s-sample-controller]: https://github.com/kubernetes/sample-controller

[kubectl]: https://kubernetes.io/docs/reference/kubectl
[kubectl-plugins]: https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins

[octant]: https://github.com/vmware-tanzu/octant
[anchore-image-validator]: https://github.com/banzaicloud/anchore-image-validator
[kube-trivy-exporter]: https://github.com/kaidotdev/kube-trivy-exporter
[container-security-operator]: https://github.com/quay/container-security-operator
[kubeaudit]: https://github.com/Shopify/kubeaudit
[openshift-console]: https://github.com/openshift/console
[popeye]: https://github.com/derailed/popeye
[polaris]: https://github.com/FairwindsOps/polaris
[etcd]: https://etcd.io
[trivy]: https://github.com/aquasecurity/trivy
[opa-rego]: https://www.openpolicyagent.org/docs/latest/policy-language/
[krew]: https://github.com/kubernetes-sigs/krew
[krew-index]: https://github.com/kubernetes-sigs/krew-index

[issue-8]: https://github.com/aquasecurity/starboard/issues/8
