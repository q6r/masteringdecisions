#!/bin/sh

# Install dependencies before running
#######################################
# uglify-js "npm install uglify-js -g"
# cleancc   "npm install clean-css -g"

echo "[+] Building release package"

export GIN_MODE=release
export SHELL=sh
export PROJHOME=$(pwd)
export VERSION=$(git log --pretty=oneline | cut -d ' ' -f 1 | cut -b 1-8 | head -n 1)
./deploy_regex.sh release

# Remove old release folder
if [ -d release/ ]
then
	echo "[+] Removed old release folder"
	rm -rf release/
fi

# Create structure, build, and copy
# files
echo "[+] Building and copying files"
mkdir release/ 2>&1 > /dev/null
GIN_MODE=release go build
mv gobackend release/
cd release/
cp -rf ../templates .
cp -rf ../static .
cp ../config.conf .
cp ../smtp.conf .

# Minify all js
echo "[+] Minify js files"
cd static/
for jsfiles in $(ls *.js | grep -v ".min.js")
do
	uglifyjs --compress -- $jsfiles > temp
	mv temp $jsfiles
done

# Minify all css
echo "[+] Minify css files"
cd css/
for cssfiles in $(ls *.css | grep -v ".min.js")
do
	cleancss $cssfiles -o $cssfiles
done

# package
echo "[+] Create release_$VERSION.tgz"
cd $PROJHOME
./deploy_regex.sh dev
tar czf "release_$VERSION.tgz" release
rm -rf release
