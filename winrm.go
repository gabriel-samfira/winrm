/*
Copyright 2013 Brice Figureau

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/masterzen/winrm/winrm"
)

func main() {
	var (
		hostname string
		user     string
		pass     string
		cmd      string
		port     int
		https    bool
		insecure bool
		cacert   string
		cert     string
		key      string
		auth     string
	)

	flag.StringVar(&hostname, "hostname", "localhost", "winrm host")
	flag.StringVar(&user, "username", "vagrant", "winrm admin username")
	flag.StringVar(&pass, "password", "vagrant", "winrm admin password")
	flag.IntVar(&port, "port", 5985, "winrm port")
	flag.BoolVar(&https, "https", false, "use https")
	flag.BoolVar(&insecure, "insecure", false, "skip SSL validation")
	flag.StringVar(&cacert, "cacert", "", "CA certificate to use")
	flag.StringVar(&cert, "cert", "", "X509 certificate to use")
	flag.StringVar(&key, "key", "", "X509 key to use")
	flag.StringVar(&auth, "auth", "basic", "Authentication method to use: cert/basic")
	flag.Parse()

	var certBytes []byte
	var err error
	var client *winrm.Client
	var certBody []byte
	var keyBody []byte
	authentication := winrm.Auth(auth)
	if cacert != "" {
		certBytes, err = ioutil.ReadFile(cacert)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		certBytes = nil
	}

	cmd = flag.Arg(0)
	if winrm.Auth(auth) == winrm.Cert {
		if cert == "" || key == "" {
			fmt.Println("Auth cert requires -cert and -key")
			os.Exit(1)
		}
		certBody, err = ioutil.ReadFile(cert)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		keyBody, err = ioutil.ReadFile(key)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	client, err = winrm.NewClient(&winrm.Endpoint{Host: hostname, Port: port, HTTPS: https, Insecure: insecure, Cert: certBody, Key: keyBody, CACert: &certBytes}, user, pass, authentication)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exitCode, err := client.RunWithInput(cmd, os.Stdout, os.Stderr, os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}
