/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "strconv"

    "github.com/jessevdk/go-flags"
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

var opts struct {
    OutputFile flags.Filename `short:"o" long:"output-file" description:"The path to write the output to. If flag is not passed, the output will be printed to stdout."`

    Args struct {
        ConfigFiles []flags.Filename `positional-arg-name:"<kubeconfig files>" required:"2" description:"A space seperated list of the kubeconfigs to merge. The first file will be the base. The remaining files will have their unique values changed if there is a duplicate"` 
    }   `positional-args:"yes" required:"yes"`
}

func parseYaml(filename string) Config {
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

func main() {
    parser := flags.NewParser(&opts, flags.Default)

    if _, err := parser.Parse(); err != nil {
        if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrRequired {
            parser.WriteHelp(os.Stdout)
            os.Exit(0)
        } else {
            os.Exit(1)
        }
    }

    base := parseYaml(string(opts.Args.ConfigFiles[0]))

    cluster_nm := 0
    context_nm := 0
    user_nm := 0

    for _, file := range opts.Args.ConfigFiles[1:] {
        tmp := parseYaml(string(file))
        namesChanged := make(map[string]string)

        for _, cluster := range tmp.Clusters {
            for _, i := range base.Clusters {
                if i.Name == cluster.Name {
                    fmt.Printf("Found a duplicate cluster name. Renaming %s from file %s.\n", i.Name, file)
                    cluster_nm ++
                    var updated_name string = cluster.Name+strconv.Itoa(cluster_nm)
                    namesChanged[cluster.Name] = updated_name
                    cluster.Name = updated_name
                    context_nm ++
                    user_nm ++
                }
            }
            base.Clusters = append(base.Clusters, cluster)
        }
        for _, user := range tmp.Users {
            for _, i := range base.Users {
                if i.Name == user.Name {
                    fmt.Printf("Found a duplicate user name, Renaming %s from file %s.\n", i.Name, file)
                    var updated_user_name string = user.Name+strconv.Itoa(user_nm)
                    namesChanged[user.Name] = updated_user_name
                    user.Name = updated_user_name
                }
            }
            base.Users = append(base.Users, user)
        }
        for _, context := range tmp.Contexts {
            for _, i := range base.Contexts {
                if i.Name == context.Name {
                    fmt.Printf("Found a duplicate context name, Renaming %s from file %s.\n", i.Name, file)
                    var updated_context_name string = context.Name+strconv.Itoa(context_nm)
                    context.Name = updated_context_name
                }
            }

            // Ensure user and cluster are updated in the context if they were changed
            if val, ok := namesChanged[context.Context.Cluster]; ok {
                context.Context.Cluster = val
            }
            if val, ok := namesChanged[context.Context.User]; ok {
                context.Context.User = val
            }
            base.Contexts = append(base.Contexts, context)
        }
    }

    out, _ := yaml.Marshal(&base)

    if len(opts.OutputFile) > 0 {
        fmt.Printf("Writing the output to %s\n", string(opts.OutputFile))
        err := ioutil.WriteFile(string(opts.OutputFile), out, 0644)

        if err != nil {
            panic(err)
        }
    } else {
        fmt.Println(string(out))
    }
}
