package models

import "time"

type Metadata struct {
	Environment  string `json:"environment" bson:"environment" binding:"required,oneof=dev staging prod"`
	Region       string `json:"region" bson:"region" binding:"required"`
	Version      int    `json:"version" bson:"version" binding:"required"`
	Developer    string `json:"developer" bson:"developer"`       // optional
	Experimental bool   `json:"experimental" bson:"experimental"` // optional
}

type Instance struct {
	ServiceName   string    `json:"serviceName" bson:"serviceName" binding:"required"`
	ID            string    `json:"id" bson:"id" binding:"required"`
	Host          string    `json:"host" bson:"host" binding:"required"`
	Port          int       `json:"port" bson:"port" binding:"required"`
	Mode          string    `json:"mode" bson:"mode" binding:"required,oneof=dev staging prod"`
	Metadata      Metadata  `json:"metadata" bson:"metadata" binding:"required"`
	Health        string    `json:"health" bson:"health"`
	LastHeartbeat time.Time `json:"lastHeartbeat" bson:"lastHeartbeat"`
}
