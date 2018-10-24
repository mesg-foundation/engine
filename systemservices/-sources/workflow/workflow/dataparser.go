package workflow

import (
	"fmt"
	"strings"
)

// dataParser parses workflow variables ($data, $configs, $services)
// and returns it's value.
// if giving data is not a variable it's directly returned.
type dataParser struct {
	configs  []ConfigDefinition
	services []ServiceDefinition
	data     map[string]interface{}
}

// Parse converts a variable or a value to value.
func (c *dataParser) Parse(value interface{}) (interface{}, error) {
	valueStr, ok := value.(string)
	if !ok {
		return value, nil
	}
	// check if it's a variable.
	object, isVar, err := c.parseVar(valueStr)
	if err != nil {
		return nil, err
	}
	// otherwise return value directly.
	if !isVar {
		return value, nil
	}
	return c.convertVar(valueStr, object)
}

func (c *dataParser) parseVar(value string) (object []string, isVar bool, err error) {
	if !strings.HasPrefix(value, "$") {
		return nil, false, nil
	}
	s := strings.Split(value, "$")
	ss := strings.Split(s[1], ".")
	if len(ss) == 0 {
		return nil, false, &invalidVarErr{value}
	}
	return ss, true, nil
}

func (c *dataParser) convertVar(str string, object []string) (interface{}, error) {
	switch object[0] {
	case "services":
		return c.convertService(str, object[1:])
	case "configs":
		return c.convertConfig(str, object[1:])
	case "data":
		return c.convertData(str, object[1:])
	default:
		return nil, &invalidVarErr{str}
	}
}

func (c *dataParser) convertService(str string, object []string) (interface{}, error) {
	for _, c := range c.services {
		if c.Name == object[0] {
			return c.ID, nil
		}
	}
	return nil, &invalidVarErr{str}
}

func (c *dataParser) convertConfig(str string, object []string) (interface{}, error) {
	for _, cf := range c.configs {
		// match with the related config.
		if cf.Key == object[0] {
			value := cf.Value
			for i, s := range object {
				if i == 0 {
					continue
				}
				m, ok := value.(map[string]interface{})
				if ok {
					value = m[s]
				} else {
					// yml supports any type as map keys.
					// for dynamically parsed maps, try asserting
					// map keys as an interface{}.
					m, ok := value.(map[interface{}]interface{})
					if ok {
						var ss interface{}
						ss = s
						value = m[ss]
					} else {
						return nil, &invalidVarErr{str}
					}
				}
			}
			return value, nil
		}
	}
	return nil, &invalidVarErr{str}
}

func (c *dataParser) convertData(str string, object []string) (interface{}, error) {
	for key, d := range c.data {
		// match with the related data.
		if key == object[0] {
			value := d
			for i, ss := range object {
				if i == 0 {
					continue
				}
				m, ok := value.(map[string]interface{})
				if ok {
					value = m[ss]
				} else {
					return nil, &invalidVarErr{str}
				}
			}
			return value, nil
		}
	}
	return nil, &invalidVarErr{str}
}

type invalidVarErr struct {
	variable string
}

func (e *invalidVarErr) Error() string {
	return fmt.Sprintf("invalid variable %q", e.variable)
}
