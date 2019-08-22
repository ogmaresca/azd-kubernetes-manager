# Azure Devops Kubernetes Manager Configuration

## Go Templating

Some configuration fields support [Go templating](https://golang.org/pkg/text/template/). The Go templates use the [Sprig](http://masterminds.github.io/sprig/) library. You can read more about Go templating from the [Helm templating guide](https://helm.sh/docs/chart_template_guide/).

Any templating in the configuration file will be tested on startup with some sample values to validate the go templating.

NOTE: When referring to fields, ALL fields are capitalized and all acronyms (ID, URI, URL, HTML, etc...) in the field name are capitalized.

## Service Hooks

See [the documentation from Microsoft](https://docs.microsoft.com/en-us/azure/devops/service-hooks/services/webhooks?view=azure-devops) on how to set up service hooks. The application will need to be exposed on an Ingress and be publicly available for the service hooks to work. The URL will be `{host}/{basePath}/serviceHooks`.

Service Hook are configured from the top-level field `serviceHooks`. Here is an example configuration that deletes namespaces with the label `azdPullRequests` that has the same value of the Pull Request ID when a Pull Request Updated event is received. This sample configuration will only be executed for:

* completed Pull Requests (statuses `completed` or `abandoned`)
* from a branch that starts with `feature/`
* is getting merged into `master`
* was created by Obama

and will only delete the namespace(s) if the label `azdPreserve` does not exist on the namespace.

``` yaml
serviceHooks:
- event: git.pullrequest.updated
  resourceFilters:
    statuses:
    - completed
    - abandoned
    reasons: []
    projects: []
    repositories: []
    sourceRefs:
    - refs/heads/feature/*
    targetRefs:
    - refs/heads/master
    templates:
    - '{{ .CreatedBy.DisplayName | title | contains "Obama" }}'
  rules:
    delete:
    - apiVersion: v1
      kind: Namespace
      selector:
        matchLabels:
          azdPullRequestId: '{{ .PullRequestID }}'
        matchExpressions:
        - key: azdPreserve
          operator: DoesNotExist
          values: []
```

See [the Service Hooks documentation from Microsoft](https://docs.microsoft.com/en-us/azure/devops/service-hooks/events?toc=%2Fazure%2Fdevops%2Fintegrate%2Ftoc.json&bc=%2Fazure%2Fdevops%2Fintegrate%2Fbreadcrumb%2Ftoc.json&view=azure-devops-2019) for the fields in Service Hook requests.

The configuration also supports summation event types. They are:

| Summation Event Type         | Event Types                                                                                          |
| ---------------------------- | ---------------------------------------------------------------------------------------------------- |
| Builds                       | build.complete                                                                                       |
| Releases                     | ms.vss-release.release-abandoned-event, ms.vss-release.release-created-event                         |
| Release Deployment Approvals | ms.vss-release.deployment-approval-completed-event, ms.vss-release.deployment-approval-pending-event |
| Release Deployments          | ms.vss-release.deployment-completed-event, ms.vss-release.deployment-started-event                   |
| Code                         | tfvc.checkin, git.push                                                                               |
| Pull Requests                | git.pullrequest.created, git.pullrequest.merged, git.pullrequest.updated                             |
| Work Items                   | workitem.created, workitem.deleted, workitem.restored, workitem.updated, workitem.commented          |

### Service Hook Configuration

| Field                           | Description                                                                                | Applicable Event Types                                      |
| ------------------------------- | ------------------------------------------------------------------------------------------ | ----------------------------------------------------------- |
| `event`                         | The Service Hook event type.                                                               | All                                                         |
| `resourceFilters.statuses`      | The resource status(es) to execute on.                                                     | build.complete, Pull Requests                               |
| `resourceFilters.reason`        | The resource reason(s) to execute on.                                                      | build.complete, Releases, Release Deployment Approvals      |
| `resourceFilters.projects`      | The Azure Devops project(s) to execute on.                                                 | All                                                         |
| `resourceFilters.releases`      | The release name(s) to execute on.                                                         | Releases, Release Deployment Approvals                      |
| `resourceFilters.environments`  | The environment(s) to execute on.                                                          | Releases, Release Deployment Approvals, Release Deployments |
| `resourceFilters.approvalTypes` | The approval type(s) to execute on.                                                        | Approvals                                                   |
| `resourceFilters.repositories`  | The Git repositor(y/ies) to execute on.                                                    | git.push, Pull Requests                                     |
| `resourceFilters.sourceRefName` | The source ref(s) to execute on.                                                           | Pull Requests                                               |
| `resourceFilters.targetRefName` | The source ref(s) to execute on.                                                           | Pull Requests                                               |
| `resourceFilters.templates`     | Filters to execute on if the templates.                                                    | All                                                         |
| `continue`                      | If set to true, then continue processing rules after the first matching rule is processed. | All                                                         |
| `rules`                         | The rules to execute for matching service hooks.                                           | All                                                         |


The source and target ref names resource filters are compiled to [POSIX ERE](https://swtch.com/~rsc/regexp/regexp2.html#posix) expressions.

The template resource filters are executed as Go Templates. The value given to the templating engine is the `resource` top-level object on the Service Hook. The template must compile to "true" (case insensitive, whitespace is ignored) for the rule(s) to execute for the service hook.

## Rules

### Configuration

| Field                 | Description                                                                                                                        | Go Templated |
| --------------------- | ---------------------------------------------------------------------------------------------------------------------------------- | ------------ |
| `delete`              | Resources to delete. This is an array of the fields below.                                                                         | No           |
| `delete[].apiVersion` | The API Version of the resources to delete.                                                                                        | No           |
| `delete[].kind`       | The Kind of the resources to delete.                                                                                               | No           |
| `delete[].namespace`  | The namespace of the resources to delete.                                                                                          | Yes          |
| `delete[].selector`   | The [LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.15/#labelselector-v1-meta) to find resources. | Yes          |


### Go Templating Values for Rules

The rule fields with Go templating will be executed using an object with the following fields:

| Field            | Description                           | Type            | Applicable Service Hook Event Types                    |
| ---------------- | ------------------------------------- | --------------- | ------------------------------------------------------ |
| `.EventType`     | The Service Hook Event Type.          | String          | All                                                    |
| `.BuildID`       | The ID of the build.                  | Nullable Int    | build.complete                                         |
| `.BuildNumber`   | The Number of the pull request.       | Nullable String | build.complete                                         |
| `.PullRequestID` | The ID of the pull request.           | Nullable Int    | Pull Requests                                          |
| `.ProjectName`   | The name of the Azure Devops project. | Nullable String | Pull Requests                                          |
| `.ResourceURL`   | The URL of the resource.              | Nullable String | build.complete, Code, Pull Requests, Work Items        |
| `.StartTime`     | The start time of the resource.       | Nullable Time   | build.complete                                         |
| `.FinishTime`    | The finish time of the resource.      | Nullable Time   | build.complete                                         |
| `.Status`        | The status of the resource.           | Nullable String | build.complete, Pull Requests                          |
| `.Reason`        | The reason of the resource.           | Nullable String | build.complete, Releases, Release Deployment Approvals |
| `.ServiceHook`   | The entire Service Hook.              | ServiceHook     | All                                                    |
