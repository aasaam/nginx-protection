# Web UI

This is the UI library for challenge page.

## Customization

For custom the UI you can use these variables:

```bash
export NODE_ENV=test # for pretty HTML
export ASSETS_PATH=challenge/assets/ # for asset path of generated ui
```

For branding support you can use:

Create folder on `additional/ir_aasaam` and overwrite the locales files on `locale/en.js` to `additional/ir_aasaam/locale/en.js` also

```bash
export BRAND_ICON=ir_aasaam # for special the icon
```

## Build

Build with nodejs

```bash
nodejs build/build.js
```
