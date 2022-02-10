#!/usr/bin/env node
/* eslint-disable import/no-extraneous-dependencies */
const fs = require('fs');
const util = require('util');
const exec = util.promisify(require('child_process').exec);
const path = require('path');

const sass = require('node-sass');
const stripComments = require('strip-comments');
const uglify = require('uglify-js');

const PROJECT_DIR = path.resolve(__dirname, '..');

const minifyJsApp = async () => {
  const mainData = await fs.promises.readFile(
    `${PROJECT_DIR}/web/js/app/app.js`,
    { encoding: 'utf8' },
  );
  const compressData = uglify.minify([mainData].join('\n;'));

  await fs.promises.writeFile(
    `${PROJECT_DIR}/static/js/app.js`,
    compressData.code,
  );
};

(async () => {
  await exec(
    [
      `rm -rf ${PROJECT_DIR}/static/*`,
      `mkdir -p ${PROJECT_DIR}/static/css`,
      `mkdir -p ${PROJECT_DIR}/static/js`,
      `mkdir -p ${PROJECT_DIR}/static/fonts`,
      `cp -rf ${PROJECT_DIR}/web/node_modules/@aasaam/noto-font/dist/*-Regular.wof* ${PROJECT_DIR}/static/fonts/`,
      `cp -rf ${PROJECT_DIR}/web/node_modules/@aasaam/noto-font/dist/*-Bold.wof* ${PROJECT_DIR}/static/fonts/`,
    ].join(' && '),
  );

  const lib = [
    await fs.promises.readFile(
      `${PROJECT_DIR}/web/node_modules/angular/angular.js`,
      { encoding: 'utf8' },
    ),
    await fs.promises.readFile(
      `${PROJECT_DIR}/web/node_modules/angular-messages/angular-messages.js`,
      { encoding: 'utf8' },
    ),
  ].join(';');

  const compressLib = uglify.minify(lib);

  await fs.promises.writeFile(
    `${PROJECT_DIR}/static/js/lib.js`,
    compressLib.code,
  );

  await minifyJsApp();

  await fs.promises.writeFile(
    `${PROJECT_DIR}/static/css/main-rtl.css`,
    stripComments(
      sass
        .renderSync({
          file: `${PROJECT_DIR}/web/css/main-rtl.scss`,
          outputStyle: 'compressed',
        })
        .css.toString()
        .trim(),
      {
        preserveNewlines: false,
      },
    ),
    { encoding: 'utf8' },
  );

  await fs.promises.writeFile(
    `${PROJECT_DIR}/static/css/main-ltr.css`,
    stripComments(
      sass
        .renderSync({
          file: `${PROJECT_DIR}/web/css/main-ltr.scss`,
          outputStyle: 'compressed',
        })
        .css.toString()
        .trim(),
      {
        preserveNewlines: false,
      },
    ),
    { encoding: 'utf8' },
  );
})();
