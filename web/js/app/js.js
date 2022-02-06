/* eslint-disable prefer-template */
/* eslint-disable prefer-arrow-callback */
/* eslint-disable no-param-reassign */
/* eslint-disable prefer-arrow-callback */
/* eslint-disable prefer-destructuring */
/* eslint-disable no-var */
(function application(window) {
  window.app.controller('js', [
    '$scope',
    '$rootScope',
    '$timeout',
    '$http',
    function JSChallenge($scope, $rootScope, $timeout, $http) {
      var waitSeconds = parseInt($rootScope.config.waitSeconds, 10);
      var start = Math.round(Date.now() / 1000);
      var end = start + waitSeconds + 2;
      var total = end - start;

      $scope.progress = 0;

      $scope.jsValue = '';

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
            $scope.solve();
          })
          .catch(function onError() {
            window.location.reload();
          });
      };

      // solve
      $scope.solve = function solve() {
        $http({
          url: window.config.baseURL + '/challenge',
          method: 'POST',
          data: {
            t: window.config.challengeToken,
            j_v: $scope.jsValue,
          },
        })
          .then(function onResponse() {
            window.location.href = $rootScope.config.protectedPath;
          })
          .catch(function onError() {
            $rootScope.reload();
          });
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
    },
  ]);
})(window);
