# Enable Powerlevel10k instant prompt. Should stay close to the top of ~/.zshrc.
# Initialization code that may require console input (password prompts, [y/n]
# confirmations, etc.) must go above this block; everything else may go below.
if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
  source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi

# Path to your oh-my-zsh installation.
export ZSH=$HOME/.oh-my-zsh

### My Local Aliases
source ~/.aliases

# Set name of the theme to load.
# Look in ~/.oh-my-zsh/themes/
# Optionally, if you set this to "random", it'll load a random theme each
# time that oh-my-zsh is loaded.
# ZSH_THEME="re5et"
# ZSH_THEME="ys"
# ZSH_THEME="xiong-chiamiov"
# ZSH_THEME="muse"
# ZSH_THEME="rkj-repos"
ZSH_THEME="powerlevel10k/powerlevel10k"

# Example aliases
# alias zshconfig="mate ~/.zshrc"
# alias ohmyzsh="mate ~/.oh-my-zsh"

# Uncomment the following line to use case-sensitive completion.
# CASE_SENSITIVE="true"

# Uncomment the following line to disable bi-weekly auto-update checks.
# DISABLE_AUTO_UPDATE="true"

# Uncomment the following line to change how often to auto-update (in days).
# export UPDATE_ZSH_DAYS=13

# Uncomment the following line to disable colors in ls.
# DISABLE_LS_COLORS="true"

# Uncomment the following line to disable auto-setting terminal title.
# DISABLE_AUTO_TITLE="true"

# Uncomment the following line to disable command auto-correction.
# DISABLE_CORRECTION="true"

# Uncomment the following line to display red dots whilst waiting for completion.
# COMPLETION_WAITING_DOTS="true"

# Uncomment the following line if you want to disable marking untracked files
# under VCS as dirty. This makes repository status check for large repositories
# much, much faster.
# DISABLE_UNTRACKED_FILES_DIRTY="true"

# Uncomment the following line if you want to change the command execution time
# stamp shown in the history command output.
# The optional three formats: "mm/dd/yyyy"|"dd.mm.yyyy"|"yyyy-mm-dd"
# HIST_STAMPS="mm/dd/yyyy"

# Would you like to use another custom folder than $ZSH/custom?
# ZSH_CUSTOM=/path/to/new-custom-folder

# Which plugins would you like to load? (plugins can be found in ~/.oh-my-zsh/plugins/*)
# Custom plugins may be added to ~/.oh-my-zsh/custom/plugins/
# Example format: plugins=(rails git textmate ruby lighthouse)
plugins=(git battery composer fabric git-extras git-flow git-prompt gitignore golang golang heroku jira laravel4 npm nvm sublime vagrant docker brew)

source $ZSH/oh-my-zsh.sh

# User configuration

export PATH="/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
export PATH=$PATH:~/.local/bin

if [ -d "/usr/local/opt/ruby/bin" ]; then
  export PATH=/usr/local/opt/ruby/bin:$PATH
  export PATH=`gem environment gemdir`/bin:$PATH
fi

# export MANPATH="/usr/local/man:$MANPATH"

# You may need to manually set your language environment
# export LANG=en_US.UTF-8

# Preferred editor for local and remote sessions
# if [[ -n $SSH_CONNECTION ]]; then
#   export EDITOR='vim'
# else
#   export EDITOR='mvim'
# fi

# Compilation flags
# export ARCHFLAGS="-arch x86_64"

# ssh
# export SSH_KEY_PATH="~/.ssh/dsa_id"

#My Aliases
# Mac Helpers
alias show_hidden="defaults write com.apple.Finder AppleShowAllFiles YES && killall Finder"
alias hide_hidden="defaults write com.apple.Finder AppleShowAllFiles NO && killall Finder"

# Helpers
alias grep='grep --color=auto' # Always highlight grep search term
alias ping='ping -c 5'      # Pings with 5 packets, not unlimited
alias df='df -h'            # Disk free, in gigabytes, not bytes
alias du='du -h -c'         # Calculate total disk usage for a folder
alias sgi='sudo gem install' # Install ruby stuff
alias clj='clj-env-dir'        # Clojure helper
alias clr='clear;echo "Currently logged in on $(tty), as $(whoami) in directory $(pwd)."'
alias tt='tt++ $HOME/.ttconf'
alias svim="sudo vim" # Run vim as super user
alias emc="emacsclient -n" # no blocking terminal waiting for edit

