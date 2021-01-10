:#!/bin/sh

# install vim-go
git clone https://github.com/fatih/vim-go.git ~/.vim/pack/plugins/start/vim-go
vim +GoInstallBinaries

# install NERDTree
git clone https://github.com/preservim/nerdtree.git ~/.vim/pack/vendor/start/nerdtree
vim -u NONE -c "helptags ~/.vim/pack/vendor/start/nerdtree/doc" -c q

# install vim-airline
git clone https://github.com/vim-airline/vim-airline ~/.vim/pack/dist/start/vim-airline
git clone https://github.com/vim-airline/vim-airline-themes ~/.vim/pack/dist/start/vim-airline-themes-

