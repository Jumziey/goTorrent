# Notes
## Why make this application?
### First step
- One torrent per instance, download then exit (note this is a test for the protocol basically, not a real client... so no need to think about sharing etc.)
- Learn the way the bittorrents works.
- Having fun with concurrency in golang.

### Second step
Not ever worth thinking about lolz, i would consider this project done when the first step is done.

## Intressting collected thoughts about torrent clients
- Mabye build one where you can choose to upload or not? Lots of hate for that but h8 iz funziz.
- Gui:s are horrible, lets make o good cross-platform one? (haha, that was stupid... Cross platform guis sucks so bad it hurts my tummy, web kinda works but sux still).
- Mabye extend to create a neat cli tool that works at one-torrent-basis. Or mabye file basis if the protocol supports it.
- Keep configs to a minimum for gods sake!

## Trackers for testing
### Opentracker
Difficult to setup and compile and use for testing. The most popular dockerhub repo has ~5k pulls if it works out its gonna be great (altough im not sure about web socket and udp support... not relevant atm though)
https://hub.docker.com/r/ephillipe/docker-opentracker/### bittorrent-tracker
A node.js implementation (bleh, why all this javascript love, its a shitty language... just bleh.). It has lots of features such as websockets and udp .Seemes reliable and the most popular dockerhub repo has 1.2k pulls (not great, but can possibly be used)
https://hub.docker.com/r/henkel/bittorrent-tracker/

## Clients for testing
### rtorrent/rutorrent
It seemes like a good combination. Has a good docker image too
https://hub.docker.com/r/linuxserver/rutorrent/

## Encoding notes
### bytestrings number dicts and lists
<number of bytes(int)>:<string> a string
i---e	a number
d---e	a dict
l---e a list
bytestrings has the number of bytes in size, not characters (for utf8 encoding its important)

## Specification of bittorrent protocol
https://wiki.theory.org/BitTorrentSpecification is the best resource by far.
