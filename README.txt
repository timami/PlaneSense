			      ━━━━━━━━━━━━
			       PLANESENSE
			      ━━━━━━━━━━━━


Table of Contents
─────────────────

1 Purpose
2 Dependencies
.. 2.1 Supported OSes
.. 2.2 Requirements
.. 2.3 Go Set up
.. 2.4 To get the libsox library to work
3 Running
4 Important files





1 Purpose
═════════

  To enable pilots to hear planes around them.


2 Dependencies
══════════════

2.1 Supported OSes
──────────────────

  Currently only gnu/linux is supported.  Any version of gnu/linux made
  in this decade should work.


2.2 Requirements
────────────────

  • A recent version of Go must be installed. We used 1.6.2.
  • An internet connection.


2.3 Go Set up
─────────────

  • Ensure that your go path and bin are set via the environment
    variables.
  • For example, adding the following to ones .bashrc will
    work. (Assuming one is using bash as a shell.)
    ┌────
    │ export GOPATH=$HOME/go
    │ export PATH=$PATH:/home/<user>/go/bin/
    └────


2.4 To get the libsox library to work
─────────────────────────────────────

  • enable source repos in your distribution.
    • This can be done either be uncommenting the source in
      /etc/apt/sources.list, and running sudo apt-get update
    • Or using the gui.
  • Then install libsox
    • sudo apt-get install libsox-dev libsox2
  • sudo apt-get build-deps sox
  • sudo apt-get source sox
  • cd into the dirctory apt-get created (It'll be sox-XXXX).
  • sudo ./confiure
  • sudo make -s
  • sudo make install
  • go get -u github.com/krig/go-sox


3 Running
═════════

  ┌────
  │ % cd ./GoPositionalAudio 
  │ % go run sound.go
  └────


4 Important files
═════════════════

  Inside the GoPositionalAudio folder
  sound.go: Main driver and implementation
