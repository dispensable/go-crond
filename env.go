package main

import (
	"fmt"
	"strings"
)

var (
	preloadEnvVars map[string]struct{} = map[string]struct{}{
		"SHELL": struct{}{},
		"PWD": struct{}{},
		"LOGNAME": struct{}{},
		"HOME": struct{}{},
		"LANG": struct{}{},
		"USER": struct{}{},
		"SHLVL": struct{}{},
		"MAILTO": struct{}{},
		"LC_ALL": struct{}{},
		"PATH": struct{}{},
		"_": struct{}{},
	}
	reservedEnvVars []string = []string{"LOGNAME", "USER"}
)

type CronEnvs struct {
	envMap map[string]string
}

func NewCronEnvs(parentEnvs []string) *CronEnvs {
	pEnvs := make(map[string]string)
	for _, e := range parentEnvs {
		splitAt := strings.Index(e, "=")
		key := e[:splitAt]
		if _, ok := preloadEnvVars[key]; ok {
			pEnvs[key] = e[splitAt+1:]
		}
	}

	return &CronEnvs{
		envMap: pEnvs,
	}
}

func (e *CronEnvs) SetEnv(key, value string) {
	e.envMap[key] = value
}

func (e *CronEnvs) GetEnv(key string) string {
	if value, ok := e.envMap[key]; ok {
		return value
	} else {
		return ""
	}
}

func (e *CronEnvs) UpdateEnvForJob(c *CrontabEntry) error {
	for _, envStr := range c.Env {
		splitAt := strings.Index(envStr, "=")
		key := envStr[:splitAt]
		for _, env := range reservedEnvVars {
			if env == key {
				return fmt.Errorf("User can't set env %s", key)
			}
		}
		e.envMap[key] = envStr[splitAt+1:]	
	}
	return nil
}

func (e *CronEnvs) GetEnvList() []string {
	result := []string{}
	for k, v := range e.envMap {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}

