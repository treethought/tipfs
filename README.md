# tipfs

tipfs is an ugly little ipfs TUI client.


[![asciicast](https://asciinema.org/a/435115.svg)](https://asciinema.org/a/435115)


## Install

``` sh
go get github.com/treethought/tipfs
```

## Features
- [x] Browse, viewing Mutable File System
- [x] Exploring  DAG nodes
- [x] Viewing peers
- [x] View supported content in terminal
- [x] Opening/copying CID in browser

## Keybindings

| key   | action                                 |
|-------|----------------------------------------|
| TAB   | Switch focus between panels            |
| o     | Open in browser                        |
| y     | Copy selected items CID                |
| j     | Move selection up                      |
| k     | Move selection down                    |
| Enter | Select file to inspect                 |
| g     | Go to top of panel                     |
| G     | Go to bottom of panel                  |
| 1     | Switch to files mode                   |
| 2     | Switch to peers mode                   |

    





