# Utilisation de l'outil FastBurn

## Fonctionnement

L'outil FastBurn lit une liste d'archives 7zip contenant le résultat d'une exécution de l'outil FastFind de la suite DFIR-ORC.
Les résultats sont décompressés, décodés et affichés.

Chaque entrée peut aussi être comparée à :
* une *whitelist* de marqueurs d'intéret particulier ;
* une *blacklist* de marqueurs à ignorer.

Au terme de l'exécution, plusieurs fichiesr CSV contenant le détail des données analysées sont générés.

`fbn` est destiné à être utilisé en ligne de commande (Windows ou Linux) et génère les affichages sur les sorties standards. Les fichiers CSV sont automatiquement générés au terme de l'exécution.


### Usage

```
   ./fbn [-debug|-trace]
      [-whitelist <whitelist.csv>] [-blacklist <blacklist.csv>]
      [-output <output file>] [-computers <machine list file>]
      [-timeline <timeline file>] [-html]
      <files>
```

#### Fichiers en entrée
* `<files>` : liste de fichiers 7zip ou de répertoires contenant ces fichiers. Si une entrée est un répertoire, celui-ci va être parcouru récursivement pour y rechercher les fichiers 7zip. Seuls les archives contenant un résultat de recherche sont traitées.
* `<whitelist.csv>` : marqueurs à mettre en valeur lors du post-traitement.
* `<blacklist.csv>` : fichier de marqueurs à ignorer lors du post-traitement.
* `<output file.csv>` : nom du fichier de sortie des résultats trouvés.

**Note** Si un fichier n'est pas une archive valide, il est ignoré mais le traitement continue.

#### Détail des options

* `debug` active le second niveau de traçabilité sur STDERR
* `trace` active le niveau maximal de traçabilité sur STDERR
* `whitelist` permet de spécifier un fichier de marqueurs à mettre en évidence
* `blacklist` permet de spécifier un fichier de marqueurs à exclure des résultats
* `output` permet de forcer le nom du fichier de résultats
* `computers` permet de forcer le nom de fichier récapitulant la liste les machines trouvées dans les archives traitées
* `timeline` permet de forcer le nom de fichier de la timeline
* `html` active la sortie au format HTML

Le format des fichiers de liste blanche et noire est le même. Il est documenté ci-dessous dans la section "Format de fichier de Flags".


### Exemple

Exécution en ligne de commande sous Linux

