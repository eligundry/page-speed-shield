# page-speed-shield

This is a marriage of [Google Page Speed Insights][page-speed] and [shields.io](https://shields.io/).

**Unfortunately** it doesn't work on GitHub READMEs because the image takes more than 4 seconds to render, which is the
gateway timeout threshold for camo.github.com.

## Usage

To show a badge for the Page Speed Insights score on your repo, you do the following:

```md
![Page Speed Insights](https://page-speed-shield.eligundry.com/desktop/https://eligundry.com)
![Page Speed Insights](https://page-speed-shield.eligundry.com/mobile/https://eligundry.com)
```

![Desktop Page Speed Insights](https://page-speed-shield.eligundry.com/desktop/https://eligundry.com)

Page Speed Insights takes a bit to run, so there is a first touch penalty of a couple seconds. After it's loaded once,
it's cached for a day both in the browser and on my server and will show up immediately.

[page-speed]: https://developers.google.com/speed/pagespeed/insights/
