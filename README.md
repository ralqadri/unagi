# unagi

video downloader discord bot based on go, using [cobalt's api](https://github.com/imputnet/cobalt)

this was made just to save memes on my personal discord channel/my phone

## todo

- [ ] file downloading failsafes
  - [x] file size requirements (max 25MB)
  - [ ] time outs
- [ ] maybe switch to a cron job for deleting downloaded files since deleting files can fail on download.go rn
- [ ] a better way to fetch the title of the stream (and potentially for stuff that are also "success" and "redirect" responses)
