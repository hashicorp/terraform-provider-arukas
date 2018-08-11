package arukas

import (
	"net"

	"github.com/yamamoto-febc/go-arukas"
)

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, string(v.(string)))
	}
	return vs
}

func expandEnvs(configured interface{}) []*arukas.Env {
	var envs []*arukas.Env
	if configured == nil {
		return envs
	}
	rawEnvs := configured.([]interface{})
	for _, raw := range rawEnvs {
		env := raw.(map[string]interface{})
		key := env["key"].(string)
		value := env["value"].(string)
		envs = append(envs, &arukas.Env{Key: key, Value: value})
	}
	return envs
}

func flattenEnvs(envs []*arukas.Env) []interface{} {
	var ret []interface{}
	for _, env := range envs {
		r := map[string]interface{}{}
		r["key"] = env.Key
		r["value"] = env.Value
		ret = append(ret, r)
	}
	return ret
}

func expandPorts(configured interface{}) []*arukas.Port {
	var ports []*arukas.Port
	if configured == nil {
		return ports
	}
	rawPorts := configured.([]interface{})
	for _, raw := range rawPorts {
		port := raw.(map[string]interface{})
		proto := port["protocol"].(string)
		num := port["number"].(int)
		ports = append(ports, &arukas.Port{Protocol: proto, Number: int32(num)})
	}
	return ports
}

func flattenPorts(ports []*arukas.Port) []interface{} {
	var ret []interface{}
	for _, port := range ports {
		r := map[string]interface{}{}
		r["protocol"] = port.Protocol
		r["number"] = port.Number
		ret = append(ret, r)
	}
	return ret
}

func flattenPortMappings(ports []*arukas.PortMapping) []interface{} {
	var ret []interface{}
	for _, port := range ports {
		r := map[string]interface{}{}
		ip := ""

		addrs, err := net.LookupHost(port.Host)
		if err == nil && len(addrs) > 0 {
			ip = addrs[0]
		}

		r["host"] = port.Host
		r["ipaddress"] = ip
		r["container_port"] = port.ContainerPort
		r["service_port"] = port.ServicePort
		ret = append(ret, r)
	}
	return ret
}
