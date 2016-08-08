seal v0
=======

**Note: Not all spec features may be in the reference implementation yet.**

"seal" is a simple container format that embeds the information needed to
check for corruption (hashes) and/or authenticity information (signatures)
into the file itself, allowing verification across filesystems, networking
protocols, and operating systems.

The goal of seal is to reduce the barrier to checking for file corruption and
file authenticity, and therefore make storing and transferring digital files
provably reliable.

seal also aims to match the integrity and security guarantees offered by signify
signed files, while having acceptable performance and being simple to use and
understand.

Header Identification
---------------------

The seal magic number is encoded in UTF-8 in the first 4 bytes of the file:


        0                      31
        +-----+-----+-----+-----+
        | 'S' | 'L' | '%' | 'v' |
        +-----+-----+-----+-----+
        | '0' | '{' | ..........


`SL%v` is the magic number for seal files and can be used to identify them.

The byte immediately following the literal `v` identifies the hex-encoded major
version number of the file. 

Header Structure
----------------

    SL%v0{variant:<claim>}

Variants and Claims
-----------------------------

There are two seal variants with different properties. The variant determines
how the claim is generated and interpreted.

### sha512

For the `sha512` variant, the claim is the hex-encoded sha512 hash of the given
file. It may optionally be truncated to the first `x` bytes, where `1 ≤ x < 64`.

The following short form may be used:

    SL%v0{<claim>}

Example, truncated to 128 bits:

    SL%v0{53331cbf3149b47ba0be481c1cfd61d6}

### signify

The `signify` variant targets compatibility with OpenBSD's signify tool for
signing and verification. The claim is the base64-encoded signature generated
by signify on a given file.

    SL%v0{signify:<signature>}

Example:

    SL%v0{signify:RWRMdbgIymBjpBudT1rr/hQivikPSRRVgTTj+0u+t5Lg1zGbz28HaseMefb9XbycbXGT0Lfm0KOc5vZbi8cydUtIpP4Txqe5GQs=}

vim: tw=80 et sw=4 sts=4
