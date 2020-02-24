/*
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.

*/
package cmd

import (
    "fmt"
    "os"
    "github.com/mjnt/k8stk/util"

    "github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
    Use:   "remove [flags] KUBECONFIG",
    Short: "Remove a context in the specified kubeconfig file",
    Long: `Removes a context from the specified kubeconfig file. Outputs to stdout
unless an output file is specified.`,
    Args: cobra.MinimumNArgs(1),
    Example: "k8stk remove -o newConfig -c kubeadmin ~/.kube/config",
    Run: func(cmd *cobra.Command, args []string) {
        removeContext(cmd, args)
    },
}

func init() {
    rootCmd.AddCommand(removeCmd)
    
    // Here you will define your flags and configuration settings.
    
    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // removeCmd.PersistentFlags().String("foo", "", "A help for foo")
    
    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    removeCmd.Flags().StringP("output", "o", "", "File to output the merged config to. If not passed the output will be sent to stdout")
    removeCmd.Flags().StringP("context", "c", "", "The context name to remove. (requried)")
    removeCmd.MarkFlagRequired("context")
}

func removeContext(cmd *cobra.Command, args []string) {
    base := util.ParseYaml(args[0])

    remove_context := cmd.Flags().Lookup("context").Value.String()

    new_config := util.Config{}
    new_config.ApiVersion = base.ApiVersion
    new_config.Kind = base.Kind
    new_config.Preferences = base.Preferences

    remove_cluster := ""
    remove_user := ""

    for _, i := range base.Contexts {
        if i.Name == remove_context {
            remove_cluster = i.Context.Cluster
            remove_user = i.Context.User
        } else {
            new_config.Contexts = append(new_config.Contexts, i)
        }
    }

    if remove_cluster == "" {
        fmt.Printf("Error: Context %s not found\n", remove_context)
        os.Exit(1)
    }

    for _, i := range base.Clusters {
        if i.Name != remove_cluster {
            new_config.Clusters = append(new_config.Clusters, i)
        }
    }

    for _, i := range base.Users {
        if i.Name != remove_user {
            new_config.Users = append(new_config.Users, i)
        }
    }

    if base.CurrentContext == remove_context {
        new_config.CurrentContext = new_config.Contexts[0].Name
    } else {
        new_config.CurrentContext = base.CurrentContext
    }

    util.OutputYaml(new_config, cmd.Flags().Lookup("output").Value.String())
}
