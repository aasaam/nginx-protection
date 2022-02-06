package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	prometheusRequestTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "protection_request_total",
		Help: "The total number of request",
	})

	prometheusRequestAuthFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "protection_request_auth_failed",
		Help: "The total number of auth failed request",
	})

	prometheusRequestError = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "protection_response_error",
		Help: "The total number of response error",
	}, []string{"status"})

	prometheusRequestAuthSuccess = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "protection_request_auth_success",
		Help: "The total number of auth success request",
	}, []string{"acl", "value"})

	prometheusRequestChallenge = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "protection_request_challenge",
		Help: "The total number of challenge request",
	}, []string{"challenge_type"})

	prometheusRequestChallengeSuccess = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "protection_request_challenge_success",
		Help: "The total number solved challenge",
	}, []string{"challenge_type"})

	prometheusRequestChallengeFailed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "protection_request_challenge_failed",
		Help: "The total number failed challenge",
	}, []string{"challenge_type"})
)

func getPrometheusRegistry() *prometheus.Registry {
	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(prometheusRequestTotal)
	promRegistry.MustRegister(prometheusRequestError)
	promRegistry.MustRegister(prometheusRequestAuthFailed)
	promRegistry.MustRegister(prometheusRequestAuthSuccess)
	promRegistry.MustRegister(prometheusRequestChallenge)
	promRegistry.MustRegister(prometheusRequestChallengeSuccess)
	promRegistry.MustRegister(prometheusRequestChallengeFailed)
	return promRegistry
}
