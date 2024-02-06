# xbm❤braille

XBM to Braille converter (the ❤ is silent)

# installation

```console
: go install chipaca.com/xbmbraille@latest
```

# commandline options

```console
: xbmbraille
Usage: xbmbraille [options] {-|filename}...
  -c	clear the terminal before printing each image
  -d duration
    	wait this much time after printing each image
  -n	negate (invert) image
  -p	print the filename before printing each image
```

# example usage

```console
: convert +dither -font Noto-Emoji -pointsize 64 label:🤩 -trim XBM:- | xbmbraille -n -
⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣤⣴⣶⠶⠶⠶⠶⢶⣶⣤⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢀⣴⡾⠟⠋⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠛⠿⣶⣄⡀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⣠⣾⠟⣿⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣴⡟⢿⣦⡀⠀⠀⠀⠀
⠀⠀⢀⣼⠟⠁⢠⣿⣿⣦⢀⣀⣀⣀⠀⠀⠀⠀⢀⣀⣀⣀⢀⣾⣿⣿⠀⠙⢿⣆⠀⠀⠀
⠀⢀⣾⣯⣤⣴⣾⣿⣿⣿⣿⣿⣿⠟⠀⠀⠀⠀⠈⢿⣿⣿⣿⣿⣿⣿⣶⣦⣤⣿⣧⠀⠀
⠀⣾⠏⠙⠻⣿⣿⣿⣿⣿⣿⣿⠁⠀⠀⠀⠀⠀⠀⠀⢹⣿⣿⣿⣿⣿⣿⡿⠟⠉⢻⣇⠀
⢸⡟⠀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣇⠀⠀⠀⠀⠀⠀⢀⣾⣿⣿⣿⣿⣿⡏⠀⠀⠀⠈⣿⡀
⣿⡇⠀⠀⠀⠀⣿⠟⠋⠀⠈⠙⠛⠀⠀⠀⠀⠀⠀⠘⠛⠉⠁⠀⠙⢿⡇⠀⠀⠀⠀⢿⡇
⣿⡇⠀⠀⠀⠀⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣸⡇
⢻⡇⠀⠀⠀⢀⣤⣤⣀⣀⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⣠⣤⣄⠀⠀⠀⠀⣿⡇
⠸⣿⠀⠀⠀⠘⣧⣄⣈⣉⡉⠉⠛⠛⠛⠛⠛⠛⠛⠛⠋⠉⢉⣉⣀⣠⡿⠀⠀⠀⢰⡿⠀
⠀⢹⣧⠀⠀⠀⠙⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠁⠀⠀⢠⣿⠃⠀
⠀⠀⠹⣷⡀⠀⠀⠀⠙⠿⣿⣿⡿⠟⠛⠋⠉⠛⠛⠿⣿⣿⡿⠟⠉⠀⠀⠀⣠⡿⠃⠀⠀
⠀⠀⠀⠘⢿⣦⡀⠀⠀⠀⠀⠉⠛⠓⠶⠶⠶⠶⠖⠛⠋⠁⠀⠀⠀⠀⢀⣼⠟⠁⠀⠀⠀
⠀⠀⠀⠀⠀⠙⠿⣦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣾⠟⠁⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠈⠙⠻⢶⣦⣤⣄⣀⣀⣀⣀⣀⣠⣤⣶⠾⠟⠋⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠉⠉⠛⠛⠛⠉⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
```
