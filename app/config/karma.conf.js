module.exports = function(config){
  config.set({
    basePath : '../',

    colors: true,
    logLevel: config.LOG_INFO,

    files : [
      'js/vendor/angular.min.js',
      'js/vendor/angular-*.js',
      'test/lib/angular-mocks.js',
      'js/*.js',
      'test/unit/**/*.js'
    ],

    exclude : [
      'js/angular/angular-loader.js',
      'js/angular/*.min.js',
      'js/angular/angular-scenario.js'
    ],

    autoWatch : true,

    frameworks: [
      'jasmine'
    ],

    browsers: [
      'Chrome',
      //'PhantomJS'
    ],

    plugins : [
      'karma-spec-reporter',
      'karma-junit-reporter',
      'karma-chrome-launcher',
      'karma-phantomjs-launcher',
      'karma-jasmine'
    ],

    junitReporter : {
      outputFile: 'test_out/unit.xml',
      suite:      'unit'
    },

    reporters: ['spec']
})}