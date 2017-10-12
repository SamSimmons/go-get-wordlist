### Abandoned

Was a project to generate a wordlist from words used on a website. The reason I stopped working on it was I realised that [cewl](https://digi.ninja/projects/cewl.php) exists and comes already installed on kali. 😂

Was a fun exercise in learning about channels in go though.

If I pick it back up, I still need to add decent regex to decide what's a valid word to include in the wordlist, and it needs to handle going more than one level deep when scraping for links


##### Setup

You need to have go installed on your system, then you just need to compile with `go build`

Then when you call the program you give it a url to hit (you can give it more than one link) eg `./go-get-wordlist https://github.com` it will dump all the strings into a file `words.txt`
