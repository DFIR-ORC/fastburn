# FastBurn
*Outil de qualification des résultats FastFind*

Ce petit outil est conçu pour traiter les résultats de FastFind en masse ou individuellement et permet de les qualifier.

Le programme peut lire une série d'archives résultant d'un FastFind, afficher un résumé et détecter les fichiers correspondant aux hachages ou aux noms référencés.

## Instructions de construction
### Unix
Un Makefile peut être utilisé pour gérer les compilations dans la plupart des environnements.
Tapez `make` pour générer le binaire `fbn` ou utilisez la ligne de commande 'go'.

```bash
go build -o fbn -v ./cmd/fastburn/  # executable: ./fbn
```

### Windows
Installer l'environnement MSYS2 à l'aide de powershell ou du programme d'installation
```bash
winget install msys2.msys2
C:\msys64\mingw64.exe
```

Utiliser le nouvel environnement MSYS2 de type `mingw64` et configurer la chaîne d'outils
```bash
pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-go
export GOROOT=/mingw64/lib/go
cd .../FastBurnt
go build -o fbn -v ./cmd/fastburn/  # executable: .\fbn.exe
```

Ajouter le support de l'outil `make` (optional)
```bash
pacman -S make
make
```

## Usage
```
fbn [--debug|--trace] [--whitelist <whitelist.csv>] [--blacklist <blacklist.csv>]
[-output <output file>] [-computers list.csv]  [-html] <fichiers>


Detailed usage:

  -blacklist string
        Indiquer un fichier CSV contenant les drapeaux à supprimer des résultats.
  -computers string
        Spécifier le nom du fichier de listage des ordinateurs
  -debug
        Activer le mode débogage
  -html
        Activer la sortie HTML
  -output string
        Spécifier le fichier de sortie
  -stats string
        Spécifier le fichier de sortie pour les statistiques
  -timeline string
        Spécifier le fichier de sortie pour la "timeline"
  -trace
        Activer le mode traces
  -version
        Afficher la version
  -whitelist string
        Spécifier un fichier CSV contenant les drapeaux à mettre en évidence dans les résultats
  <fichiers>
        Un ou plusieurs fichier(s) d'archive 7z ou un dossier les contenant
```

Exemple d'exécution de la ligne de commande
```log
 ./fbn Resultats/*7z
  INFO[0000] File 'Resultats/ORC_WorkStation_DESKTOP-LCINJKL_FastFind.7z', Hostname DESKTOP-LCINJKL matches: 0
  INFO[0000] File 'Resultats/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z', Hostname DESKTOP-LCINQGJ matches: 4
  INFO[0000] File 'Resultats/ORC_WorkStation_DESKTOP-JKLNQGJ_FastFind.7z', Hostname DESKTOP-JKLNQGJ matches: 1
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\FalsePositive\foo.dll] : backdoor SOLARBURST - Archive 'Resultats/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\Webshell\Aie.dll] : webshell SUPERNOVA - Archive 'Resultats/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\TruePositive\Solarwinds.Orion.Core.Businesslayer.dll] : installation SolarWinds Orion, backdoor SOLARBURST - Archive 'Resultats/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\TruePositive\Solarwinds.Other.Businesslayer.dll] : backdoor SOLARBURST - Archive 'Resultats/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-JKLNQGJ [\Users\user\Documents\SolarWindsSain\Solarwinds.Orion.Core.Businesslayer.dll] : installation SolarWinds Orion - Archive 'Resultats/ORC_WorkStation_DESKTOP-JKLNQGJ_FastFind.7z'
  INFO[0000] Matches exported to '2020-12-31T00_09_21Z-fastfound.csv'
```

Voir aussi [Utilisation avancée](usage.md)

## License

Le contenu de ce dépôt est disponible sous licence LGPL2.1+, tel qu'indiqué [ici](LICENSE).
Le nom DFIR-ORC, Fastburn, fbn et le logo associé appartiennent à l'ANSSI, aucun usage n'est permis sans autorisation expresse.

---

The contents of this repository is available under [LGPL2.1+ license](LICENSE).
The name DFIR-ORC, Fastburn, fbn and the associated logo belongs to ANSSI, no use is permitted without express approval.
