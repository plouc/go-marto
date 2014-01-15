'use strict';

describe('Scenario edit controller', function () {

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
            get: function(scenario, callbackFn) {
                callbackFn({

                });
            }
        };
        spyOn(scenarioService, 'get').andCallThrough();;

        location = {
            path: function(path) {
                console.log('PATH', path);
            }
        };
        spyOn(location, 'path');

        ctrl  = $controller('ScenarioEditCtrl', {
            $scope:       scope,
            $location:    location,
            Scenario:     scenarioService,
            $routeParams: {
                id: 'test'
            }
        });
    }));


    it('Change current section', function () {
        expect(scope.section).toBe('edit-scenario');
    });


    it('Should redirect to scenarios when cancelling', function () {
        expect(scope.cancel).toBeDefined();

        scope.cancel();

        expect(location.path).toHaveBeenCalledWith('/scenarios');
    });
});