```log
 ./fbn Resultats
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

### Format des fichiers de sortie

#### Fichier de liste de machines inspectées

Le fichier des machines inspectées fournit un récapitulatif des résultats examinés.

Les données sont au format CSV, séparées par le caractère ',' et où les chaines de caractères sont encadrées par des guillemets '"'.

Si aucun nom n'est spécifié sur la ligne de commande, le nom par défaut du fichier est `<timestamp>-fastburn_computers.csv`

#### Fichier des détections

Le fichier des détections liste le détail de chaque entrée détectée par FastFind.

Les données sont au format CSV, séparées par le caractère ',' et où les chaines de caractères sont encadrées par des guillemets '"'.

Si aucun nom n'est spécifié sur la ligne de commande, le nom par défaut du fichier est `<timestamp>-fastburn_matches.csv`

Chaque détection fait l'objet d'une ligne d'information.

Chaque ligne est constituée des champs suivants:

* `Ignore`:                      positionné à `true` si l'entrée correspond à un élément de la blacklist, à `false` sinon
* `Computer`:                    nom de la machine sur laquelle l'entrée a été détectée
* `ComputerRole`:                fonction de la machine
* `ComputerOS`:                  système d'exploitation de la machine
* `ORCVersion`:                  version de DFIR-ORC sur laquelle l'outil FastFind utilisé était basée
* `MatchType`:                   type de détection qui a déterminée la sélection de l'entrée par FastFind
* `Software`
* `Reason`:                      critères pour lequel l'entrée a été sélectionnée par FastFind
* `Filename`:                    chemin absolu du fichier
* `AltName`:                     nom alternatif du fichier
* `RegKey`:                      nom de clé de base de registre
* `RegType`:                     type de clé de base de registre
* `RegValue`:                    valeur de la clé de registre
* `FileSize`:                    taille du fichier en nombre d'octets
* `MD5`:                         condensat MD5 encodé en hexadécimal
* `SHA1`:                        condensat SHA1 encodé en hexadécimal
* `SHA256`:                      condensat SHA256 encodé en hexadécimal
* `FileCreation`:      	         date de création de fichier
* `FileLastModification`:        date de dernière modification de fichier
* `FileLastEntryChange`:         date de dernier changement des méta-informations du fichier
* `FileLastAccess`:              date de dernier accès au fichier
* `FilenameCreation`:	         date de création de l'entrée `$FN` du fichier
* `FilenameLastModification`:    date de dernière modification de l'entrée `$FN` du fichier
* `FilenameLastEntryChange`:     date de dernière modification aux méta-informations de l'entrée `$FN` du fichier
* `FilenameLastAccess`:          date de dernier accès aux méta-informations de l'entrée `$FN` du fichier
* `AltFilenameCreation`:         date de création du nom alternatif de fichier
* `AltFilenameLastModification`: date de dernière modification de la dernière modification du nom alternatif de fichier
* `AltFilenameLastEntryChange`:  date de dernière modification de la dernière modification des méta-informations du nom alternatif de fichier
* `AltFilenameLastAccess`:	      date de dernier accès au nom alternatif de fichier
* `VolumeID`:                    identifiant du volume du système de fichier dans lequel la recherche est effectuée
* `SnapshotID`:                  identifiant de l'instantané du système de fichier dans lequel la recherche est effectuée
* `ArchiveName`:                 chemin de l'archive contenant le résultat

#### Fichier de timeline

Le fichier est au format MACB couramment utilisé par les outils d'investigation forensique.

Les données sont au format CSV, séparées par le caractère ',' et où les chaines de caractères sont encadrées par des guillemets '"'.

Si aucun nom n'est spécifié sur la ligne de commande, le nom par défaut du fichier est `<timestamp>-fastburn_timeline.csv`

Ce format met toutes les informations relatives au changement d'un fichier ou une entité système sur une seule ligne.

Les entrées de ce format sont destinées à être intégrées à une *timeline* reconstituant des séquences d'action.

Une bonne description du format peut être trouvée ici <https://andreafortuna.org/2017/10/06/macb-times-in-windows-forensic-analysis/>

Pour chaque entrée les informations suivantes sont générées:

* `Timestamp`: date du changement
* `SI_MACB`: code du changement de l'entrée  au format MACB
* `FN_MACB`: code du changement du `$FILENAME` de l'entrée au format MACB
* `ComputerName`: nom de la machine sur laquelle l'entrée est identifiée
* `File`: nom du fichier ou de l'entrée système
* `ParentName`: nom du répertoire auquel l'entrée est associée
* `FullName`: chemin complet de l'entrée
* `Extension`: extension du fichier (derniers caractères après le dernier '.')
* `SizeInBytes`: taille de l'entrée en nombre d'octets
* `CreationDate`: date de création de l'entrée
* `LastModificationDate`: date de dernière modification
* `LastAccessDate`: date de dernier accès
* `LastAttrChangeDate`: date de dernier changement des attributs
* `FileNameCreationDate`: date de création de l'entrée `$FN`
* `FileNameLastModificationDate`: date de dernière modification de l'entrée `$FN`
* `FileNameLastAccessDate`: date de dernier accès à l'entrée `$FN`
* `FileNameLastAttrModificationDate`: date de dernière modification des attributs de l'entrée `$FN`
* `MD5`: condensat MD5 encodé en hexadécimal
* `SHA1`: condensat SHA1 encodé en hexadécimal
* `SHA256`: condensat SHA256 encodé en hexadécimal
* `Reason`: raison de la sélection de l'entrée par FastFind
* `ArchiveName`: nom de l'archive

##### Format de date

Toutes les dates sont exprimées au format `YYYY-MM-DD HH:MN:SS.MS`

Celles-ci sont sur le fuseau horaire UTC.

##### Format MACB

Le format MACB est une chaine de 4 caractères qui définit les changements constatés sur les méta-informations de l'entrée.

Chacun de ces caractères peut être

* `M` pour *Modified* date de dernière modification
* `A` pour *Accessed* date de dernier accès
* `C` pour *Changed* date de dernier changement de l'entrée `$MFT`
* `B` pour *Birth* date de création de l'entrée

Si le changement concerné n'est pas applicable à l'entrée la lettre est remplacée par le caractère `.`.

Exemple création de fichier:

```csv
"2020-12-23 00:09:16.944";"..C.";"..C.";"DESKTOP-LCINQGJ";"Solarwinds.Core.Businesslayer.dll";"\Users\user\Documents\TruePositive";"\Users\user\Documents\TruePositive\Solarwinds.Core.Businesslayer.dll";"dll";"1028072";"2020-12-23 00:09:32.117";"2020-12-23 16:39:01.000";"2020-12-23 00:09:38.147";"2020-12-23 00:11:44.288";"2020-12-23 00:09:32.117";"2020-12-23 16:39:01.000";"2020-12-23 00:09:32.117";"";"846E27A652A5E1BFBD0DDD38A16DC865";"D130BD75645C2433F88AC03E73395FBA172EF676";"CE77D116A074DAB7A22A0FD4F2C1AB475F16EEC42E1DED3C0B0AA8211FE858D6";"Size=1028072, SHA256=CE77D116A074DAB7A22A0FD4F2C1AB475F16EEC42E1DED3C0B0AA8211FE858D6";"share/samples/FastFind/ORC_WorkStation_DESKTOP-LCINQGJ_FastFind.7z"
```


#### Fichier de statistiques

Le fichier de statistiques récapitule les principales métriques de l'analyse.

Les données sont au format CSV, séparées par le caractère ',' et où les chaines de caractères sont encadrées par des guillemets '"'.

Si aucun nom n'est spécifié sur la ligne de commande, le nom par défaut du fichier est `<timestamp>-fastburn_stats.csv`

Le fichier de statistiques donne les nombres de résultats en fonction des critères de :

* machines
* types de machines
* système d'exploitation
* domaine Windows
* règle de détection
* nom de fichier
* condensat
* taille de fichier
* jour de création et dernière modification
* mois de création, et dernière modification

### Format du fichier de Flags

Le fichier de "whitelist"/"blacklist" est un CSV séparé par des virgules et utilisant des guillemets (`"`) comme délimiteur de chaine.

