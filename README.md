# Watchman

A simple program for running commands on file write events, its super small and lightweight. I'm using this to reload my `wpaperd` config right now, which can be done like so:

```fish
./watchman -file ~/.config/wpaperd/wallpaper.toml -command 'killall wpaperd; wpaperd' -split-then
```

Heres the full usage:

```
Usage of ./watchman:
  -command string
    	The command to run on file change events
  -file string
    	The file to watch for changes
  -split-and
    	Splits command string by the "and" operator ("&&"), using each new string as a command. Successive commands will only run upon the previous's success
  -split-then
    	Splits command string by the "then" operator (";"), using each new string as a command. Successive commands will run regardless the previous's success
```

Have fun :)

