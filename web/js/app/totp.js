/* eslint-disable prefer-template */
/* eslint-disable prefer-arrow-callback */
/* eslint-disable no-param-reassign */
/* eslint-disable prefer-arrow-callback */
/* eslint-disable prefer-destructuring */
/* eslint-disable no-var */
(function application(window) {
  window.app.controller('totp', [
    '$scope',
    '$rootScope',
    '$timeout',
    '$interval',
    '$http',
    function TOTPChallenge($scope, $rootScope, $timeout, $interval, $http) {
      var waitSeconds = parseInt($rootScope.config.waitSeconds, 10);
      var start = Math.round(Date.now() / 1000);
      var end = start + waitSeconds + 2;
      var total = end - start;

      $scope.progress = 0;
      $scope.jsValue = '';
      $scope.formData = {
        totpValue: '',
      };

      // jsValue
      $scope.jsValueSolve = function jsValueSolve() {
        $http({
          url: window.config.baseURL + '/challenge',
          method: 'PATCH',
          data: {
            t: window.config.challengeToken,
          },
        })
          .then(function onResponse(response) {
            var data = response.data;
            window.config.challengeToken = data.t;
            // eslint-disable-next-line no-eval
            eval(window.atob(data.js));
            // eslint-disable-next-line vars-on-top
            var value = window.jsv();
            delete window.jsv;
            window.jsv = undefined;
            $scope.jsValue = value;
            $scope.captcha = data.captcha;
            document.getElementById('inputTOTP').focus();
          })
          .catch(function onError() {
            $rootScope.reload(3000);
          });
      };

      // solve
      $scope.solve = function solve(event) {
        $http({
          url: window.config.baseURL + '/challenge',
          method: 'POST',
          data: {
            t: window.config.challengeToken,
            j_v: $scope.jsValue,
            t_v: $scope.formData.totpValue.toString(),
          },
        })
          .then(function onResponse() {
            window.location.href = $rootScope.config.protectedPath;
          })
          .catch(function onError() {
            $rootScope.reload(3000);
          });
        if (event) {
          event.stopPropagation();
          event.preventDefault();
        }
        return false;
      };

      // progressUpdate interval then solve
      $scope.progressUpdate = function progressUpdate() {
        $timeout(function progress() {
          var state = Math.round(Date.now() / 1000) - start;
          var currentProgress = Math.round((state / total) * 100);
          if (currentProgress < 100) {
            $scope.progress = Math.floor(currentProgress);
            $scope.progressUpdate();
          } else {
            $scope.progress = 100;
            $scope.jsValueSolve();
          }
        }, 500);
      };

      $scope.progressUpdate();
      $scope.totpProgress = 0;
      $scope.totpProgressClass = '';

      $interval(function totpSeconds() {
        var seconds = new Date().getSeconds();
        var progress = Math.round((seconds / 30) * 100);
        if (seconds > 30) {
          progress = Math.round(((seconds - 30) / 30) * 100);
        }
        // eslint-disable-next-line vars-on-top
        var diff = Math.abs(50 - progress);
        $scope.totpProgress = progress;
        if (diff > 40) {
          $scope.totpProgressClass = 'is-danger';
        } else if (diff > 30) {
          $scope.totpProgressClass = 'is-warning';
        } else {
          $scope.totpProgressClass = 'is-success';
        }
      }, 500);
      // $scope.totpProgressInterval
    },
  ]);
})(window);
