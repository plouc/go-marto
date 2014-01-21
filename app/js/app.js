'use strict';

angular.module('martoApp',
  [
    'martoApp.services',
    'martoApp.controllers',
    'ngRoute'
  ]).

  config(['$routeProvider', function($routeProvider) {
    $routeProvider.
    when('/scenarios', {
      templateUrl: 'partials/scenarios.html',
      controller:  'ScenariosCtrl'
    }).
    when('/scenarios/new', {
      templateUrl: 'partials/scenario-form.html',
      controller:  'ScenarioCreateCtrl'
    }).
    when('/scenarios/:id/edit', {
      templateUrl: 'partials/scenario-form.html',
      controller:  'ScenarioEditCtrl'
    }).
    when('/scenarios/:id/remove', {
      templateUrl: 'partials/scenario-remove.html',
      controller:  'ScenarioRemoveCtrl'
    }).
    when('/scenarios/:id/run', {
        templateUrl: 'partials/scenario-run.html',
        controller:  'ScenarioRunCtrl'
    }).
    otherwise({
      redirectTo: '/scenarios'
    });
  }]).

  factory('Marto', ['$q', '$rootScope', function($q, $rootScope) {

    // We return this object to anything injecting our service
    var Service = {};

    // Keep all pending requests here until they get responses
    var callbacks = {};

    // Store all event listeners
    var listeners = {};

    // Create a unique callback ID to map requests to responses
    var currentCallbackId = 0;

    // Create our websocket object with the address to the websocket
    var ws = new WebSocket("ws://localhost:8182/marto");
    
    ws.onopen = function(){  
        console.log("Socket has been opened!");  
    };
    
    ws.onmessage = function(message) {
    	console.log('onmessage');
      listener(JSON.parse(message.data));
    };

    function sendRequest(request) {
      console.log('sendRequest', request);
      var defer = $q.defer();
      var callbackId = getCallbackId();
      callbacks[callbackId] = {
        time: new Date(),
        cb:defer
      };
      request.callback_id = callbackId;
      console.log('Sending request', request);
      ws.send(JSON.stringify(request));
      return defer.promise;
    }

    function listener(data) {
      var messageObj = data;
      console.log("Received data from websocket: ", messageObj);

      if (data.hasOwnProperty("type") && listeners.hasOwnProperty(data.type)) {
        listeners[data.type].forEach(function(listener) {
          $rootScope.$apply(listener(data));
        });
      }

      // If an object exists with callback_id in our callbacks object, resolve it
      if (callbacks.hasOwnProperty(messageObj.callback_id)) {
        console.log("Registered callback");
        console.log(callbacks[messageObj.callback_id]);
        $rootScope.$apply(callbacks[messageObj.callback_id].cb.resolve(messageObj.data));
        delete callbacks[messageObj.callbackID];
      } else {
        console.log("No registered callback");
      }
    }

    // This creates a new callback ID for a request
    function getCallbackId() {
      currentCallbackId += 1;
      if (currentCallbackId > 10000) {
        currentCallbackId = 0;
      }

      return currentCallbackId;
    }

    // Define a "getter" for getting customer data
    Service.postScenario = function(scenario) {

      var event = {
        "type":    "scenario.add",
        "payload": scenario
      };

      // Storing in a variable for clarity on what sendRequest returns
      var promise = sendRequest(event); 
      return promise;
    }

    Service.on = function(eventType, callback) {
      if (!listeners.hasOwnProperty(eventType)) {
        listeners[eventType] = [];
      }

      listeners[eventType].push(callback);
    };

    Service.off = function(eventType) {
      delete listeners[eventType];
    };

    return Service;
  }]);