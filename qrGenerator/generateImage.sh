#!/bin/bash

# usage: ./generateImage.sh <url> <Title> <Remove Date> <Output file>

temp_file=$(mktemp -q)
qrencode -s 30 -d 30 -o "$temp_file" "$1"

magick \
 \( -background none -size 874x100 canvas:none \) \
  \( -background none -fill black -font "Montserrat-Bold" -gravity center -pointsize 60 -size 874x caption:"$2" \) \
  \( -background none -size 1x40 canvas:none \) \
  \( "$temp_file" -resize 600x600 \) \
  \( -background none -fill black -font "Montserrat-Bold" -gravity center -pointsize 36 -size 874x caption:"$1" \) \
  \( -background none -size 1x200 canvas:none \) \
  \( -background none -fill black -font "Montserrat-Bold" -gravity center -pointsize 28 -size 800x caption:"Dieser Sticker ist Teil einer Schnitzeljagd und wird am $3 wieder entfernt." \) \
  -background white -append \
  -gravity north -extent 874x1240 \
  "$4"

rm "$temp_file"