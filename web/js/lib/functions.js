/* eslint-disable prefer-arrow-callback */
/* eslint-disable prefer-template */
/* eslint-disable no-var */
/* eslint-disable vars-on-top */
/* eslint-disable no-param-reassign */

(function mainFunctions(window) {
  /**
   * RC4
   * @param {String} key
   * @param {String} str
   * @return {String}
   */
  window.RC4 = function RC4(key, str) {
    var s = [];
    var j = 0;
    var x;
    var res = '';
    var i;
    for (i = 0; i < 256; i += 1) {
      s[i] = i;
    }
    for (i = 0; i < 256; i += 1) {
      j = (j + s[i] + key.charCodeAt(i % key.length)) % 256;
      x = s[i];
      s[i] = s[j];
      s[j] = x;
    }
    i = 0;
    j = 0;
    for (var y = 0; y < str.length; y += 1) {
      i = (i + 1) % 256;
      j = (j + s[i]) % 256;
      x = s[i];
      s[i] = s[j];
      s[j] = x;
      res += String.fromCharCode(
        // eslint-disable-next-line no-bitwise
        str.charCodeAt(y) ^ s[(s[i] + s[j]) % 256],
      );
    }
    return res;
  };

  window.decipherJS = function decipherJS(payload) {
    try {
      // eslint-disable-next-line no-eval
      eval(window.atob(payload));
      var result = window.jsChallengeRuntime();
      window.jsChallengeRuntime = undefined;
      return result;
    } catch (e) {
      // nothing
    }
    return '';
  };
})(window);
