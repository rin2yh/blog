# Theme images

Drop the following PNGs here before submitting to https://themes.gohugo.io/:

| File | Size | Purpose |
|---|---|---|
| `screenshot.png` | 1500 × 1000 px (3:2) | Header image on the theme listing page |
| `tn.png` | 900 × 600 px (3:2) | Thumbnail in the theme grid |

Easiest capture flow:

```bash
cd exampleSite
hugo serve --themesDir ../..
# open http://localhost:1313 in a 1500-wide browser window
# capture full-viewport screenshots, downscale tn.png to 900x600
```

Both files are referenced from `README.md` and the Hugo Themes submission process — they need to exist before publishing.
