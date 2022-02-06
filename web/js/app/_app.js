/* eslint-disable prefer-arrow-callback */
/* eslint-disable prefer-template */
/* eslint-disable no-var */
/* eslint-disable vars-on-top */
/* eslint-disable no-param-reassign */

(function mainFunctions(window) {
  window.errorLog = function errorLog(e) {
    // eslint-disable-next-line no-console
    console.error(e);
  };
  window.app = window.angular.module('p', ['ngMessages']).run([
    '$rootScope',
    '$timeout',
    function run($rootScope, $timeout) {
      $rootScope.config = window.config;

      $rootScope.ClientJSInfo = function ClientJSInfo() {
        return {
          devicePixelRatio: window.devicePixelRatio,
          screenPixelDepth: window.screen.pixelDepth,
          screenColorDepth: window.screen.colorDepth,
          screenWidth: window.screen.width,
          screenHeight: window.screen.height,
          userAgent: navigator.userAgent,
        };
      };

      $rootScope.getCountryName = function getCountryName(code) {
        if ('Intl' in window && 'DisplayNames' in Intl) {
          try {
            const title = new Intl.DisplayNames([window.config.lang], {
              type: 'region',
            });
            return title.of(code);
          } catch (e) {
            // nothing
          }
        }
        return code;
      };

      $rootScope.failed = false;

      $rootScope.reload = function reload(timeoutSeconds) {
        $rootScope.failed = true;
        $timeout(function reloadTimeout() {
          // eslint-disable-next-line no-self-assign
          window.location = window.location;
        }, timeoutSeconds);
      };

      $rootScope.getCountryFlag = function getCountryFlag(code) {
        if (code) {
          try {
            return code.toUpperCase().replace(/./g, function replace(char) {
              return String.fromCodePoint(char.charCodeAt(0) + 127397);
            });
          } catch (e) {
            // nothing
          }
        }
        return '';
      };

      $rootScope.mailTo = function mailTo() {
        var host = window.location.hostname;
        return [
          'mailto:',
          window.config.supportInfo.email,
          '?subject=',
          encodeURIComponent(
            'report protection ' +
              $rootScope.config.challengeType +
              ' on ' +
              host,
          ),
          '&body=',
          encodeURIComponent(
            [
              'node: ' + window.config.ipData.nodeID,
              'ip: ' + window.config.ipData.ip,
              'host: ' + host,
              'country: ' + window.config.ipData.country,
              'asn: ' + window.config.ipData.asn,
              'asn name: ' + window.config.ipData.asn_org,
            ].join('\r\n'),
          ),
        ].join('');
      };

      $rootScope.urlSupport = function urlSupport() {
        var host = window.location.hostname;
        return (
          window.config.supportInfo.url +
          '?' +
          [
            'mode=protection',
            'node=' + encodeURIComponent(window.config.ipData.nodeID),
            'type=' + encodeURIComponent($rootScope.config.challengeType),
            'host=' + encodeURIComponent(host),
            'ip=' + encodeURIComponent(window.config.ipData.ip),
            'country=' + encodeURIComponent(window.config.ipData.country),
            'asn=' + encodeURIComponent(window.config.ipData.asn),
            'asn_org=' + encodeURIComponent(window.config.ipData.asn_org),
          ].join('&')
        );
      };

      // show page
      $timeout(function timeout() {
        document.querySelector('#app').style.display = 'block';
      }, 200);
    },
  ]);
})(window);
