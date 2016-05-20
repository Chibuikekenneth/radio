# An attempt at creating a simultaneous streaming service, similar to internet radio


The plays music in a "music directory at the root of the app".

##How to run 
First, pull the repository and dependencies with go get.
`
    go get github.com/tonyalaribe/radio
`


go to the code directory
`
    cd /PATH/TO/GO/WORKSPACE/github/tonyalaribe/radio
`

create a music folder and copy music into the folder
`
    mkdir music
    cp ~/Music/*.mp3 music

` 

run the code and access the strean at http://localhost:8080
`
    go run main.go
`
