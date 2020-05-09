module github.com/ontio/ontology-tool

go 1.12

require (
	github.com/FactomProject/basen v0.0.0-20150613233007-fe3947df716e // indirect
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/blang/semver v3.5.1+incompatible
	github.com/hashicorp/golang-lru v0.5.3
	github.com/ontio/ontology v1.9.0
	github.com/ontio/ontology-crypto v1.0.8
	github.com/ontio/ontology-eventbus v0.9.1
	github.com/ontio/ontology-go-sdk v1.11.1
	github.com/scylladb/go-set v1.0.2
)

replace github.com/ontio/ontology v1.9.0 => ../ontology
