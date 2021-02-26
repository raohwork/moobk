// This file is part of moobk.
//
// moobk is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 2 of the License, or
// (at your option) any later version.
//
// moobk is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with moobk.  If not, see <https://www.gnu.org/licenses/>.
package cmd

import (
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "moobk",
	Short: "moobk is a simple tool helps you to backup CoW filesystem",
	Long: `moobk, stands for "MOOOOOOOOOOO Backup", is a simple tool to backup CoW filesystems. Currently only btrfs and zfs are supported.

moobk aims to help you automating daily backup work with ease. Say you have a laptop with two btrfs subvolumes mounted at root and /home, and an btrfs-formatted external storage. You may periodically run this to take snapshot (like, hourly):

  moobk snap -t btrfs -r local:///.backup / rootfs
  moobk snap -t btrfs -r local:///.backup /home

Then, write an udev rule to transfer new snapshots to external storage (and delete snapshots older than 2 weeks) when plugged

  moobk sync -t btrfs -r local:///.root.backup local:///path/to/external/storage
  moobk purge -t btrfs -r local:///.root.backup local:///path/to/external/storage 2w

Run "moobk help" for further info.
`,
	SuggestionsMinimumDistance: 3,
	Version:                    "0.0.1",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
