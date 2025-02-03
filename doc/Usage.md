# Utilisation de l'outil FastBurnt


## Fonctionnement

L'outil Fasturnt lit une liste d'archives 7zip contenant le résultat d'une exécution de l'outil FastFind de la suite DFIR-ORC.
Les résultats sont décompressés, décodés et affichés.

Chaque entrée peut aussi être comparée à :
* une *whitelist* de marqueurs d'intéret particulier
* une *blacklist* de marqueurs à ignorer


Au terme de l'exécution, un fichier CSV contenant le détail des données analysées est généré.

Deux déclinaisons de l'outil d'ont fournies
* `fastburnt_cli` est destiné à être utilisé en ligne de commande (Windows ou Linux) et génère les affichages sur les sorties standard. Le fichier CSV est automatiquement généré au terme de l'exécution.
* `fastburnt_ui` permet une exécution via interface graphique. La liste des fichiers à traiter peut être éditée interactivement. L'export CSV est déclenché par un choix utilisateur.


### Usage

```
   ./fastburnt_cli [--debug|--trace] [--whitelist <whitelist.csv>] [--blacklist <blacklist.csv>] [-output <output file>] [-computers list.csv] <7zArchive1 ... n>
```


ou `<7zArchive1...n>` est une liste de fichiers de résultats 7zip ou de répertoires contenant ces fichiers. Si une entrée est un répertoire, celui-ci va être parcouru récursivement pour y rechercher les fichiers 7zip. Seuls les fichiers 7zip contenant un résultat de recherche seront traités.

ou `<whitelist.csv>` est un fichier de marqueurs à mettre en valeur lors du post-traitement

ou `<blacklist.csv>` est un fichier de marqueurs à ignorer lors du post-traitement

ou `<output file.csv>` est le nom du fichier de sortie des résultats matchés

Si un fichier n'est pas une archive valide, il est ignoré mais le traitement continue.


Détail des options:

* `debug` active le second niveau de traçabilité sur STDERR
* `trace` active le niveau maximal de traçabilité sur STDERR
* `whitelist` permet de spécifier un fichier de marqueurs à mettre en évidence
* `blacklist` permet de spécifier un fichier de marqueurs à exclure des résultats
* `output` permet de forcer le nom du fichier de résultats
* `computers` permet de forcer le nom de fichier récapitulant la liste les machines trouvées dans les archives traitées

Le format des fichiers de liste blanches et noire est le même. Il est documenté ci-dessous dans la section "Format de fichier de Flags" ci dessous.

### Exemple

Exécution en ligne de commande sous Linux

```
 ./fastburnt_cli Resultats
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

### Format du fichier de Flags

Le fichier de "whitelist"/"blacklist" est un CSV séparé par des virgules et utilisant des doubles quotes comme délimiteur de chaine.

Il doit contenir les colonnes suivantes (même vides):

* `sha256`      - condensat
* `sha1`        - condensat
* `md5`         - condensat
* `file_re`     - REGEXP à appliquer sur le chemin de fichier
* `description` - description de la détection

Une détection correspond a une ligne de "flag" si n'importe lequel des marqueurs correspond.
Le critère `file_re` est appliqué sur le champ `Fullname` du fichier de résultat.

Exemples d'expressions régulières:

* Matcher toutes les DLL du répertoire `Program Files (x86)\Adobe\Acrobat Reader DC\Reader\AcroCEF`
```
^Program Files \(x86\)\\Adobe\\Acrobat Reader DC\\Reader\\AcroCEF\\.*\.dll$
```

* Matcher les exécutables du répertoire d'installation de WinRAR indépendamment de la casse

```
(?i)^\\Program Files \(x86\)\\WinRAR\\.*exe$
```

### Note sur l'usage pour traiter des grandes quantités de fichiers

La génération de traces vers le terminal a un impact non négligeable sur les performances.
Il est conseillé de rediriger la sortie d'erreur vers un fichiers lorsque l'on traite des dizaines ou centaines de milliers de fichiers.

Exemple:
```

 ./fastburnt_cli Resultats

# Sera lent si 'Resultats' est une arborescence contenant
# beaucoup de fichiers de résultats.
# Pour aller plus vite:

 ./fastburnt_cli -debug Resultats 2> fastburnt.log

# et va générer les journaux d'exécution détaillés dans le
# fichier 'fastburnt.log'

```

