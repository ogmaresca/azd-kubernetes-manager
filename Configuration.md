# Azure Devops Kubernetes Manager Configuration

## Go Templating

Some configuration fields support [Go templating](https://golang.org/pkg/text/template/). The Go templates use the [Sprig](http://masterminds.github.io/sprig/) library. You can read more about Go templating from the [Helm templating guide](https://helm.sh/docs/chart_template_guide/).

Any templating in the configuration file will be tested on startup with some sample values to validate the go templating.

NOTE: When referring to fields, ALL fields are capitalized and all acronyms (URI, URL, HTML, etc...) are capitalized.

## Service Hooks

Service Hook are configured from the top-level field `serviceHooks`. Here is an example configuration that deletes namespaces with the label `azdPullRequests` that has the same value of the Pull Request ID when a Pull Request Updated event is received. This sample configuration will only be executed for completed Pull Requests (statuses `completed` or `abandoned`) from a branch that starts with `feature/`, is getting merged into `master`, and was created by Obama, and will only delete the namespace(s) if the label `azdPreserve` does not exist on the namespace.

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

## Label Selector Template Values

| Field            | Description                 | When Populated                                                                                               |
| ---------------- | --------------------------- | ------------------------------------------------------------------------------------------------------------ |
| `.PullRequestID` | The ID of the pull request. | Service Hook event types `git.pullrequest.created`, `git.pullrequest.merged`, and `git.pullrequest.updated`. |
