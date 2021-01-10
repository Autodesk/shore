# Shore PRD

## Problem

Developing CI/CD pipelines is a hard task. They need to be resilient, easy to use and extendable pieces of software that are run by another system (AKA Spinnaker, Jenkins, Tekton, ArgoCD, Waypoint, etc...)

When developing, extending or testing pipelines on any 3rd party system (not the developers computer or a VM) there are a few key problems that come up.

* **Local development can be very hard**, due to setup or complexity of the downstream systems.
* **Stages/Pipelines are not testable**, most CI/CD systems don't have a built in testing framework.
* **Collaboration and code reusability are flawed**, Sharing & Testing code requires jumping through hoops or stacking functionality in one centralized location.

## High Level Goals

* **Reduce time to release/deploy of features**
* **Testable pipelines** - promote a culture of testing and resilience.
* **Scaling & Collaboration through code reuse**.
* **Improve Developer Experience** - Reduce friction when developing pipelines
* **Pipeline development best practices** - create a shared mental model for pipeline development (AKA Framework)
* Open source the frameworks and tooling!

## Out Of Scope

* Support For Other CI/CD solutions (`Spinnaker` is the only backend we support, more on this later)

**Out of scope items may be made part of the project's scope in the future*

## People & Roles

* Simon Chammah (Project Owner)
* Eyal Mor (Lead Developer)
* Sergey Liberman (Developer)
* Daniel Kirillov (Developer)

## Context

### Problem Deep Dive

>This is a deep dive into specifics, and may carry many technology specific terms.\
>The document will try to provide a link to each new term as it is encountered.

