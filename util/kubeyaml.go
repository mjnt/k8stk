/*
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.

*/
package util

import (
    "fmt"
    "io/ioutil"

    "gopkg.in/yaml.v3"
)

type Config struct {
    ApiVersion string `yaml:"apiVersion"`
    Clusters []struct {
        Cluster struct {
            CaData string `yaml:"certificate-authority-data"`
            Server string
        }
        Name string
    }
    Contexts []struct {
        Context struct {
            Cluster string
            User string
        }
        Name string
    }
    CurrentContext string `yaml:"current-context"`
    Kind string
    Preferences map[string]string
    Users []struct {
        Name string
        User struct {
            CertData string `yaml:"client-certificate-data"`
            KeyData string `yaml:"client-key-data"`
        }
    }
}

func ParseYaml(filename string) Config {
    source, err := ioutil.ReadFile(filename)

    if err != nil {
        panic(err)
    }

    c := Config{}
    err = yaml.Unmarshal(source, &c)
    if err != nil {
        panic(err)
    }

    return c
}

func OutputYaml(base Config, output_file string) {
    out, _ := yaml.Marshal(&base)

    if output_file != "" {
        fmt.Printf("Writing the output to %s\n", output_file)
        err := ioutil.WriteFile(output_file, out, 0644)

        if err != nil {
            panic(err)
        }
    } else {
        fmt.Println(string(out))
    }
}
