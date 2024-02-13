package k8sclient

import (
	corev1 "k8s.io/api/core/v1"
)

type option struct {
	labelSelector   map[string]string
	replica         *int32
	containerOption map[string]containerOption
}

type containerOption struct {
	imagePullPolicy corev1.PullPolicy
}

func newOption() *option {
	return &option{
		labelSelector: map[string]string{
			CreatedByLabel: AppIdentifier,
		},
		containerOption: make(map[string]containerOption),
	}
}

type Option func(*option)

func WithLabelSelector(labelSelector map[string]string) Option {
	return func(o *option) {
		o.labelSelector = labelSelector
	}
}

func WithCreatedByLabel(name string) Option {
	return func(o *option) {
		o.labelSelector[CreatedByLabel] = name
	}

}

func WithRepositoryLabel(repository string) Option {
	return func(o *option) {
		o.labelSelector[RepositryLabel] = repository
	}
}

func WithPrIDLabel(prID string) Option {
	return func(o *option) {
		o.labelSelector[PrIDLabel] = prID
	}
}

func WithBaseBranchLabel(baseBranch string) Option {
	return func(o *option) {
		o.labelSelector[BaseBranchLabel] = baseBranch
	}
}

func WithEnvirionmentLabel(envirionment string) Option {
	return func(o *option) {
		o.labelSelector[EnvirionmentLabel] = envirionment
	}
}

func WithReplica(replica int32) Option {
	return func(o *option) {
		o.replica = &replica
	}
}

func WithImagePullPolicy(containerName, policy string) Option {
	return func(o *option) {
		o.containerOption[containerName] = containerOption{
			imagePullPolicy: corev1.PullPolicy(policy),
		}
	}
}
