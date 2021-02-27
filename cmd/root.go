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
	"os"

	"github.com/raohwork/moobk/moodrvs"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "moobk",
	Short: "moobk is a simple tool helps you to backup CoW filesystem",
	Long: `moobk, stands for "MOOOOOOOOOOO Backup", is a simple tool to backup CoW filesystems. Currently only btrfs and zfs are supported.

moobk aims to help you automating daily backup work with ease. It's roughly a frontend of zfs/btrfs. Say you have a laptop with two btrfs subvolumes mounted at root and /home, and an btrfs-formatted external storage. You may periodically run this to take snapshot (like, hourly):

  moobk snap -t btrfs -r local:///.backup / rootfs
  moobk snap -t btrfs -r local:///.backup /home

Then, write an udev rule to transfer new snapshots to external storage (and delete snapshots older than 2 weeks) when plugged

  moobk sync -t btrfs -r local:///.root.backup local:///path/to/external/storage
  moobk purge -t btrfs -r local:///.root.backup local:///path/to/external/storage 2w

Those are long lines suitable for scripts which only write once. If you like to run it manually (like, for testing), moobk supports repo alias. Create a file named ".moobk.yaml" in your home directory (or /root if you run moobk with sudo) with following lines

me: local:///path/to/local/backup
nas: ssh+sudo://user@nas.my.home/path/to/backup?ssh_i=/home/user/.ssh/id_ed25519
usb: local:///path/to/usb/backup

Then you can run

  moobk list -r me     # equals to moobk list -r local:///path/to/local/backup
  moobk sync -r me nas

To get some idea about what "repo" is, run "moobk repo". You might also want to try "moobk driver".

moobk is free software: it is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more details.

`,
	SuggestionsMinimumDistance: 3,
	Version:                    "0.0.1",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		loadAliases()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	f := rootCmd.PersistentFlags()
	f.StringVarP(&aliasFile, "alias", "a", "", "specify yaml file stores repo alias, default to $HOME/.moobk.yaml (watchout if use with sudo)")
	f.StringVarP(&fs, "type", "t", "btrfs", "driver type, see 'moobk driver' for detail")
}

var aliasFile string
var fs string

// skip if failed
func loadAliases() {
	if aliasFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return
		}
		aliasFile = home + "/.moobk.yaml"
	}

	f, err := os.Open(aliasFile)
	if err != nil {
		return
	}
	defer f.Close()

	var a map[string]string
	dec := yaml.NewDecoder(f)
	if err = dec.Decode(&a); err != nil {
		return
	}

	moodrvs.SetAlias(a)
}
