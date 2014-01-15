'use strict';

angular.module('martoApp.controllers', []).

  /**
   * Main controller
   */
  controller('MainCtrl',
    ['$scope', 'Marto',
    function ($scope, Marto) {

      $scope.section = null;
      $scope.setSection = function(section) {
          $scope.section = section;
          console.log($scope.section);
      };

      $scope.logs = [];

      $scope.httpMethods = [
          { method: "GET" },
          { method: "POST" },
          { method: "PUT" },
          { method: "DELETE" },
      ];

      $scope.removeRequest = function(index) {
          $scope.scenario.requests.splice(index, 1);
      };

      $scope.run = function() {
          $scope.scenario.requests.forEach(function(request) {
              request.method = request.method.method;
          });
          console.log($scope.scenario);
          Marto.postScenario($scope.scenario);
      };

      Marto.on("request.started", function(event) {
          $scope.logs.push({ type: event.type });
      });
      Marto.on("response.finished", function(event) {
          $scope.logs.push({ type: event.type });
      });
  }]).

  /**
   * Scenarios controller
   */
  controller('ScenariosCtrl',
    ['$scope', 'Scenarios',
    function ($scope, Scenarios) {

      $scope.setSection('scenarios');

      $scope.scenarios = [];

      $scope.loadScenarios = function() {
        Scenarios.get({},
        function(response) {
          console.log(response);
          $scope.scenarios = response;
        },
        function(e) {
          console.log(e);
        });
      };
      $scope.loadScenarios();
  }]).

  /**
   * Scenario create controller
   */
  controller('ScenarioCreateCtrl',
    ['$scope', '$location', 'Marto', 'Scenario',
    function ($scope, $location, Marto, Scenario) {
      $scope.setSection('new-scenario');

      $scope.scenario = {
        id:    '',
        repeat: 1,
        every:  0,
        requests: []
      };

      $scope.addRequest = function() {
        $scope.scenario.requests.push({
          method: $scope.httpMethods[0],
          url:    '/',
          delay:  0
        });
      };
      $scope.addRequest();

      $scope.saveScenario = function() {
        $scope.scenario.requests.forEach(function(request) {
          request.method = request.method.method;
        });

        Scenario.save($scope.scenario, function(response) {
          $location.path('/scenarios');
        });
      };

      $scope.cancel = function() {
        $location.path('/scenarios');
      };
  }]).


  /**
   * Scenario remove controller
   */
  controller('ScenarioRemoveCtrl',
    ['$scope', '$location', '$routeParams', 'Scenario',
    function ($scope, $location, $routeParams, Scenario) {
      $scope.setSection('remove-scenario');

      Scenario.get({
        id: $routeParams.id
      }, function(response) {
        $scope.scenario = response;
      }, function() {
      });

      $scope.removeScenario = function(scenario) {
        console.log('remove scenario', scenario);
      };

      $scope.cancel = function() {
        $location.path('/scenarios');
      };
  }]).


  /**
   * Scenario edit controller
   */
  controller('ScenarioEditCtrl',
    ['$scope', '$location', '$routeParams', 'Marto', 'Scenario',
    function ($scope, $location, $routeParams, Marto, Scenario) {
      $scope.setSection('edit-scenario');

      Scenario.get({
        id: $routeParams.id
      }, function(response) {
        $scope.scenario = response;
      }, function() {
      });

      $scope.cancel = function() {
        $location.path('/scenarios');
      };
   }]);