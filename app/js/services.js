'use strict';

angular.module('martoApp.services', ['ngResource']).

  factory('Scenarios', ['$resource', function($resource){
    return $resource('/api/scenarios', {}, {
      get: { method:'GET', isArray: true }
    });
  }]).

  factory('Scenario', ['$resource', function($resource){
    return $resource('/api/scenarios', {}, {
      save: { method:'POST' }
    });
  }]).

  factory('Scenario', ['$resource', function($resource){
    return $resource('/api/scenarios/:id', {}, {
      get: { method:'GET', params: { id: 'id' } }
    });
  }]);