@echo off
IF NOT EXIST %APPDATA%"\journal" md %APPDATA%"\journal"
copy journal.exe %APPDATA%"\journal\journal.exe"
SET PATH=%PATH%;%APPDATA%"\journal\journal.exe"
echo Done
