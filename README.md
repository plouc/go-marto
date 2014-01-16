go-marto
========

go-marto is an http stress tool written in golang.

[![Build Status](https://travis-ci.org/plouc/go-marto.png?branch=master)](https://travis-ci.org/plouc/go-marto)

it's composed of various go package and one angular app:

  * **marto/** contains the main go package, you can use it as a standalone library,
    visit [doc](https://github.com/plouc/go-marto/tree/master/marto)
  * **api/** add a storage for scenarios
  * **app/** contains the angularjs app


Requirements
------------

  * go 1.2
  * mongodb


Launch the app
--------------

    $ go run main.go


Testing
-------

Testing requires additional dependencies:

  * nodejs
  * npm packages (see package.json)
  * phantomjs

To run all the tests (go + js), run:

    chmod +x test_all
    $ bash test_all


