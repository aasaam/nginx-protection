/* eslint-disable import/no-extraneous-dependencies */
// @ts-check

const util = require('util');
const fs = require('fs');
const { resolve, parse } = require('path');
const exec = util.promisify(require('child_process').exec);

const { Organization } = require('@aasaam/information');

const flat = require('flat');
const { merge, uniq } = require('lodash');
const { to } = require('await-to-js');
const pug = require('pug');

const projectDir = resolve(__dirname, '..');
const rootDir = resolve(projectDir, '..');
const assetDir = resolve(rootDir, 'assets');
const { log } = console;
const fsp = fs.promises;

const aasaamLogo = fs.readFileSync(
  `${projectDir}/node_modules/@aasaam/information/logo/aasaam-mono.svg`,
  { encoding: 'utf-8' }
);

const BRAND_ICON = process.env.BRAND_ICON ? process.env.BRAND_ICON : '_';

const rtlLanguages = ['ar', 'dv', 'fa', 'he', 'ps', 'ur', 'yi'];

let locales = {};
const localeDir = `${projectDir}/locale`;
const localeAddonDir = `${projectDir}/additional/${BRAND_ICON}/locale`;
const sizes = [];
fs.readdirSync(localeDir).forEach((file) => {
  const { name: lang } = parse(file);
  const p = resolve(localeDir, file);
  const addonPath = resolve(localeAddonDir, file);

  // eslint-disable-next-line global-require, import/no-dynamic-require
  locales[lang] = require(p);

  if (fs.existsSync(addonPath)) {
    // eslint-disable-next-line global-require, import/no-dynamic-require
    const addon = require(addonPath);
    locales[lang] = merge(locales[lang], addon);
  }

  if (Organization[lang]) {
    locales[lang].aasaam = {
      name: Organization[lang].name,
      description: Organization[lang].description,
      url: Organization[lang].url
    };
  } else {
    locales[lang].aasaam = {
      name: Organization.en.name,
      description: Organization.en.description,
      url: Organization.en.url
    };
  }

  locales[lang].dir = rtlLanguages.includes(lang) ? 'rtl' : 'ltr';
  locales[lang].lang = lang;

  sizes.push(Object.keys(flat(locales[lang])).length);
});

// @ts-ignore
if (uniq(sizes).length !== 1) {
  log('Seems translation fields are not equal');
  process.exit();
}

const languageData = {};
Object.keys(locales).forEach((lang) => {
  // @ts-ignore
  const regionNamesInNative = new Intl.DisplayNames([lang], {
    type: 'language'
  });

  if (!languageData[lang]) {
    languageData[lang] = {};
  }
  languageData[lang][lang] = regionNamesInNative.of(lang);
  Object.keys(locales).forEach((l2) => {
    // @ts-ignore
    const regionNamesInOther = new Intl.DisplayNames([l2], {
      type: 'language'
    });

    languageData[lang][l2] = regionNamesInOther.of(lang);
  });
});

if (process.env.LANGUAGES) {
  const validLanguages = process.env.LANGUAGES.split(',');
  const validLocales = {};
  Object.keys(locales).forEach((lang) => {
    if (validLanguages.includes(lang)) {
      validLocales[lang] = locales[lang];
    }
  });
  locales = validLocales;
}

const templates = [];
const templateDir = `${projectDir}/templates`;
fs.readdirSync(templateDir).forEach((file) => {
  const m = file.match(/([a-z0-9-]+)\.pug$/i);
  if (!m) {
    return;
  }
  templates.push({
    appName: m[1],
    file
  });
});

