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
		channel := make(chan error)
		go rh.handleApply(rule, args, channel)
		channels = append(channels, channel)
	}

	for _, rule := range rules.Delete {
		channel := make(chan error)
		go rh.handleDelete(rule, args, channel)
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
func (rh RuleHandlerImpl) handleApply(rule config.ApplyResourceRule, args templating.Args, channel chan<- error) {
	defer func() {
		if err := recover(); err != nil { //catch
			channel <- fmt.Errorf("Recovered from panic when executing apply resource rule: %v", err)
		}
	}()

	logger.Alert("Apply resource rule is not implemented")

	channel <- nil
}

// handleDelete executes Delete Resource rules
func (rh RuleHandlerImpl) handleDelete(rule config.DeleteResourceRule, args templating.Args, channel chan<- error) {
	defer func() {
		if err := recover(); err != nil { //catch
			channel <- fmt.Errorf("Recovered from panic when executing delete resource rule: %v", err)
		}
	}()

	logger.Debugf("Processing delete resource rule:\n%s", rule.Describe())

	templatedSelector, err := rule.Selector.ToTemplatedKubernetesLabelSelector(args)
	if err != nil {
		channel <- fmt.Errorf("Error templating delete resource rule:\n%s\nError: %s", rule.Describe(), err.Error())
		return
	}

	err = rh.client.Sync().Delete(rule.APIVersion, rule.Kind, rule.Namespace, templatedSelector)
	if err != nil {
		channel <- fmt.Errorf("Error applying delete resource rule:\n%s\nError: %s", rule.Describe(), err.Error())
		return
	}

	channel <- nil
}
