package server

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
)

// universal config keys:
// max_cpu
func (this *Server) LoadConfig(fn string) *Server {
	log.Info("Server[%s.%s] loading config file %s", this.Name, BuildID, fn)
	this.configFile = fn

	var err error
	this.Conf, err = conf.Load(fn)
	if err != nil {
		panic(err)
	}

	return this
}
