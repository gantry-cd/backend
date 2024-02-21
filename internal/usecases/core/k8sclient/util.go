package k8sclient

import v1 "k8s.io/api/core/v1"

// toEnvVar key, valueからは環境変数を生成する .
func toEnvVar(k, v string) v1.EnvVar {
	return v1.EnvVar{
		Name:  k,
		Value: v,
	}
}
