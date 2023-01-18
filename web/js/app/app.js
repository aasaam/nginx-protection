/* eslint-disable sonarjs/no-identical-functions */
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
    // eslint-disable-next-line sonarjs/no-duplicate-string
    '$rootScope',
    '$timeout',
    '$http',
    // eslint-disable-next-line sonarjs/cognitive-complexity
    function run($rootScope, $timeout, $http) {
      $rootScope.config = window.config;

      $rootScope.getTimeAccuracy = function getTimeAccuracy() {
        if ('Intl' in window && 'RelativeTimeFormat' in Intl) {
          // eslint-disable-next-line prefer-destructuring
          var timeAccuracy = window.config.timeAccuracy;
          var fmt = new Intl.RelativeTimeFormat(window.config.lang, {
            style: 'short',
          });
          if (timeAccuracy < 5) {
            return '';
          }
          if (timeAccuracy <= 120) {
            return fmt.format(timeAccuracy, 'seconds');
          }
          if (timeAccuracy <= 3600) {
            return fmt.format(timeAccuracy, 'minute');
          }
          if (timeAccuracy <= 86400) {
            return fmt.format(timeAccuracy, 'hour');
          }
          return fmt.format(timeAccuracy, 'day');
        }
        return window.config.timeAccuracy;
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

      $rootScope.reload = function reload(timeoutSeconds) {
        $timeout(function reloadTimeout() {
          document.querySelector('#app').style.display = 'none';
          // eslint-disable-next-line no-self-assign
          window.location = window.location;
        }, timeoutSeconds);
      };

      $rootScope.justReload = function justReload() {
        document.querySelector('#app').style.display = 'none';
        // eslint-disable-next-line no-self-assign
        window.location = window.location;
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

      // formData
      $rootScope.formData = {
        failed: false,
        failedReload: false,
        jsValue: '',
        captchaValue: '',
        totpValue: '',
        username: '',
        password: '',
      };

      $rootScope.challengePATCH = function challengePATCH(
        callbackFunction,
        errorCallback,
      ) {
        $http({
          url: window.config.baseURL + '/challenge',
          method: 'PATCH',
          data: {
            t: window.config.challengeToken,
          },
        })
          .then(function onResponsePath(response) {
            // eslint-disable-next-line prefer-destructuring
            var data = response.data;
            window.config.challengeToken = data.t;
            // eslint-disable-next-line no-eval
            eval(window.atob(data.js));
            // eslint-disable-next-line vars-on-top
            var jsSolved = window.jsv();
            delete window.jsv;
            window.jsv = undefined;
            $rootScope.formData.jsValue = jsSolved;
            if (data.captcha && data.captcha.length >= 1) {
              $rootScope.captchaImage = data.captcha;
            }
            callbackFunction();
          })
          .catch(function onErrorPath(e) {
            if (errorCallback) {
              errorCallback(e);
            } else {
              window.location.reload();
            }
          });
      };

      // challengePOST
      $rootScope.challengePOSTCalled = false;
      $rootScope.challengePOST = function challengePOST(
        dataObject,
        errorCallback,
      ) {
        if ($rootScope.challengePOSTCalled) {
          return;
        }
        $rootScope.challengePOSTCalled = true;
        $http({
          url: window.config.baseURL + '/challenge',
          method: 'POST',
          data: dataObject,
        })
          .then(function onResponsePost() {
            window.location.href = $rootScope.config.protectedPath;
          })
          .catch(function onErrorPost() {
            if (errorCallback) {
              errorCallback();
            } else {
              $rootScope.justReload();
            }
          });
      };

      $rootScope.times = {
        timeAccuracy: Math.abs(window.config.timeAccuracy) < 5,
        progress: 0,
        waitSeconds: parseInt($rootScope.config.waitSeconds, 10),
        timeoutTotal: parseInt($rootScope.config.timeoutSeconds, 10),
        start: Math.round(Date.now() / 1000),
        timeoutProgress: parseInt($rootScope.config.waitSeconds, 10),
      };

      $rootScope.times.end =
        $rootScope.times.start + $rootScope.times.waitSeconds + 1;
      $rootScope.times.total = $rootScope.times.end - $rootScope.times.start;

      $rootScope.captchaImage =
        'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7';

      $rootScope.timeoutProgressFunc = function timeoutProgressFunc() {
        $timeout(function timeoutProgressFuncTimeout() {
          if (
            $rootScope.times.timeoutProgress < $rootScope.times.timeoutTotal
          ) {
            $rootScope.times.timeoutProgress += 1;
            $rootScope.timeoutProgressFunc();
          } else {
            $rootScope.times.timeoutProgress = $rootScope.times.timeoutSeconds;
            $rootScope.formData.failed = true;
            $rootScope.formData.failedReload = true;
          }
        }, 1000);
      };

      $rootScope.progressUpdate = function progressUpdate(callbackFunction) {
        $timeout(function progress() {
          var state = Math.round(Date.now() / 1000) - $rootScope.times.start;
          var currentProgress = Math.round(
            (state / $rootScope.times.total) * 100,
          );
          if (currentProgress < 100) {
            $rootScope.times.progress = Math.floor(currentProgress);
            $rootScope.progressUpdate(callbackFunction);
          } else {
            $rootScope.times.progress = 100;
            callbackFunction();
          }
        }, 500);
      };

      // show page
      $timeout(function timeout() {
        document.querySelector('#app').style.display = 'block';
      }, 200);
    },
  ]);

  // BlockChallenge
  window.app.controller('block', [function BlockChallenge() {}]);

  // JSChallenge
  window.app.controller('js', [
    '$rootScope',
    function JSChallenge($rootScope) {
      $rootScope.progressUpdate(function updateChallenge() {
        $rootScope.challengePATCH(function challengePATCHApplied() {
          $rootScope.challengePOST({
            t: window.config.challengeToken,
            j: $rootScope.formData.jsValue,
          });
        });
      });
    },
  ]);

  // CaptchaChallenge
  window.app.controller('captcha', [
    '$scope',
    '$rootScope',
    function CaptchaChallenge($scope, $rootScope) {
      $rootScope.progressUpdate(function updateChallenge() {
        $rootScope.challengePATCH(function challengePATCHApplied() {
          document.getElementById('inputCaptcha').focus();
        });
      });

      $rootScope.timeoutProgressFunc();

      $scope.submit = function submit() {
        $rootScope.challengePOST(
          {
            t: window.config.challengeToken,
            j: $rootScope.formData.jsValue,
            c: $rootScope.formData.captchaValue,
          },
          function onError() {
            $rootScope.formData.failed = true;
            $rootScope.formData.failedReload = false;
            $rootScope.reload(3000);
          },
        );
      };
    },
  ]);

  // TOTPChallenge
  window.app.controller('totp', [
    '$scope',
    '$rootScope',
    function TOTPChallenge($scope, $rootScope) {
      $rootScope.progressUpdate(function updateChallenge() {
        $rootScope.challengePATCH(function challengePATCHApplied() {
          document.getElementById('inputTOTP').focus();
        });
      });

      $rootScope.timeoutProgressFunc();

      $scope.submit = function submit() {
        $rootScope.challengePOST(
          {
            t: window.config.challengeToken,
            j: $rootScope.formData.jsValue,
            totp: $rootScope.formData.totpValue.toString(),
          },
          function onError() {
            $rootScope.formData.failed = true;
            $rootScope.formData.failedReload = false;
            $rootScope.reload(3000);
          },
        );
      };
    },
  ]);

  // LDAPChallenge
  window.app.controller('ldap', [
    '$scope',
    '$rootScope',
    function LDAPChallenge($scope, $rootScope) {
      $rootScope.progressUpdate(function updateChallenge() {
        $rootScope.challengePATCH(function challengePATCHApplied() {
          document.getElementById('inputUsername').focus();
        });
      });

      $rootScope.timeoutProgressFunc();

      $scope.submit = function submit() {
        $rootScope.challengePOST(
          {
            t: window.config.challengeToken,
            j: $rootScope.formData.jsValue,
            c: $rootScope.formData.captchaValue,
            lu: $rootScope.formData.username,
            lp: $rootScope.formData.password,
          },
          function onError() {
            $rootScope.formData.failed = true;
            $rootScope.formData.failedReload = false;
            $rootScope.reload(3000);
          },
        );
      };
    },
  ]);
})(window);
