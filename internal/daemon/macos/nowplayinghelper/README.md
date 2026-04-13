# Now Playing helper

Utility for macOS, used for displaying track information in the control center


## Build

```sh
swiftc main.swift -o tuitunes-nowplayinghelper -Xlinker -sectcreate -Xlinker __TEXT -Xlinker __info_plist -Xlinker Info.plist
```
