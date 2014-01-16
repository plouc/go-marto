'use strict';

describe('Scenarios controller', function () {

    var ctrl, scope, scenariosService;

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

        scenariosService = {
            get: function(scenario, successFn, errorFn) {
                successFn([]);
            }
        };
        spyOn(scenariosService, 'get').andCallThrough();

        ctrl  = $controller('ScenariosCtrl', {
            $scope:    scope,
            Scenarios: scenariosService
        });
    }));


    it('Should change current section', function () {
        expect(scope.section).toBe('scenarios');
    });


    it('Should load scenarios when called', function () {
        expect(scope.loadScenarios).toBeDefined();
        expect(scenariosService.get).toHaveBeenCalled();
    });
});