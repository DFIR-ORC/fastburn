# Fastburnt - Build


## Makefile

Un fichier Makefile permet de gérer les compilations sur la plupart des environnements.

Si vous avez une ligne de commande ouverte dans le répertoire `fastburnt` il vous suffit de taper `make` pour générer le binaire `fastburnt_cli`.

## Compilation manuelle UNIX



Le logiciel est compilé avec la chaine de compilation Go.

```
   go build  -ldflags="" -v "dfir-orc/fastburnt/cmd/fastburnt_cli"
```

Un  exécutables est produit: `fastburnt_cli`

## Compilation manuellle Windows

### Installations des prérequis

#### Installer l'environnement de compilation Go

* Source de téléchargement <https://golang.org/>
  * Testé avec `go1.15.6.windows.amd64.msi`
  * Assurez vous d'avoir la version `Windows-AMD64`

### Installation des paquets Go

Lancer une console `Cmd.exe`.

Et y exécuter comme pour pour Unix:
```
  go get "github.com/kjk/lzmadec"
  go get "github.com/sirupsen/logrus"
  go get "github.com/andlabs/ui"
  go get "github.com/gen2brain/go-unarr"

```

### Build des exécutables

Lancer une console `Cmd.exe` et y exécuter:
```
set PATH=%PATH%;%CD%
go build -v "fastburnt/cmd/fastburnt_cli"
```

Un exécutables est produit: `fastburnt_cli.exe`


