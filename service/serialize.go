package service

import (
	"strconv"

	"github.com/mesg-foundation/engine/hash/hashserializer"
)

// HashSerialize returns the hashserialized string of this type
func (data *Service) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("1", data.Name)
	ser.AddString("2", data.Description)
	ser.Add("5", serviceTasks(data.Tasks))
	ser.Add("6", serviceEvents(data.Events))
	ser.Add("7", serviceDependencies(data.Dependencies))
	ser.Add("8", data.Configuration)
	ser.AddString("9", data.Repository)
	ser.AddString("12", data.Sid)
	ser.AddString("13", data.Source)
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data Service_Configuration) HashSerialize() string {
	ser := hashserializer.New()
	ser.AddStringSlice("1", data.Volumes)
	ser.AddStringSlice("2", data.VolumesFrom)
	ser.AddStringSlice("3", data.Ports)
	ser.AddStringSlice("4", data.Args)
	ser.AddString("5", data.Command)
	ser.AddStringSlice("6", data.Env)
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Service_Task) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("1", data.Name)
	ser.AddString("2", data.Description)
	ser.Add("6", serviceParameters(data.Inputs))
	ser.Add("7", serviceParameters(data.Outputs))
	ser.AddString("8", data.Key)
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Service_Parameter) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("1", data.Name)
	ser.AddString("2", data.Description)
	ser.AddString("3", data.Type)
	ser.AddBool("4", data.Optional)
	ser.AddString("8", data.Key)
	ser.AddBool("9", data.Repeated)
	ser.Add("10", serviceParameters(data.Object))
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Service_Event) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("1", data.Name)
	ser.AddString("2", data.Description)
	ser.Add("3", serviceParameters(data.Data))
	ser.AddString("4", data.Key)
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Service_Dependency) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("1", data.Image)
	ser.AddStringSlice("2", data.Volumes)
	ser.AddStringSlice("3", data.VolumesFrom)
	ser.AddStringSlice("4", data.Ports)
	ser.AddString("5", data.Command)
	ser.AddStringSlice("6", data.Args)
	ser.AddString("8", data.Key)
	ser.AddStringSlice("9", data.Env)
	return ser.HashSerialize()
}

type serviceTasks []*Service_Task

// HashSerialize returns the hashserialized string of this type
func (data serviceTasks) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}

type serviceParameters []*Service_Parameter

// HashSerialize returns the hashserialized string of this type
func (data serviceParameters) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}

type serviceEvents []*Service_Event

// HashSerialize returns the hashserialized string of this type
func (data serviceEvents) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}

type serviceDependencies []*Service_Dependency

// HashSerialize returns the hashserialized string of this type
func (data serviceDependencies) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}
