//
// Copyright (C) 2020 OpenSIPS Solutions
//
// Call API is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Call API is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.
//

package main

import (
	"flag"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/OpenSIPS/call-api/pkg/cmd"
	"github.com/OpenSIPS/call-api/pkg/proxy"
	"github.com/OpenSIPS/call-api/pkg/config"
)

func usage(prog string) {
	logrus.Fatalf("Usage: %s command [arguments...]", prog)
}

func main() {

	cfgPath, err := config.ParseFlags("call-cmd")
	if err != nil {
		logrus.Fatal(err)
	}

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		logrus.Fatal(err)
	}

	/* initialize logging */
	logfile, err := config.InitLogging(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	if logfile != nil {
		defer logfile.Close()
	}

	if flag.NArg() < 1 {
		logrus.Error("no command specified!")
		usage(os.Args[0])
	}

	proxy := proxy.NewProxy(cfg)
	if proxy == nil {
		logrus.Fatal("could not initialize SIP proxy")
	}
	command := flag.Arg(0)
	logrus.Debugf("Running command %s", command)
	c := cmd.New(command, "", proxy)
	if c == nil {
		logrus.Fatalf("could not initialize %s command", command)
	}
	var arguments = map[string]interface{}{}
	for _, arg := range flag.Args()[1:] {
		param := strings.Split(arg, "=")
		arguments[param[0]] = strings.Join(param[1:], "=")
	}
	c.Run(arguments)
	for {
		event := <-c.Wait()
		if event == nil {
			break
		} else if event.IsError() {
			logrus.Fatal(event.String())
		} else {
			logrus.Printf("%s:%s: %+v", c.Command, c.ID, event.String())
		}
	}
}
