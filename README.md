shield
======

**Warning: This project is under heavy development. Breaking changes will most
likely be introduced in the spec and the tool itself until 1.0.**

`shield` is a container format that lets you check for file corruption without
dealing with separate checksum files.

Distilled, shield is just a file with a prepended sha512 hash.

The header has the following format:

    SHD%v0{cf83e135}

Above is a fully valid shielded file (ending in a newline), where the
contents are completely empty.

Check the examples folder for more.

( I was disturbingly able to trash many bits within squirrel.jpg before the
corruption was visible when I opened the file. Hence why this exists. )

Installation
------------

    ; go get -v github.com/crasm/shield
    ; go install -v github.com/crasm/shield

Usage
-----

Once installed, run shield with no arguments.

    ; shield

Examples
--------

    # Extracts to squirrel.jpg.
    ; shield -X squirrel.jpg.shd 

    # Prints to stdout. (Be careful with binary.)
    ; shield -C < LICENSE

    # Shields the text and then extracts it.
    ; echo 'shield pipe!' | shield -C | shield -X

Abstract goals
-------------

1. Make checking file integrity easy and easier to automate, so that it can be
   done across filesystems, networking protocols, and operating systems.
2. Keep the format simple enough to generate, extract, and verify file contents
   "by hand" with basic \*nix tools.

Concrete goals
-------------

- [ ] Signify support.
- [ ] Backup client that uses shield to verify integrity while copying files.
- [ ] Browser plugins and apps for automatic verification and extraction of downloads.
- [ ] HTTP middleware for go.

Manual shield header generation
-------------------------------

To generate a shield file by hand:

    ; sha512sum LICENSE
    53331cbf3149b47ba0be481c1cfd61d60282ce13652909a17a25626...  LICENSE
    ; echo 'SHD%v0{<paste the hash>}' > LICENSE.shd
    ; cat LICENSE >> LICENSE.shd

To check a shielded file by hand:

    ; head -n 1 LICENSE.shd
    ; tail -n +2 LICENSE.shd > LICENSE
    ; sha512sum LICENSE
      <compare the hashes starting from the left>


vim: tw=80 et sw=4 sts=4
