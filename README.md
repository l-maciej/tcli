# tcli
This is a simple tool used to communicate with Tedee Smart lock infrastructure.

## To run it:
* Build the executable using go build
* Add API_KEY BRIDGE_URL env variables to store the API key and Bridge IP address

## Optional
* Add the FAV - it will store Favourite lock to quickly open/lock it 

## Usage
Just use tcli and any of the available commands :
   bridge-state, bs  Check the bridge status
   lock, l           Lock operation
   callback, c       Callback operation
   help, h           Shows a list of commands or help for one command
Some of those have additionals arguments - please chceck the help within the tool using tcli *command* help
