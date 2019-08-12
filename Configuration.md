# Azure Devops Kubernetes Manager Configuration


## Label Selector Template Values

| Field            | Description                 | When Populated                                                                                               |
| ---------------- | --------------------------- | ------------------------------------------------------------------------------------------------------------ |
| `.PullRequestID` | The ID of the pull request. | Service Hook event types `git.pullrequest.created`, `git.pullrequest.merged`, and `git.pullrequest.updated`. |
