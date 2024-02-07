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
  -ms-gap int
    	Adds a minimum gap between events, any events received less than this many milliseconds after the previous event will be ignored. Useful if writes are infrequent but duplicate from some programs.
  -split-then
    	Splits command string by the "then" operator (";"), using each new string as a command. Successive commands will run regardless the previous's success
```

Have fun :)

