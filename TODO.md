# TODO

## Zip Functions

- [ ] Add Entries, support for glob
- [ ] Update Entries, support for glob, add if file does not exist already
- [x] [Delete Entries](#delete-entries)
  - [X] Add support for glob
- [ ] [Freshen Entries](#freshen-entries)
- [ ] [Junk Paths](#junk-paths)

### Delete Entries

```text
-d
--delete
    Remove (delete) entries from a zip archive. For example:

        zip -d foo foo/tom/junk foo/harry/\* \*.o

    will remove the entry foo/tom/junk, all of the files that start with foo/harry/, and all of the files that end with .o (in any path). Note that shell pathname expansion has been inhibited with backslashes, so that zip can see the asterisks, enabling zip to match on the contents of the zip archive instead of the contents of the current directory. (The backslashes are not used on MSDOS-based platforms.) Can also use quotes to escape the asterisks as in

        zip -d foo foo/tom/junk "foo/harry/*" "*.o"

    Not escaping the asterisks on a system where the shell expands wildcards could result in the asterisks being converted to a list of files in the current directory and that list used to delete entries from the archive.

    Under MSDOS, -d is case sensitive when it matches names in the zip archive. This requires that file names be entered in upper case if they were zipped by PKZIP on an MSDOS system. (We considered making this case insensitive on systems where paths were case insensitive, but it is possible the archive came from a system where case does matter and the archive could include both Bar and bar as separate files in the archive.) But see the new option -ic to ignore case in the archive.
```

### Freshen Entries

```text
-f
--freshen
    Replace (freshen) an existing entry in the zip archive only if it has been modified more recently than the version already in the zip archive; unlike the update option (-u) this will not add files that are not already in the zip archive. For example:

        zip -f foo

    This command should be run from the same directory from which the original zip command was run, since paths stored in zip archives are always relative.

    Note that the timezone environment variable TZ should be set according to the local timezone in order for the -f, -u and -o options to work correctly.

    The reasons behind this are somewhat subtle but have to do with the differences between the Unix-format file times (always in GMT) and most of the other operating systems (always local time) and the necessity to compare the two. A typical TZ value is ''MET-1MEST'' (Middle European time with automatic adjustment for ''summertime'' or Daylight Savings Time).

    The format is TTThhDDD, where TTT is the time zone such as MET, hh is the difference between GMT and local time such as -1 above, and DDD is the time zone when daylight savings time is in effect. Leave off the DDD if there is no daylight savings time. For the US Eastern time zone EST5EDT.
```

### Junk Paths

```text
-j
--junk-paths
Store just the name of a saved file (junk the path), and do not store directory names. By default, zip will store the full path (relative to the current directory).
```

## Unzip Functions

- [ ] [Freshen Entries](#freshen-entries-1)
- [ ] [Junk Paths](#junk-paths-1)

### Freshen Entries

```text
-f

freshen existing files, i.e., extract only those files that already exist on disk and that are newer than the disk copies. By default unzip queries before overwriting, but the -o option may be used to suppress the queries. Note that under many operating systems, the TZ (timezone) environment variable must be set correctly in order for -f and -u to work properly (under Unix the variable is usually set automatically). The reasons for this are somewhat subtle but have to do with the differences between DOS-format file times (always local time) and Unix-format times (always in GMT/UTC) and the necessity to compare the two. A typical TZ value is ''PST8PDT'' (US Pacific time with automatic adjustment for Daylight Savings Time or ''summer time'').
```

### Junk Paths

```text
-j

junk paths. The archive's directory structure is not recreated; all files are deposited in the extraction directory (by default, the current one).
```
