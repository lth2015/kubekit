package appmachinery

import (
	"regexp"
	"strings"
	"fmt"
)

/*
Workload example:
	default/pcbt-starter-demo001-1544492700-cmfgq         CronJob      pcbt-starter-demo001
	ingress-nginx/default-http-backend-68ff788cd7-lz2rm   Deployment   default-http-backend
	litest2/final-mysql-slave-0                           StatefulSet  final-mysql-slave
	weave/weave-scope-agent-z6mxl   					  DaemonSet    weave-scope-agent
*/

type WorkloadRegex struct {
	hashRegex *regexp.Regexp
	randomRegex *regexp.Regexp
	sequnceRegex *regexp.Regexp
}

func NewWorkloadRegex() *WorkloadRegex {
	wlr := &WorkloadRegex{}
	wlr.hashRegex, _ = regexp.Compile(Hash)
	wlr.randomRegex, _ = regexp.Compile(Random)
	wlr.sequnceRegex, _ = regexp.Compile(Sequence)
	return wlr
}

func (this *WorkloadRegex) GetAppName(pod string) (string, error) {
	ss := strings.Split(pod, "-")
	if len(ss) >= 2 && this.randomRegex.MatchString(ss[len(ss)-1]) {
		if this.hashRegex.MatchString(ss[len(ss)-2]) {
			// Deployment or CronJob
			index := strings.LastIndex(pod, "-")
			app := pod[:index]
			index = strings.LastIndex(app, "-")
			app = pod[:index]
			return app, nil
		} else {
			// DaemonSet
			index := strings.LastIndex(pod, "-")
			app := pod[:index]
			return app, nil
		}
	}

	if len(ss) >= 1 && !this.hashRegex.MatchString(ss[len(ss)-2]) && this.sequnceRegex.MatchString(ss[len(ss)-1]) {
		// StatefulSet
		index := strings.LastIndex(pod, "-")
		app := pod[:index]
		return app, nil
	}

	return "", fmt.Errorf("Unknown controller of this pod: pod=%s", pod)
}