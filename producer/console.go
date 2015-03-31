// Copyright 2015 trivago GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package producer

import (
	"fmt"
	"github.com/trivago/gollum/core"
	"github.com/trivago/gollum/shared"
	"os"
	"runtime"
	"strings"
	"sync"
)

// Console producer plugin
// Configuration example
//
//   - "producer.Console":
//     Enable: true
//     Console: "stderr"
//
// Console may either be "stdout" or "stderr"
type Console struct {
	core.ProducerBase
	console *os.File
}

func init() {
	shared.RuntimeType.Register(Console{})
}

// Configure initializes this producer with values from a plugin config.
func (prod *Console) Configure(conf core.PluginConfig) error {
	err := prod.ProducerBase.Configure(conf)
	if err != nil {
		return err
	}

	console := conf.GetString("Console", "stdout")

	switch strings.ToLower(console) {
	default:
		fallthrough
	case "stdout":
		prod.console = os.Stdout
	case "stderr":
		prod.console = os.Stderr
	}

	return nil
}

func (prod Console) printMessage(msg core.Message) {
	prod.Formatter().PrepareMessage(msg)
	fmt.Fprint(prod.console, prod.Formatter().String())
}

func (prod Console) flush() {
	for prod.NextNonBlocking(prod.printMessage) {
		runtime.Gosched()
	}
}

// Produce writes to stdout or stderr.
func (prod Console) Produce(workers *sync.WaitGroup) {
	defer func() {
		prod.flush()
		prod.WorkerDone()
	}()

	prod.AddMainWorker(workers)
	prod.DefaultControlLoop(prod.printMessage, nil)
}