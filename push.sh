#!/bin/bash

#cd ..
if [[ $# = 1 ]] 
then 
# cd ..
    dev=$(git status | awk 'NR=='1'{print $3}')
    echo "$dev"
    git add .
    git commit -m "$1"
    # git commit --amend --no-edit # Оставляем прежний коммит
    git push -u origin $dev
    # cat "$(date +"%Y.%m.%d_%H:%M")"
else 
    echo "Введите коммит"    
fi
