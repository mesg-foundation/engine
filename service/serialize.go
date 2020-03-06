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
	return hashserializer.New().
		AddString("1", data.Name).
		AddString("2", data.Description).
		Add("5", serviceTasks(data.Tasks)).
		Add("6", serviceEvents(data.Events)).
		Add("7", serviceDependencies(data.Dependencies)).
		Add("8", data.Configuration).
		AddString("9", data.Repository).
		AddString("12", data.Sid).
		AddString("13", data.Source).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data Service_Configuration) HashSerialize() string {
	return hashserializer.New().
		AddStringSlice("1", data.Volumes).
		AddStringSlice("2", data.VolumesFrom).
		AddStringSlice("3", data.Ports).
		AddStringSlice("4", data.Args).
		AddString("5", data.Command).
		AddStringSlice("6", data.Env).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Service_Task) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("1", data.Name).
		AddString("2", data.Description).
		Add("6", serviceParameters(data.Inputs)).
		Add("7", serviceParameters(data.Outputs)).
		AddString("8", data.Key).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Service_Parameter) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("1", data.Name).
		AddString("2", data.Description).
		AddString("3", data.Type).
		AddBool("4", data.Optional).
		AddString("8", data.Key).
		AddBool("9", data.Repeated).
		Add("10", serviceParameters(data.Object)).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Service_Event) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("1", data.Name).
		AddString("2", data.Description).
		Add("3", serviceParameters(data.Data)).
		AddString("4", data.Key).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Service_Dependency) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("1", data.Image).
		AddStringSlice("2", data.Volumes).
		AddStringSlice("3", data.VolumesFrom).
		AddStringSlice("4", data.Ports).
		AddString("5", data.Command).
		AddStringSlice("6", data.Args).
		AddString("8", data.Key).
		AddStringSlice("9", data.Env).
		HashSerialize()
}

type serviceTasks []*Service_Task

// HashSerialize returns the hashserialized string of this type
func (data serviceTasks) HashSerialize() string {
	if data == nil || len(data) == 0 {
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
	if data == nil || len(data) == 0 {
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
	if data == nil || len(data) == 0 {
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
	if data == nil || len(data) == 0 {
		return ""
	}
	ser := hashserializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}
