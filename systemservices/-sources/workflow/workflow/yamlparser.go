package workflow

import (
	"io"
	"sort"

	"gopkg.in/yaml.v2"
)

// YAMLDefinition represents the definition of Workflow in YML.
type YAMLDefinition struct {
	// Name of workflow.
	// e.g.:
	//   name: discord-invites
	Name string `yaml:"name"`

	// Description of workflow.
	// e.g.:
	//   description: a long description
	Description string `yaml:"description"`

	// Services describes service name, service id pair.
	// e.g.:
	//   discord: c6651756923f8569ef17bf31b0baef78
	Services map[string]string `yaml:"services"`

	// Configs is for setting constants for workflow.
	// constants values can be a nested object as well as basic data types.
	// config name, value pair.
	// to access configs use $ identifier.
	// e.g.:
	//   key: STBGu4DY3AMH1aw
	//   defaultUser:
	//     company: mesg
	//     languages:
	//       - tr
	//       - en
	//       - go
	Configs map[string]interface{} `yaml:"configs"`

	// When keeps the actual workflow.
	// service name, when pair.
	When map[string]When `yaml:"when"`
}

// When keeps the actual workflow.
// e.g.:
//   when:
//     webhook:
//       event:
//         request:
//           map:
//             email: $data.data.email
//             sendgridAPIKey: $configs.key
//           execute:
//             discord: send
type When struct {
	// Event keeps the event-execute flow.
	// event name, event listener pair.
	Event map[string]EventListener `yaml:"event"`
}

// EventListener keeps the event-execute flow.
type EventListener struct {
	// Map is to map event data to task inputs.
	// event data, value pair.
	Map map[string]interface{} `yaml:"map"`

	// Execute keeps the service key, task key information to execute
	// a task after event received.
	// service name, task name pair.
	Execute map[string]string `yaml:"execute"`
}

// ParseYAML parses a workflow yml file from an io.Reader to Definition.
func ParseYAML(r io.Reader) (WorkflowDefinition, error) {
	return newYAMLParser(r).parse()
}

// yamlParser parses workflow yml.
type yamlParser struct {
	r io.Reader
}

// newYAMLParser returns a yml parser.
func newYAMLParser(r io.Reader) *yamlParser {
	return &yamlParser{r}
}

func (p *yamlParser) parse() (WorkflowDefinition, error) {
	def := WorkflowDefinition{}

	ydef, err := p.parseYAML()
	if err != nil {
		return def, err
	}

	def.Name = ydef.Name
	def.Description = ydef.Description
	def.Services = p.toServices(ydef.Services)
	def.Configs = p.toConfigs(ydef.Configs)
	def.Events = p.toEvents(ydef.When)

	return def, nil
}

// parseYAML parses yaml file into yaml definition.
func (p *yamlParser) parseYAML() (YAMLDefinition, error) {
	ydef := YAMLDefinition{}
	if err := yaml.NewDecoder(p.r).Decode(&ydef); err != nil {
		return ydef, err
	}
	return ydef, nil
}

// toServices converts service name, id pairs to a Service slice.
// sorting is made be consistent in tests.
func (p *yamlParser) toServices(s map[string]string) []ServiceDefinition {
	sorted := make([]string, 0)
	for key := range s {
		sorted = append(sorted, key)
	}
	sort.Strings(sorted)

	var services []ServiceDefinition
	for _, name := range sorted {
		services = append(services, ServiceDefinition{
			Name: name,
			ID:   s[name],
		})
	}
	return services
}

// toConfigs converts config key, value pairs to a Config slice.
// sorting is made be consistent in tests.
func (p *yamlParser) toConfigs(s map[string]interface{}) []ConfigDefinition {
	sorted := make([]string, 0)
	for key := range s {
		sorted = append(sorted, key)
	}
	sort.Strings(sorted)

	var configs []ConfigDefinition
	for _, key := range sorted {
		configs = append(configs, ConfigDefinition{
			Key:   key,
			Value: s[key],
		})
	}
	return configs
}

// toEvents converts service name, when pairs to an Event slice.
// sorting is made be consistent in tests.
func (p *yamlParser) toEvents(s map[string]When) []EventDefinition {
	sorted := make([]string, 0)
	for key := range s {
		sorted = append(sorted, key)
	}
	sort.Strings(sorted)

	var events []EventDefinition
	for _, serviceName := range sorted {
		event := EventDefinition{
			ServiceName: serviceName,
		}

		for eventKey, eventListener := range s[serviceName].Event {
			event.EventKey = eventKey
			event.Map = p.toMaps(eventListener.Map)
			event.Execute = p.toExecute(eventListener.Execute)
		}

		events = append(events, event)
	}
	return events
}

// toMaps converts input data key, value pairs to a Map slice.
// sorting is made be consistent in tests.
func (p *yamlParser) toMaps(s map[string]interface{}) []MapDefinition {
	sorted := make([]string, 0)
	for key := range s {
		sorted = append(sorted, key)
	}
	sort.Strings(sorted)

	var maps []MapDefinition
	for _, key := range sorted {
		maps = append(maps, MapDefinition{
			Key:   key,
			Value: s[key],
		})
	}
	return maps
}

// toExecute converts service name, task key pairs to an Execute.
func (p *yamlParser) toExecute(s map[string]string) ExecuteDefinition {
	var execute ExecuteDefinition
	for serviceName, taskKey := range s {
		execute = ExecuteDefinition{
			ServiceName: serviceName,
			TaskKey:     taskKey,
		}
	}
	return execute
}
