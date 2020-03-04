# Ritchie-CLI-For-Windows

Binaries to install a ritchie cli terminal on Windows. We all know that windows terminals don't behave like linux bash. To make ritchie autocomplete works on Windows, we maked a custom terminal that works like a linux bash for Windows users. This is maked with a custom Cmder Terminal

<img src="docs/img/ritchie-cli-architecture.png">

## How it works

The installer will extract files on a windows temp and execute a bat file called rit-install.bat and this bat will make the "magic" of copy files and set windows environment variables to make rit and rit autocomplete works.

The bat file will make...

 - Necessary Folders
 - Unzip zipped custom terminal folder on %systemdrive%\tools\ (Most of time %systemdrive% is C:)
 - Copy exe file and shortcut file to necessary folders
 - Set environment variables on windows path 

When all this this is doned a custom bash of ritchie will be opened on your windows user home. You can test rit command and press <tab> to see if this work.

To open this terminal again, you can go in your windows start and find for ritchie, will be found a Ritchie bash with ritchie icon. Otherwise you can go at the binaries folder on %systemdrive%\tools\ and open ritchie.exe.

On another windows terminals like CMD or Powershell rit command will be work but not with autocomplete.

## How to maintain the installer?

To maintain the installer, edit the bat file rit-install.bat on folder rit-install-$VERSION. Generate the exe installer again with bat2exe on folder utils, increase version make a changelog file and send to repository.

If you want to edit the custom terminal, unzip the ritchie.zip of folder rit-install-$VERSION make your editions and zip the file again. Generate the exe installer again with bat2exe on folder utils, increase version make a changelog file and send to repository.

## Contributors

* [@viniciusaparecidozup](https://github.com/viniciusaparecidozup) 
* [@sandokandias](https://github.com/sandokandias) 
* [@marcosgmgm](https://github.com/marcosgmgm) 
* [@viniciusramosdefaria](https://github.com/viniciusramosdefaria) 
* [@kaduartur](https://github.com/kaduartur) 
* [@erneliojuniorzup](https://github.com/erneliojuniorzup)
