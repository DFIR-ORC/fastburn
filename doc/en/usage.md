# Using the FastBurn tool

## Operation

The FastBurn tool reads a list of 7zip archives containing results from the FastFind tool of the DFIR-ORC suite.
The results are decompressed, decoded and displayed.
Each entry can be compared against:

* a *whitelist* of markers of particular interest
* a *blacklist* of markers to ignore

At the end of execution, a CSV file containing detailed analysis data is generated.
fbn is designed to be used from the command line (Windows or Linux) and generates output to standard output. The CSV file is automatically generated at the end of execution.

### Usage

```
   ./fbn [--debug|-trace]
      [-whitelist <whitelist.csv>] [-blacklist <blacklist.csv>]
      [-output <output file>] [-computers <machine list file>]
      [-timeline <timeline file>] [-html]
      <files>
```

#### Output files

* `<7zArchive1...n>`: list of 7zip result files or directories containing these files. If an entry is a directory, it will be recursively traversed to search for 7zip files. Only 7zip files containing search results will be processed.
* `<whitelist.csv>`: markers to be highlighted during post-processing
* `<blacklist.csv>`: marker file to be ignored during post-processing
* `<output file.csv>`: name of output file for matched results

**Note**  If a file is not a valid archive, it is ignored, but processing continues.

#### Option details

* `debug`: activates the second level of traceability on STDERR
* `trace`: activates the maximum level of traceability on STDERR
* `whitelist`: allows you to specify a file of markers to be highlighted
* `blacklist`: allows you to specify a marker file to be excluded from results
* `output`: allows you to force the name of the results file
* `computers`: allows you to force the name of the file summarizing the list of machines found in the processed archives
* `timeline`: allows you to force the timeline file name
* `html`: enable the output of an HTML report
* `files`: list of 7z archive files or folder containing them

The format of whitelist and blacklist files is the same. It is documented in the “Flags file format” section below.

### Example

Command-line execution under Linux

```log
 ./fbn Results
  INFO[0000] File 'Results/ORC_WorkStation_DESKTOP-LCINJKL_FastFind.7z', Hostname DESKTOP-LCINJKL matches: 0
  INFO[0000] File 'Results/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z', Hostname DESKTOP-LCINQGJ matches: 4
  INFO[0000] File 'Results/ORC_WorkStation_DESKTOP-JKLNQGJ_FastFind.7z', Hostname DESKTOP-JKLNQGJ matches: 1
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\FalsePositive\Pouet.dll] : backdoor SOLARBURST - Archive 'Results/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\Webshell\Aie.dll] : webshell SUPERNOVA - Archive 'Results/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\TruePositive\Solarwinds.Orion.Core.Businesslayer.dll] : installation SolarWinds Orion, backdoor SOLARBURST - Archive 'Results/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\TruePositive\Solarwinds.Other.Businesslayer.dll] : backdoor SOLARBURST - Archive 'Results/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-JKLNQGJ [\Users\user\Documents\SolarWindsSain\Solarwinds.Orion.Core.Businesslayer.dll] : installation SolarWinds Orion - Archive 'Results/ORC_WorkStation_DESKTOP-JKLNQGJ_FastFind.7z'
  INFO[0000] Matches exported to '2020-12-31T00_09_21Z-fastfound.csv'
```

### Output file format

#### Inspected machines list file

The inspected machines file provides a summary of the results examined.

The data is in CSV format, separated by the `,` character and where strings are enclosed in quotation marks: `"`. If no name is specified on the command line, the default file name is `<timestamp>-fastburn_computers.csv`.

Each detection is the subject of an information line.

Each line consists of the following fields:

* `Ignore`: set to `true` if the entry corresponds to a blacklist item, `false` otherwise
* `Computer`: name of the machine on which the entry was detected
* `ComputerRole`: machine function
* `ComputerOS`: machine operating system
* `ORCVersion`: version of DFIR-ORC on which the FastFind tool used was based
* `MatchType`: type of detection that determined FastFind's input selection
* Software
* `Reason`: criteria for which the entry was selected by FastFind
* `Filename`: absolute path of the file
* `AltName`: alternative file name
* `RegKey`: registry key name
* `RegType`: registry key type
* `RegValue`: registry key value
* `FileSize`: file size in number of bytes
* `MD5`: hexadecimal encoded MD5 digest
* `SHA1`: hexadecimal encoded SHA1 digest
* `SHA256`: hexadecimal-encoded SHA256 condensate holding the result
* `FileCreation`: file creation date
* `FileLastModification`: date of last file modification
* `FileLastEntryChange`: date file meta-information last changed
* `FileLastAccess`: date of last access to file
* `FilenameCreation`: creation date of the `$FN` entry in the file
* `FilenameLastModification`: last modification date of the `$FN` file entry
* `FilenameLastEntryChange`: date of last meta-information change to the `$FN` file entry
* `FilenameLastAccess`: date of last access to the meta-information of the `$FN` file entry
* `AltFilenameCreation`: alternative file name creation date
* `AltFilenameLastModification`: date of last modification of the alternative file name
* `AltFilenameLastEntryChange`: date of last modification of alternative filename meta information
* `AltFilenameLastAccess`: date of last access to alternate filename
* `VolumeID`: identifier of the file system volume in which the search is performed
* `SnapshotID`: identifier of the file system snapshot in which the search is performed
* `ArchiveName`: path of the archive containing the result

