/* eslint-disable sonarjs/no-identical-functions */
/* eslint-disable prefer-destructuring */
/* eslint-disable vars-on-top */
/* eslint-disable object-shorthand */
/* eslint-disable sonarjs/no-duplicate-string */
/* eslint-disable prefer-arrow-callback */
/* eslint-disable prefer-template */
/* eslint-disable no-var */
/* eslint-disable no-param-reassign */
/* eslint prettier/prettier: ["error", { trailingComma: "none" }] */

window.angular
  .module('a', ['ngMessages'])
  .run([
    '$rootScope',
    '$timeout',
    function run($rootScope, $timeout) {
      var thresholdTime = 2;
      $rootScope.waitTime = window.wait * 1000 + thresholdTime * 1000;
      $rootScope.timeOut = window.timeout * 1000 - thresholdTime * 1000;

      $rootScope.timerProgress = {
        max: window.wait + thresholdTime,
        current: 0,
        percent: '',
        run: function runProgress() {
          if (
            $rootScope.timerProgress.current <= $rootScope.timerProgress.max
          ) {
            $timeout(function countDown() {
              $rootScope.timerProgress.current += 1;
              $rootScope.timerProgress.percent =
                Math.round(
                  ($rootScope.timerProgress.current /
                    $rootScope.timerProgress.max) *
                    100
                ).toString() + '%';
              $rootScope.timerProgress.run();
            }, 1000);
          }
        }
      };

      $rootScope.RC4 = function RC4(key, str) {
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
            str.charCodeAt(y) ^ s[(s[i] + s[j]) % 256]
          );
        }
        return res;
      };

      $rootScope.noMoreWait = false;
      $rootScope.loading = false;
      $rootScope.noMoreWaitCall = function noMoreWaitCall() {
        $rootScope.noMoreWait = true;
        document.getElementById('main').style.display = 'block';
      };
    }
  ])
  .controller('JSChallenge', [
    '$scope',
    '$timeout',
    '$http',
    '$rootScope',
    function JSChallenge($scope, $timeout, $http, $rootScope) {
      $rootScope.timerProgress.run();
      $scope.getJSChallenge = function getJSChallenge() {
        $timeout(function t1() {
          $http({
            url: window.baseURL + '/challenge/js',
            method: 'POST',
            data: {
              token: window.token
            }
          })
            .then(function r1(rs1) {
              var token1 = rs1.data.token;
              var keyContent = rs1.data.content.split(':');
              var decrypted = $rootScope.RC4(
                window.atob(keyContent[0]),
                window.atob(keyContent[1])
              );
              $http({
                url: window.baseURL + '/challenge/js',
                method: 'PATCH',
                data: {
                  token: token1,
                  content: decrypted
                }
              })
                .then(function r2() {
                  var u = new URL(window.location.href);
                  if (u.searchParams.has('url')) {
                    window.location.href = u.searchParams.get('url');
                  } else {
                    window.location.href = '/';
                  }
                })
                .catch(function e2() {
                  window.location.replace(window.location.href);
                });
            })
            .catch(function e1() {
              window.location.replace(window.location.href);
            });
        }, $rootScope.waitTime);
      };
      $scope.getJSChallenge();
    }
  ])
  .controller('CaptchaChallenge', [
    '$scope',
    '$timeout',
    '$http',
    '$rootScope',
    function CaptchaChallenge($scope, $timeout, $http, $rootScope) {
      $rootScope.timerProgress.run();
      $scope.token = '';
      $scope.errors = {};
      $timeout(function t1() {
        $scope.getTokenChallenge(function after() {
          $rootScope.noMoreWaitCall();
        });
      }, $rootScope.waitTime);
      $scope.getTokenChallenge = function getTokenChallenge(after) {
        $rootScope.loading = true;
        $http({
          url: window.baseURL + '/challenge/captcha',
          method: 'POST',
          data: {
            token: window.token
          }
        })
          .then(function r1(rs1) {
            $scope.token = rs1.data.token;
            $scope.captchaImage = rs1.data.content;
            if (after) {
              after();
            }
            $rootScope.loading = false;
          })
          .catch(function e1() {
            window.location.replace(window.location.href);
          });
      };
      $scope.submit = function submit() {
        $rootScope.loading = true;
        $scope.errors = {};
        $http({
          url: window.baseURL + '/challenge/captcha',
          method: 'PATCH',
          data: {
            Token: $scope.token,
            content: $scope.captchaValue
          }
        })
          .then(function r1() {
            var u = new URL(window.location.href);
            if (u.searchParams.has('url')) {
              window.location.href = u.searchParams.get('url');
            } else {
              window.location.href = '/';
            }
          })
          .catch(function e1() {
            $scope.errors['invalid-captcha'] = true;
            $rootScope.loading = false;
          });
      };
    }
  ])
  .controller('OTPChallenge', [
    '$scope',
    '$timeout',
    '$http',
    '$rootScope',
    function OTPChallenge($scope, $timeout, $http, $rootScope) {
      $rootScope.timerProgress.run();
      $scope.token = '';
      $scope.errors = {};
      $timeout(function t1() {
        $rootScope.noMoreWaitCall();
      }, $rootScope.waitTime);
      $scope.submit = function submit() {
        $rootScope.loading = true;
        $scope.errors = {};
        $http({
          url: window.baseURL + '/challenge/otp',
          method: 'PATCH',
          data: {
            Token: window.token,
            content: $scope.otp
          }
        })
          .then(function r1() {
            var u = new URL(window.location.href);
            if (u.searchParams.has('url')) {
              window.location.href = u.searchParams.get('url');
            } else {
              window.location.href = '/';
            }
          })
          .catch(function e1() {
            $scope.errors['invalid-otp'] = true;
            $rootScope.loading = false;
          });
      };
    }
  ]);
