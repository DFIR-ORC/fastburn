# Fastburnt - outil de qualification des retours FastFind


## Objet

Ce petit outil est destiné à permettre de traiter des retours FastFind en masse ou à l'unitée et permettre de les qualifier.

Le programme permet de lire une série d'archive résultat d'un FastFind, d'en afficher une synthèse et d'y détecter les fichiers correspondant à des hash ou des noms référencés.


## Build

Voir le fichier `doc/Build.md`

## Execution

### Usage

```
   ./fastburnt_cli [--debug|--trace] [--whitelist <whitelist.csv>] [--blacklist <blacklist.csv>] [-output <output file>] [-computers list.csv] <7zArchive1 ... n>
```

### Exemple d'exécution en ligne de commande:

```
 ./fastburnt_cli Resultats/*7z
  INFO[0000] File 'Resultats/ORC_WorkStation_DESKTOP-LCINJKL_FastFind.7z', Hostname DESKTOP-LCINJKL matches: 0
  INFO[0000] File 'Resultats/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z', Hostname DESKTOP-LCINQGJ matches: 4
  INFO[0000] File 'Resultats/ORC_WorkStation_DESKTOP-JKLNQGJ_FastFind.7z', Hostname DESKTOP-JKLNQGJ matches: 1
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\FalsePositive\Pouet.dll] : backdoor SOLARBURST - Archive 'Resultats/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\Webshell\Aie.dll] : webshell SUPERNOVA - Archive 'Resultats/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\TruePositive\Solarwinds.Orion.Core.Businesslayer.dll] : installation SolarWinds Orion, backdoor SOLARBURST - Archive 'Resultats/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-LCINQGJ [\Users\user\Documents\TruePositive\Solarwinds.Other.Businesslayer.dll] : backdoor SOLARBURST - Archive 'Resultats/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z'
  WARN[0000] - DESKTOP-JKLNQGJ [\Users\user\Documents\SolarWindsSain\Solarwinds.Orion.Core.Businesslayer.dll] : installation SolarWinds Orion - Archive 'Resultats/ORC_WorkStation_DESKTOP-JKLNQGJ_FastFind.7z'
  INFO[0000] Matches exported to '2020-12-31T00_09_21Z-fastfound.csv'
```

Pour plus de détail `Doc/Usage.md`
