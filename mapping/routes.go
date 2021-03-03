package mapping

import (
	"encoding/json"
	"errors"
	"fmt"
	protoStorage "github.com/kulycloud/protocol/storage"
	"strings"
)

var ErrInvalidReference = errors.New("invalid reference")
var ErrInvalidName = errors.New("invalid name")

func MapIncomingRoute(incoming *IncomingRoute) (*protoStorage.Route, error) {
	route := &protoStorage.Route{
		Host:  incoming.Host,
		Steps: make([]*protoStorage.RouteStep, 0, len(incoming.Steps)),
	}

	ids := make(map[string]uint32, len(incoming.Steps))

	// add first step
	firstStep, ok := incoming.Steps[incoming.Start]
	if !ok {
		return nil, fmt.Errorf("cannot find initial step %s: %w", firstStep, ErrInvalidReference)
	}

	step, err := MapIncomingStep(incoming.Start, &firstStep)
	if err != nil {
		return nil, err
	}

	route.Steps = append(route.Steps, step)
	ids[incoming.Start] = 0

	var nextId uint32 = 1

	// First pass: Convert steps and assign IDs
	for name, incomingStep := range incoming.Steps {
		if name == incoming.Start {
			continue
		}

		step, err := MapIncomingStep(name, &incomingStep)
		if err != nil {
			return nil, err
		}

		route.Steps = append(route.Steps, step)
		ids[name] = nextId
		nextId++
	}

	// Second pass: Resolve references
	for name, incomingStep := range incoming.Steps {
		step := route.Steps[ids[name]]
		for refName, target := range incomingStep.References {
			step.References[refName], ok = ids[target]
			if !ok {
				return nil, fmt.Errorf("cannot find referenced step %s (referenced from step %s): %w", firstStep, name, ErrInvalidReference)
			}
		}
	}

	return route, nil
}

func MapIncomingStep(name string, incomingStep *IncomingStep) (*protoStorage.RouteStep, error) {
	service, err := MapToNamespacedName(incomingStep.Service)
	if err != nil {
		return nil, err
	}

	configJson := []byte("")
	if incomingStep.Config != nil {
		configJson, err = json.Marshal(incomingStep.Config)
		if err != nil {
			return nil, err
		}
	}

	return &protoStorage.RouteStep{
		Service:    service,
		Config:     string(configJson),
		Name:       name,
		References: make(map[string]uint32, len(incomingStep.References)),
	}, nil
}

func MapToNamespacedName(name string) (*protoStorage.NamespacedName, error) {
	parts := strings.SplitN(name, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("cannot parse name %s: %w", name, ErrInvalidName)
	}

	return &protoStorage.NamespacedName{
		Namespace: parts[0],
		Name:      parts[1],
	}, nil
}

func MapFromNamespacedName(namespacedName *protoStorage.NamespacedName) string {
	return fmt.Sprintf("%s:%s", namespacedName.Namespace, namespacedName.Name)
}


func MapRoute(route *protoStorage.Route) (*IncomingRoute, error) {
	inc := &IncomingRoute{
		Host:  route.Host,
		Start: route.Steps[0].Name,
		Steps: make(map[string]IncomingStep, len(route.Steps)),
	}

	for _, step := range route.Steps {
		conf := make(map[string]interface{})

		if step.Config != "" {
			err := json.Unmarshal([]byte(step.Config), &conf)
			if err != nil {
				return nil, err
			}
		}

		references := make(map[string]string, len(step.References))
		for name, ptr := range step.References {
			references[name] = route.Steps[ptr].Name
		}

		inc.Steps[step.Name] = IncomingStep{
			Service:    MapFromNamespacedName(step.Service),
			Config:     conf,
			References: references,
		}
	}

	return inc, nil
}

