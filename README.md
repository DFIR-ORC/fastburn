# FastBurn
*FastFind Results Qualification Tool* - [Version française](doc/fr/README.md)

[![LGPL licensed][img-license]](./LICENSE)

This small tool is designed to process FastFind results either in bulk or individually and allows for their qualification.

The program can read a series of archives resulting from a FastFind, display a summary, and detect files corresponding to referenced hashes or names.

## Build instructions
### Unix
A Makefile can be used to manage compilations on most environments.
Type `make` to generate the `fbn` binary or use the 'go' command line.

```
go build -o fbn -v ./cmd/fastburn/  # output: ./fbn
```

### Windows
Install MSYS2 enviroment using powershell or installer
```
winget install msys2.msys2
```

Use the new `MSYS2` environment of type `mingw64` and setup the toolchain
```
C:\msys64\msys2_shell.cmd -defterm -here -no-start -mingw64  # or 'MSYS2 MINGW64' from the menu
pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-go make
export GOROOT=/mingw64/lib/go  # or open another mingw64 terminal
cd .../FastBurn
go build -o fbn.exe -v ./cmd/fastburn/  # output: .\fbn.exe
```

Optional: setup development environment
```
winget install microsoft.visualstudiocode
# open another terminal to refresh environment variables
code --install-extension golang.go
C:\msys64\msys2_shell.cmd -defterm -here -no-start -mingw64
"${LOCALAPPDATA}/Programs/Microsoft VS Code/bin/Code.cmd" --install-extension golang.go
"${LOCALAPPDATA}/Programs/Microsoft VS Code/bin/Code.cmd" --add FastBurn/
```

## Usage
```
fbn [--debug|--trace] [--whitelist <whitelist.csv>] [--blacklist <blacklist.csv>]
[-output <output file>] [-computers list.csv]  [-html] <7zArchive1 ... n>

Detailed usage:

  -blacklist string
        Specify a CSV file containing flags to suppress from the results
  -computers string
        Specify computers listing filename
  -debug
        Enable debug mode
  -html
        Enable HTML output
  -output string
        Specify output filename
  -stats string
        Specify statistics filename
  -timeline string
        Specify a filename for timeline output
  -trace
        Enable trace mode
  -version
        Show version and exit
  -whitelist string
        Specify a CSV file containing flags to highligth in the results
```

Example of command line execution
```log
 ./fbn Results/*7z
  INFO[0000] File 'Results/ORC_WorkStation_DESKTOP-LCINJKL_FastFind.7z', Hostname DESKTOP-LCINJKL matches: 0
  INFO[0000] File 'Results/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z', Hostname DESKTOP-LCINQGJ matches: 4
  INFO[0000] File 'Results/ORC_WorkStation_DESKTOP-JKLNQGJ_FastFind.7z', Hostname DESKTOP-JKLNQGJ matches: 1
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\FalsePositive\foo.dll] : backdoor SOLARBURST - Archive 'Results/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\Webshell\Aie.dll] : webshell SUPERNOVA - Archive 'Results/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\TruePositive\Solarwinds.Orion.Core.Businesslayer.dll] : installation SolarWinds Orion, backdoor SOLARBURST - Archive 'Results/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\TruePositive\Solarwinds.Other.Businesslayer.dll] : backdoor SOLARBURST - Archive 'Results/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-JKLNQGJ [\Users\user\Documents\SolarWindsSain\Solarwinds.Orion.Core.Businesslayer.dll] : installation SolarWinds Orion - Archive 'Results/ORC_WorkStation_DESKTOP-JKLNQGJ_FastFind.7z'
  INFO[0000] Matches exported to '2020-12-31T00_09_21Z-fastfound.csv'
```

See also [Advanced usage](doc/en/usage.md)

## License

The contents of this repository is available under [LGPL2.1+ license](LICENSE).
The name DFIR-ORC, Fastburn, fbn and the associated logo belongs to ANSSI, no use is permitted without express approval.

---

Le contenu de ce dépôt est disponible sous licence LGPL2.1+, tel qu'indiqué [ici](LICENSE).
Le nom DFIR-ORC, Fastburn, fbn et le logo associé appartiennent à l'ANSSI, aucun usage n'est permis sans autorisation expresse.

[img-license]: https://img.shields.io/github/license/DFIR-ORC/fastburn
