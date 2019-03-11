// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/33cn/chain33/cmd/tools/strategy"
	"github.com/33cn/chain33/cmd/tools/types"
	"github.com/spf13/cobra"
)

//GenDappCmd advance cmd
func GenDappCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gendapp",
		Short: "auto generate chain33 dapp base code",
		Run:   genDapp,
	}
	addGenDappFlag(cmd)
	return cmd
}

func addGenDappFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("name", "n", "", "dapp name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringP("output", "o", "", "output directory for new dapp, default is $GOPATH/src/github.com/33cn/plugin/plugin/dapp/")
	cmd.Flags().StringP("proto", "p", "", "dapp protobuf file")
}

func genDapp(cmd *cobra.Command, args []string) {

	dappName, _ := cmd.Flags().GetString("name")
	outDir, _ := cmd.Flags().GetString("output")
	propFile, _ := cmd.Flags().GetString("proto")

	s := strategy.New(types.KeyGenDapp)
	if s == nil {
		fmt.Println(types.KeyGenDapp, "Not support")
		return
	}

	s.SetParam(types.KeyExecutorName, dappName)
	s.SetParam(types.KeyDappOutDir, outDir)
	s.SetParam(types.KeyProtobufFile, propFile)
	s.Run()
}
