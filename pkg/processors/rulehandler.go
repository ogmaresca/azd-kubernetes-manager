package processors

import (
	"fmt"
	"strings"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/config"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/kubernetes"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/templating"
)

// RuleHandler handles Kubernetes rules
type RuleHandler interface {
	Handle(rules config.Rules, args templating.Args) error
}

// RuleHandlerImpl is the default implementation of RuleHandler
type RuleHandlerImpl struct {
	client kubernetes.ClientAsync
}

// NewRuleHandler creates a RuleHandler
func NewRuleHandler(client kubernetes.ClientAsync) RuleHandler {
	return RuleHandlerImpl{client}
}

// Handle executes configuration rules
func (rh RuleHandlerImpl) Handle(rules config.Rules, args templating.Args) error {
	if rules.IsEmpty() {
		logger.Infof("[%s] No rules were defined.", args.ServiceHook.Describe())
		return nil
	}

	var channels []chan error
	for _, rule := range rules.Apply {
		channel := make(chan error, len(rules.Apply)+len(rules.Delete))
		rh.handleApply(rule, args, channel)
		channels = append(channels, channel)
	}

	var errors []string
	for _, channel := range channels {
		err := <-channel
		if err != nil {
			errors = append(errors, fmt.Sprintf("- %s", strings.ReplaceAll(err.Error(), "\n", "\n  ")))
		}
	}

	var err error
	if len(errors) > 0 {
		err = fmt.Errorf("%s", strings.Join(errors, "\n"))
	}

	return err
}

// handleApply executes Apply Resource rules
func (rh RuleHandlerImpl) handleApply(rules config.ApplyResourceRule, args templating.Args, channel chan<- error) {
	channel <- nil
}

// handleDelete executes Delete Resource rules
func (rh RuleHandlerImpl) handleDelete(rules config.DeleteResourceRule, args templating.Args, channel chan<- error) {
	channel <- nil
}
