/* eslint-disable prefer-destructuring */
/* eslint-disable vars-on-top */
/* eslint-disable object-shorthand */
/* eslint-disable sonarjs/no-duplicate-string */
/* eslint-disable prefer-arrow-callback */
/* eslint-disable prefer-template */
/* eslint-disable no-var */
/* eslint-disable no-param-reassign */
/* eslint prettier/prettier: ["error", { trailingComma: "none" }] */

if (!('URL' in window) || !('URLSearchParams' in window)) {
  document.write(
    '<script src="' +
      window.assets +
      'js/_url-polyfill.' +
      window.version +
      '.min.js' +
      '"></script>'
  );
}
