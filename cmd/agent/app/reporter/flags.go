// Copyright (c) 2018 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package reporter

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/cmd/all-in-one/setupcontext"
	"github.com/jaegertracing/jaeger/cmd/flags"
)

const (
	// Whether to use grpc, tchannel or http reporter.
	reporterType = "reporter.type"
	// AgentTagsDeprecated is a configuration property name for adding process tags to incoming spans.
	AgentTagsDeprecated = "jaeger.tags"
	agentTags           = "agent.tags"
	// GRPC is name of gRPC reporter.
	GRPC Type = "grpc"
	// HTTP is name of http reporter.
	HTTP Type = "http"
)

// Type defines type of reporter.
type Type string

// Options holds generic reporter configuration.
type Options struct {
	ReporterType Type
	AgentTags    map[string]string
}

// AddFlags adds flags for Options.
func AddFlags(flags *flag.FlagSet) {
	flags.String(reporterType, string(GRPC), fmt.Sprintf("Reporter type to use: %s, %s", string(GRPC), string(HTTP)))
	if !setupcontext.IsAllInOne() {
		flags.String(AgentTagsDeprecated, "", "(deprecated) see --"+agentTags)
		flags.String(agentTags, "", "One or more tags to be added to the Process tags of all spans passing through this agent. Ex: key1=value1,key2=${envVar:defaultValue}")
	}
}

// InitFromViper initializes Options with properties retrieved from Viper.
func (b *Options) InitFromViper(v *viper.Viper, logger *zap.Logger) *Options {
	b.ReporterType = Type(v.GetString(reporterType))
	if !setupcontext.IsAllInOne() {
		if len(v.GetString(AgentTagsDeprecated)) > 0 {
			logger.Warn("Using deprecated configuration", zap.String("option", AgentTagsDeprecated))
			b.AgentTags = flags.ParseJaegerTags(v.GetString(AgentTagsDeprecated))
		}
		if len(v.GetString(agentTags)) > 0 {
			b.AgentTags = flags.ParseJaegerTags(v.GetString(agentTags))
		}
	}
	return b
}