(async () => {
  const version = Math.random().toString(36).substring(2);

  let commands = [
    `mkdir -p ${projectDir}/tmp`,
    `mkdir -p ${projectDir}/public/js`,
    `rm -rf ${projectDir}/public/js/_*`,
    `rm -rf ${projectDir}/public/images`,
    `rm -rf ${projectDir}/public/fonts`,
    `mkdir -p ${projectDir}/public/fonts`,
    `mkdir -p ${projectDir}/public/images`,
    `rm -rf ${projectDir}/public/*.html`
  ];

  await exec(commands.join(' && '));

  const logoPath = `${projectDir}/additional/logo.svg`;
  const [, logoExist] = await to(fsp.stat(logoPath));
  if (process.env.BRAND_ICON) {
    commands.push(
      `cp -rf ${projectDir}/node_modules/@aasaam/brand-icons/svg/${process.env.BRAND_ICON}.svg ${projectDir}/public/images/logo.${version}.svg`
    );
  } else if (logoExist) {
    commands.push(
      `cp -rf ${projectDir}/additional/logo.svg ${projectDir}/public/images/logo.${version}.svg`
    );
  } else {
    commands.push(
      `cp -rf ${projectDir}/node_modules/@aasaam/brand-icons/svg/ir_aasaam.svg ${projectDir}/public/images/logo.${version}.svg`
    );
  }

  await exec(commands.join(' && '));

  const organizationLogo = fs.readFileSync(
    `${projectDir}/public/images/logo.${version}.svg`,
    {
      encoding: 'utf-8'
    }
  );

  commands = [
    // font
    `cp -rf ${projectDir}/node_modules/@aasaam/noto-font/dist/*Regular*.wof* ${projectDir}/public/fonts/`,
    `cp -rf ${projectDir}/node_modules/@aasaam/noto-font/dist/*Bold*.wof* ${projectDir}/public/fonts/`,
    `cd ${projectDir}/public/fonts/`,
    `for f in *.woff; do mv "$f" "\${f%.woff}.${version}.woff"; done`,
    `for f in *.woff2; do mv "$f" "\${f%.woff2}.${version}.woff2"; done`,
    // url-polyfill
    `cat ${projectDir}/node_modules/url-polyfill/url-polyfill.js > ${projectDir}/public/js/_url-polyfill.js`,
    `cd ${projectDir}/public/js`,
    `${projectDir}/node_modules/.bin/uglifyjs --compress --mangle --output _url-polyfill.${version}.min.js _url-polyfill.js`,
    // framework
    `cat ${projectDir}/node_modules/angular/angular.js > ${projectDir}/public/js/_framework.js`,
    `cat ${projectDir}/node_modules/angular-messages/angular-messages.js >> ${projectDir}/public/js/_framework.js`,
    `cd ${projectDir}/public/js`,
    `${projectDir}/node_modules/.bin/uglifyjs --compress --mangle --output _framework.${version}.min.js _framework.js`,
    // polyfill
    `cd ${projectDir}/public/js`,
    `${projectDir}/node_modules/.bin/uglifyjs --compress --mangle --output ${projectDir}/tmp/polyfill.min.js polyfill.js`,
    // init
    `cd ${projectDir}/public/js`,
    `${projectDir}/node_modules/.bin/uglifyjs --compress --mangle --output _app.${version}.min.js app.js`,
    // version
    `echo "\\$version: \\"${version}\\";" > ${projectDir}/public/css/__version.scss`
  ];

  await exec(commands.join(' && '));

  await exec(
    [
      `rm -rf ${projectDir}/public/css/*.css`,
      `rm -rf ${projectDir}/public/css/*.map`,
      `cd ${projectDir}/public/css`,
      `${projectDir}/node_modules/.bin/node-sass main-rtl.scss main-rtl.${version}.css --output-style compressed`,
      `${projectDir}/node_modules/.bin/node-sass main-ltr.scss main-ltr.${version}.css --output-style compressed`
    ].join(' && ')
  );

  const polyfillScript = await fsp.readFile(
    `${projectDir}/tmp/polyfill.min.js`,
    {
      encoding: 'utf-8'
    }
  );

  let promises = [];
  Object.keys(locales).forEach((lang) => {
    templates.forEach(({ file }) => {
      promises.push(
        new Promise((res, rej) => {
          pug.renderFile(
            resolve(projectDir, 'templates', file),
            {
              version,
              polyfillScript: polyfillScript.trim(),
              aasaamLogo: aasaamLogo.trim(),
              organizationLogo: organizationLogo.trim(),
              i18n: locales[lang],
              assetPath: process.env.ASSETS_PATH ? process.env.ASSETS_PATH : '',
              versionJson: JSON.stringify(version),
              languageData: JSON.stringify(languageData),
              pretty: process.env.NODE_ENV === 'test'
            },
            (e, html) => {
              if (e) {
                rej(e);
              } else {
                const fileName = file.replace('.pug', `.${lang}.html`);
                const path = resolve(projectDir, 'public', fileName);

                res({
                  path,
                  fileName,
                  file,
                  html
                });
              }
            }
          );
        })
      );
    });
  });

  const htmlList = await Promise.all(promises);

  promises = [];
  htmlList.forEach(({ path, html }) => {
    promises.push(
      new Promise((res, rej) => {
        fsp.writeFile(path, html).then(res).catch(rej);
      })
    );
  });

  await Promise.all(promises);

  commands = [
    `rm -rf ${assetDir}`,
    `mkdir -p ${assetDir}/public/js`,
    `mkdir -p ${assetDir}/public/fonts`,
    `mkdir -p ${assetDir}/public/css`,
    `mkdir -p ${assetDir}/public/images`,
    `mkdir -p ${assetDir}/templates`,
    `cp -rf ${projectDir}/public/*.html ${assetDir}/templates/`,
    `cp -rf ${projectDir}/public/js/_*.min.js ${assetDir}/public/js/`,
    `cp -rf ${projectDir}/public/css/*.css ${assetDir}/public/css/`,
    `cp -rf ${projectDir}/public/fonts/* ${assetDir}/public/fonts/`,
    `cp -rf ${projectDir}/public/images/* ${assetDir}/public/images/`
  ];

  await exec(commands.join(' && '));

  log('UI GENERATED SUCCESSFULLY');
})();