# GIT ALIASES
alias g=git
alias ga='git add'
alias gb='git branch'
alias gba='git branch -a'
alias gc='git commit -v'
alias gl='git pull'
alias gp='git push'
alias gst='git status -sb'
alias gsd='git svn dcommit'
alias gsr='git svn rebase'
alias gs='git stash'
alias gsa='git stash apply'
alias gr='git stash && git svn rebase && git svn dcommit && git stash pop' # git refresh
alias gd='git diff | $GIT_EDITOR -'
alias gmv='git mv'
alias gho='$(git remote -v 2> /dev/null | grep github | sed -e "s/.*git\:\/\/\([a-z]\.\)*/\1/" -e "s/\.git$//g" -e "s/.*@\(.*\)$/\1/g" | tr ":" "/" | tr -d "\011" | sed -e "s/^/open http:\/\//g")'

# Functions
extract () {
    if [ -f $1 ] ; then
        case $1 in
            *.tar.bz2)        tar xjf $1        ;;
            *.tar.gz)         tar xzf $1        ;;
            *.bz2)            bunzip2 $1        ;;
            *.rar)            unrar x $1        ;;
            *.gz)             gunzip $1         ;;
            *.tar)            tar xf $1         ;;
            *.tbz2)           tar xjf $1        ;;
            *.tgz)            tar xzf $1        ;;
            *.zip)            unzip $1          ;;
            *.Z)              uncompress $1     ;;
            *)                echo "'$1' cannot be extracted via extract()" ;;
        esac
    else
        echo "'$1' is not a valid file"
    fi
}

dls () {
 # directory LS
 echo `ls -l | grep "^d" | awk '{ print $9 }' | tr -d "/"`
}

killit() {
    # Kills any process that matches a regexp passed to it
    ps aux | grep -v "grep" | grep "$@" | awk '{print $2}' | xargs sudo kill
}

tree () {
    find $@ -print | sed -e 's;[^/]*/;|____;g;s;____|; |;g'
}

exip () {
    # gather external ip address
    echo -n "Current External IP: "
    curl -s -m 5 http://myip.dk | grep "ha4" | sed -e 's/.*ha4">//g' -e 's/<\/span>.*//g'
}

ips () {
    # determine local IP address
    ifconfig | grep "inet " | awk '{ print $2 }'
}

if [ -z "\${which tree}" ]; then
  tree () {
      find $@ -print | sed -e 's;[^/]*/;|____;g;s;____|; |;g'
  }
fi

mcd () {
    mkdir "$@" && cd "$@"
}


### This loads nvm
[ -s "/Users/sdagostino/.nvm/nvm.sh" ] && . "/Users/sdagostino/.nvm/nvm.sh"

### Homebrew
# export PATH="/usr/local/sbin:$PATH"
export PATH=$PATH:/opt/homebrew/bin
# eval "$(/usr/local/bin/brew shellenv)"

### PHP
# export PATH="$(brew --prefix homebrew/php/php56)/bin:$PATH"

### Added by the Heroku Toolbelt
export PATH="/usr/local/heroku/bin:$PATH"

export PATH="$PATH:$HOME/.rvm/bin" # Add RVM to PATH for scripting

#Postgres.app
export PATH=$PATH:/Applications/Postgres.app/Contents/Versions/latest/bin

# Add Visual Studio Code (code)
export PATH=$PATH:/Applications/Visual\ Studio\ Code.app/Contents/Resources/app/bin

#Go Path
export GOPATH=~/Dev/golang
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN

# export PATH=$PATH:/usr/local/opt/go/libexec/bin

test -e "${HOME}/.iterm2_shell_integration.zsh" && source "${HOME}/.iterm2_shell_integration.zsh"

# Docker
dockerEnter() { docker exec -it $1 /bin/bash ;};
alias de=dockerEnter

### My Custom Aliases
source ~/.aliases-custom

### My Local Functions
source ~/.functions

### Apple shortcuts
alias fixOSXCamera='sudo killall AppleCameraAssistant'

### My Custom Exports
source ~/.exports
autoload -U +X bashcompinit && bashcompinit
complete -o nospace -C /usr/local/bin/terraform terraform
export PATH="/usr/local/opt/ruby/bin:$PATH"
TERM='xterm'
export TERM

# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh.
[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh
