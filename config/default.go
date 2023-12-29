package config

const defaultJson string = `{
	"service": {
		"address": ":9713",
		"ttl": 15,
		"interval": 10
	},
	"logger": {
		"level": "info",
		"file": "logs/server.log",
		"std": false
	},
	"database": {
		"name": "rgsCloud",
		"ip": "192.168.1.10",
		"port": "27017",
		"user": "root",
		"password": "pass2019",
		"type": "mongodb"
	},
	"album":{
		"person":{
			"count":200,
			"size": 2097152
		},
		"group":{
			"count":600,
			"size": 10485760
		}
	},
	"basic": {
		"synonym": 6,
		"tag": 6
	}
}
`