Il doit contenir les colonnes suivantes (même vides) :

* `sha256`      - condensat
* `sha1`        - condensat
* `md5`         - condensat
* `file_re`     - REGEXP à appliquer sur le chemin de fichier
* `description` - description de la détection

Une détection correspond a une ligne de "flag" si n'importe lequel des marqueurs correspond.
Le critère `file_re` est appliqué sur le champ `Fullname` du fichier de résultat.

Exemple de fichier "blacklist"
```
"sha256","sha1","md5","file_re","description"
"","DA39A3EE5E6B4B0D3255BFEF95601890AFD80709","","",""
"","","","\\windows\system32\explorer.exe",""
```

#### Exemples d'expressions régulières:

Identifier toutes les DLL du répertoire `Program Files (x86)\Adobe\Acrobat Reader DC\Reader\AcroCEF`
```
^Program Files \(x86\)\\Adobe\\Acrobat Reader DC\\Reader\\AcroCEF\\.*\.dll$
```

Identifier les exécutables du répertoire d'installation de WinRAR indépendamment de la casse
```
(?i)^\\Program Files \(x86\)\\WinRAR\\.*exe$
```

### Note sur l'usage pour traiter des grandes quantités de fichiers

La génération de traces vers le terminal a un impact non négligeable sur les performances.
Il est conseillé de rediriger la sortie d'erreur vers un fichier lorsque l'on traite des dizaines ou centaines de milliers de fichiers.

Exemple:

```sh

# Sera lent si 'Resultats' est une arborescence contenant beaucoup de fichiers de résultats.
 ./fbn Resultats

# Pour aller plus vite et générer les journaux d'exécution détaillés dans le
# fichier 'fastburn.log'
 ./fbn -debug Resultats 2> fastburn.log

```