The space of [CI/CD](https://en.wikipedia.org/wiki/CI/CD) is saturated with solutions. Some are [cloud native](https://en.wikipedia.org/wiki/Cloud_native_computing) specific, while others are more general purpose.

Examples ([CI/CD](https://en.wikipedia.org/wiki/CI/CD)):

* [Jenkins](https://www.jenkins.io/) / [Jenkins-X](https://jenkins-x.io/)
* [Spinnaker](https://spinnaker.io/)
* [ArgoCD](https://argoproj.github.io/argo-cd/)
* [Tekton](https://tekton.dev/)
* [Waypoint](https://www.waypointproject.io/)

All these systems solve CI/CD, by running code in a managed service based on an external event. Each of the example systems solve [event handling](https://en.wikipedia.org/wiki/Event_(computing)), system management and much more in very different ways.

If we treat a pipeline like a [Web API endpoint](https://en.wikipedia.org/wiki/Web_API#Endpoints), for a web request, the pipeline is just a function/method for which the developer doesn't control the means of execution (e.g.: [Step-Functions](https://aws.amazon.com/step-functions/), [Lambda](https://aws.amazon.com/lambda/)).

However, unlike developing Web APIs, the dev, debug & test cycle for CI/CD gets complicated.

Lack of built in debugging & testing tools in these CI/CD solutions hurdles the release of features in a timely manner.

As it stands, there is no formal way to test pipelines, some solutions in the wild are either not available in Autodesk or not ready at all.

<details>

<summary>The need for testing CI/CD solutions is widely known but hasn't been adequately implemented so far, <b>expand for more details</b>.</summary>

CI/CD solutions struggle with testing and best practices.

This is emphasized by open Github and available tooling.

* [Jenkins Unit Testing](https://github.com/jenkinsci/JenkinsPipelineUnit)
* [Tekton Testsing tools issue](https://github.com/tektoncd/pipeline/issues/1289)
* [Tekton debugging tools issue](https://github.com/tektoncd/pipeline/issues/2069)

Spinnaker doesn't have open issues on the matter ([spinnaker GH issues](https://github.com/spinnaker/spinnaker/issues)).

</details>

The next pain point is opinionated, but has a massive affect on developer productivity.

Software engineers prefer controlling their dev environment. Some develop in a Docker container, some install software on their personal machine, and some will develop in the Cloud (on a VM).

Others developers may mix and match.

When developing pipelines the only option is to use the Web-UI provided by the CI/CD system or use a specific CLI that only covers creating and triggering pipelines.

As it stands cli tools (specifically `spin-cli`) need to be extended to support the same level of UX/DX that the web-uis provide. But web-uis are clunky and hard to automate tests around the GUI.

## The `Shore` Framework

`Shore` is a [framework](https://en.wikipedia.org/wiki/Software_framework) to make developing, testing and deploying pipelines.

**(Marketing phrase ahead, you have been warned!)** `shore` allows developers to build pipelines that **deploy** software to safe **shores**.

The framework consists of:

1. `shore-core` - The orchestration layer
2. `shore-cli` - CLI tool that runs on developers machine.
3. `shore standard libraries` - Standard Reusable Libraries.

### Shore Core

#### Goals

* Interface with rendering engines ([Jsonnet](https://jsonnet.org/), [CUELang](https://cuelang.org/), [Dhall-Lang](https://dhall-lang.org/), [HCL2](https://github.com/hashicorp/hcl/tree/hcl2/))
* Interface with CI/CD backend (e.g. `Spinnaker`)
* Shareable and reusable layer for other projects to consume.
* Load 3rd party libraries (if the rendering engine doesn't support it, e.g. `Jsonnet`, `HCL`)
* Provide facilities for Unit/E2E testing

#### Context

`shore-core` is the programable layer for `shore`.

It can be used in other projects and should be easy to interface with.

It provides the **tooling** necessary to build, execute & test pipelines.

### Shore CLI

#### Goals

* Render pipelines in the terminal
* Save pipelines to the chosen backend
* Execute pipelines with specific arguments.
* E2E testing of pipelines.
* Unit testing of custom rendering functionality.

#### Context

`shore-cli` is the user interface for any `shore` project (backed by `shore-core`).

The CLI serves 2 purposes:

1. An easy to work with tool for developing `shore` projects
2. A `shore-core` reference implementation for other project to consume `shore-core` easily.

The CLI coupled with standard libraries should also provide a mental model on how `shore` projects are structured.

### Shore Standard Shared Libraries

#### Goals

* Abstract boilerplate from developers.
* Promote structures for reusable, testable pipelines.

#### Context

Standard libraries are a great way to get projects of the ground quickly.

These libraries are corner stones of any technical solution (e.g. `glibc`, `Python STDLIB`)

Libraries provide common useful functionality that is usually necessary for many projects.

Some specific examples for common boilerplate code in Spinnaker pipeline development:

* `Pipeline`
* `Stage` (`WaitStage`, `JobStage`, `DeployStage`)
* `NestedPipeline`
* `ParallelStage`

Some of these are feature specific data structures, some are nice to have functions.

## User Stories

### Reusable Pipeline Core Developer

#### User Description

This user is developing a reusable standard deployment pipeline. Examples of reusable pipelines: [EC2](https://aws.amazon.com/ec2/), [ECS](https://aws.amazon.com/ecs/), [Kubernetes](https://aws.amazon.com/eks/), [Lambda](https://aws.amazon.com/lambda/), CDN - FE ([Cloudfront](https://aws.amazon.com/cloudfront/) + [S3](https://aws.amazon.com/s3/)).

#### User requirements

I want to be able to develop, debug, test and maintain pipelines easily.

I would like to mock up ideas and validate them out quickly.

Testing is very important for me and my team. This pipeline could be the most essential piece of software for any deployment in our company. We must have a way to formally validate it after every change (bug fix or new feature).

Our pipeline may be very complex. If possible, we would like to have mostly business logic code in our main code repository.

### Generic Pipeline Developer

#### User Description

This user is developing a pipeline for a specific use case for their team.
The standard pipelines don't meet their needs, and changing their deployment process isn't a short term option. This developer must deliver this pipeline as quickly as possible.

#### User requirements

I need to create a deployment pipeline for my team.

This pipeline must be completed in 1 week.

I can't migrate migrate the dev team to the company's standard deployment at the moment, as we need to release features to customers, and the team doesn't have time to plan out a migration to standard tooling.

As an ad-hoc solution, we would like to have this midterm pipeline to support regionalized deployments via the CPN CD infra we already have.

The pipeline will deploy the whole tech stack together. Lambdas, EC2, Infra and Step Functions.

If it is possible, I would like to reuse the existing tested code as much as possible, as I don't have time to test or validate the pipeline, it just has to work sufficiently well.

We do plan to move to standard tooling in the future, but customer features are currently more important.

### Platform developer

#### User Description

This developer is building a deployment platform for their company (surprising!)

The API interface for platform consists of a bunch of files that represent the infra/config as code in a VCS repository.

#### User requirements

I am building a reusable deployment platform for my company.

We plan on deploying any/all AWS/GCP/Azure components our customer base may need.

Because of the complexity of this project, we will be using a custom Web API to manage our customers specific product needs.

I would like to reuse whatever pipelines and/or tools that were already developed in this space.

This platform needs to be 100% formally validated and have regression tests.

We prefer not to use a CLI, as we are writing an API, we would like our formal coding tools to validate the APIs we are consuming from an Interface.

### A Developer Constrained By Jsonnet

#### User Description

This developer is tired of Jsonnet.

#### User requirements

I would like a way to change the JSON rendering language.

It has no formal types, making our developer experience horrible. It's slow and unbearable.

Our project has outgrown the language and it's time to switch.

We would like to have a way to switch runtime languages. Our team has the capacity to implement the rendering layer, but not the backend.

The tooling must allow A/B testing, validating we are 100% backwards compatible and haven't missed a single feature while re-writing our project.
