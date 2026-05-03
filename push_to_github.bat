@echo off
git init
git remote add origin https://github.com/ZINTENT/zintent.git
git add .
git commit -m "Initial Release v2.1.0"
git branch -M main
git push -u origin main
pause
