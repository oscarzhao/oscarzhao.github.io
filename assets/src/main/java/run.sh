#! /bin/bash

find . -name "*.java" > sources.txt
javac @sources.txt
java com.oscarzhao.Main

# clean up
find . -name "*.class" | xargs rm -rf 
find . -name "sources.txt" | xargs rm -rf

