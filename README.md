# unagi

Video downloader Discord bot based on Go, using [cobalt's api](https://github.com/imputnet/cobalt)

## To-do

- [x] Use [gobalt](https://github.com/lostdusty/gobalt) instead
- [ ] File downloading failsafes
  - [x] Size requirements (max 25MB)
  - [ ] **Time outs**
- [ ] **Token stashing method**
- [ ] Cron job for deleting downloaded files
- [x] A better way to fetch the title of the `redirect`/`stream` files (note: apparently this shit works now but not for redirect i guess)
- [ ] Allow users to set cobalt video downlaoding setting parameters instead of using the default ones everytime
