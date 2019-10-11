# armulator

The goal of this project is to facilitate development of ARM applications. It
is not meant to be a full fledged VM but rather, an easy way to test ARM
compiled binaries. An arbitrary objective would be to run a Go compiler
that was build with `GOARCH=arm64` and use it to compile itself again
within this emulator.

Basically, this means:
* broad support of logic, arithmetic, and memori instructions
* minimal support for vector instruction
* minimal support for syscalls

Non goals:
* precise hardware emulation
* bootloader, os, customization
    * idealy we just forward syscalls to the host
    * we will probably have to implement something to read ELF files
    * we will probably have to implement our version of `sys_brk` or `mmap`

These things would not necessarily be helpful with C where system library and
headers are needed. But, for statically linked programs (like Go programs)
this emulator should be handy.

