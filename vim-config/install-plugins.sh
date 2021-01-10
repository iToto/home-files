#!/usr/bin/env bash

# install vim-go
if [[ ! -d "$HOME/.vim/pack/plugins/start/vim-go"  ]]; then
    echo "installing plugin vim-go"
    git clone https://github.com/fatih/vim-go.git ~/.vim/pack/plugins/start/vim-go
    vim +GoInstallBinaries
fi

# install NERDTree
if [[ ! -d "$HOME/.vim/pack/vendor/start/nerdtree"  ]]; then
    echo "installing plugin NERDTree"
    git clone https://github.com/preservim/nerdtree.git ~/.vim/pack/vendor/start/nerdtree
    vim -u NONE -c "helptags ~/.vim/pack/vendor/start/nerdtree/doc" -c q
fi

# install vim-airline
if [[ ! -d "$HOME/.vim/pack/dist/start/vim-airline"  ]]; then
    echo "installing plugin vim-airline"
    git clone https://github.com/vim-airline/vim-airline ~/.vim/pack/dist/start/vim-airline
    git clone https://github.com/vim-airline/vim-airline-themes ~/.vim/pack/dist/start/vim-airline-themes-
fi


# install git gutter
if [[ ! -d "$HOME/.vim/pack/airblade/start"  ]]; then
    echo "installing plugin airblade-vim-gutter"
    mkdir -p ~/.vim/pack/airblade/start
    cd ~/.vim/pack/airblade/start
    git clone https://github.com/airblade/vim-gitgutter.git
    vim -u NONE -c "helptags vim-gitgutter/doc" -c q
fi
