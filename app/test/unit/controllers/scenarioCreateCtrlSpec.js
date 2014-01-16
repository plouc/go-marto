'use strict';

describe('Scenario create controller', function () {

  var ctrl, scope, scenarioService, location;

  beforeEach(module('martoApp'));
  beforeEach(inject(function ($rootScope, $controller) {
    scope = $rootScope.$new();
    scope.httpMethods = [
      { method: "GET" },
      { method: "POST" },
      { method: "PUT" },
      { method: "DELETE" }
    ];
    scope.setSection = function(section) {
      this.section = section;
    };

    scenarioService = {
      save: function(scenario, callbackFn) {
        callbackFn();
      }
    };
    spyOn(scenarioService, 'save').andCallThrough();;

    location = {
      path: function(path) {
        console.log('PATH', path);
      }
    };
    spyOn(location, 'path');

    ctrl  = $controller('ScenarioCreateCtrl', {
      $scope:    scope,
      $location: location,
      Marto:     {},
      Scenario:  scenarioService
    });
  }));


  it('Should change current section', function () {
    expect(scope.section).toBe('new-scenario');
  });


  it('Has a default scenario set with one request', function () {
    expect(scope.scenario).toEqual({
      id :       '',
      repeat :   1,
      every :    0,
      requests : [
        { method : { method : 'GET' }, url : '/', delay : 0 }
      ]
    });
  });


  it('Should a save function defined', function () {
    expect(scope.saveScenario).toBeDefined();
  });


  it('Should add request when addRequest is called', function () {
    expect(scope.addRequest).toBeDefined();

    scope.addRequest();

    expect(scope.scenario).toEqual({
      id :       '',
      repeat :   1,
      every :    0,
      requests : [
        { method : { method : 'GET' }, url : '/', delay : 0 },
        { method : { method : 'GET' }, url : '/', delay : 0 }
      ]
    });
  });


  it('Should save scenario and redirect to scenarios', function () {
    expect(scope.saveScenario).toBeDefined();

    scope.saveScenario();

    expect(scenarioService.save).toHaveBeenCalled();
    expect(scope.scenario).toEqual({
      id :       '',
      repeat :   1,
      every :    0,
      requests : [
        { method : 'GET', url : '/', delay : 0 }
      ]
    });
    expect(location.path).toHaveBeenCalledWith('/scenarios');
  });


  it('Should redirect to scenarios when cancelling', function () {
    expect(scope.cancel).toBeDefined();

    scope.cancel();

    expect(location.path).toHaveBeenCalledWith('/scenarios');
  });
});