#### Timeline file

The file is in MACB format, commonly used by forensic investigation tools.

Data is in CSV format, separated by the `,` character and where strings are enclosed in quotation marks: `"`.

If no name is specified on the command line, the default file name is `<timestamp>-fastburn_timeline.csv`.

This format puts all information relating to the change of a file or system entity on a single line.

Entries in this format are intended to be integrated into a *timeline* reconstructing action sequences.

A good description of the format can be found here <https://andreafortuna.org/2017/10/06/macb-times-in-windows-forensic-analysis/>

For each entry the following information is generated:

* `Timestamp`: date of change
* `SI_MACB`: MACB format entry change code
* `FN_MACB`: change code for `$FILENAME` entry in MACB format
* `ComputerName`: name of the machine on which the entry is identified
* `File`: name of the file or system entry
* `ParentName`: name of the directory to which the entry is associated
* `FullName`: full path of the entry
* `Extension`: file extension (last characters after `.`)
* `SizeInBytes`: entry size in number of bytes
* `CreationDate`: entry creation date
* `LastModificationDate`: date of last modification
* `LastAccessDate`: last access date
* `LastAttrChangeDate`: date attributes last changed
* `FileNameCreationDate`: creation date of `$FN` entry
* `FileNameLastModificationDate`: date entry was last modified `$FN` * `FileNameLastAccessDate` date entry was last changed
* `FileNameLastAccessDate`: date entry was last accessed `$FN` * `FileNameLastAccessDate` date entry was last accessed
* `FileNameLastAttrModificationDate`: date of last modification to attributes of entry `$FN`
* `MD5`: hexadecimal-encoded MD5 digest
* `SHA1`: hexadecimal-encoded SHA1 digest
* `SHA256`: hexadecimal-encoded SHA256 digest
* `Reason`: reason for input selection by FastFind
* `ArchiveName`: archive name


##### Date format

All dates are expressed in the format `YYYY-MM-DD HH:MN:SS.MS`.

These are in the UTC time zone.

##### MACB format

The MACB format is a 4-character string that defines changes in the input's meta-information.

Each of these characters can be

* `M` for *Modified* last modification date
* `A` for *Accessed* date last accessed
* `C` for *Changed* date the `$MFT` entry was last changed
* `B` for *Birth* date entry created

If the change in question is not applicable to the entry, the letter is replaced by the `.` character.

#### Output file example

```csv
"2020-12-23 00:09:16.944";"..C.";"..C.";"DESKTOP-LCINQGJ";"Solarwinds.Core.Businesslayer.dll";"\Users\user\Documents\TruePositive";"\Users\user\Documents\TruePositive\Solarwinds.Core.Businesslayer.dll";"dll";"1028072";"2020-12-23 00:09:32.117";"2020-12-23 16:39:01.000";"2020-12-23 00:09:38.147";"2020-12-23 00:11:44.288";"2020-12-23 00:09:32.117";"2020-12-23 16:39:01.000";"2020-12-23 00:09:32.117";"";"846E27A652A5E1BFBD0DDD38A16DC865";"D130BD75645C2433F88AC03E73395FBA172EF676";"CE77D116A074DAB7A22A0FD4F2C1AB475F16EEC42E1DED3C0B0AA8211FE858D6";"Size=1028072, SHA256=CE77D116A074DAB7A22A0FD4F2C1AB475F16EEC42E1DED3C0B0AA8211FE858D6";"share/samples/FastFind/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z"
```

#### Statistics file

The statistics file summarizes the main metrics of the analysis.

Data is in CSV format, separated by the `,` character and where strings are enclosed in quotation marks: `"`.

If no name is specified on the command line, the default file name is `<timestamp>-fastburn_stats.csv`.

The statistics file gives the number of results according to :

* machines
* machine types
* operating system
* Windows domain
* detection rule
* file name
* condensate
* file size
* day created and last modified
* month of creation and last modification

### Flags file format

The “whitelist”/“blacklist” file is a comma-separated CSV using double quotes as string delimiters.

It must contain the following columns (even if empty):

* `sha256` - condensate
* `sha1` - condensate
* `md5` - condensate
* `file_re` - REGEXP to be applied to the file path
* `description` - detection description

A detection corresponds to a “flag” line if any of the markers match.
The `file_re` criterion is applied to the `Fullname` field of the result file.

#### Examples of regular expressions

Match all DLLs in the directory `Program Files (x86)\Adobe\Acrobat Reader DC\Reader\AcroCEF``.
```
^Program Files (x86)\Adobe\Acrobat Reader DC\Reader\AcroCEF\.*\.dll$
```

Match executables in the WinRAR installation directory regardless of case.
```
(?i)^\Program Files \(x86\)\WinRAR\\.*exe$
```

### Note on use when processing large quantities of files

The generation of traces to the terminal has a significant impact on performance.
It is advisable to redirect error output to a file when processing tens or hundreds of thousands of files.

Example:

```sh

# Will be slow if 'Results' is a tree containing many result files.
 ./fbn Results

# To go faster and to generate detailed execution logs in the file 'fastburn.log'
 ./fbn -debug Resultats 2> fastburn.log
```

