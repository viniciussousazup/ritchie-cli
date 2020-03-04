@echo off
title INSTALL RITCHIE
mode 100,30
echo.
echo.
echo ====================================================================================================
echo                                            Zup Innovations
echo                                         Ritchie is installing!
echo                                            Please wait..
echo ====================================================================================================
echo.
echo.
	md %USERPROFILE%\.rit 2>nul
	md %systemdrive%\tools\ 2>nul
	md %systemdrive%\ProgramData\Ritchie\ 2>nul
	tar -xf ritchie.zip -C %systemdrive%\tools\
	xcopy "rit.exe" "%USERPROFILE%\.rit" /K /D /H /Y 2>nul
	xcopy "Ritchie Bash.lnk" "%systemdrive%\ProgramData\Ritchie\" /K /D /H /Y 2>nul
	echo %PATH% | find /C /I "%USERPROFILE%\.rit" 2>nul || setx /m PATH "%PATH%;%USERPROFILE%\.rit
	echo %PATH% | find /C /I "%systemdrive%\tools\ritchie\cygwin\bin" 2>nul || setx /m PATH "%PATH%;%systemdrive%\tools\ritchie\cygwin\bin
	echo %PATH% | find /C /I "%USERPROFILE%\.rit" 2>nul || setx /m PATH "%PATH%;%USERPROFILE%\.rit
	cls
echo.
echo.
echo ====================================================================================================
echo                                            Zup Innovations
echo                                         Ritchie is ready to use!
echo                                               Enjoy!
echo ====================================================================================================
echo.
echo.
	%systemdrive%\tools\ritchie\ritchie.exe
	

