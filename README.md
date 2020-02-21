[![Go Report Card](https://goreportcard.com/badge/github.com/stiks/gobs?style=flat-square)](https://goreportcard.com/report/github.com/stiks/gobs)
[![Build Status](https://travis-ci.org/stiks/gobs.svg?branch=master)](https://travis-ci.org/stiks/gobs)
[![Codecov](https://codecov.io/gh/stiks/gobs/branch/master/graph/badge.svg)](https://codecov.io/gh/stiks/gobs)
[![Maintainability](https://api.codeclimate.com/v1/badges/4f69faf9cf3186f85943/maintainability)](https://codeclimate.com/github/stiks/gobs/maintainability)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/stiks/gobs/master/LICENSE)

This is not yet another Go REST framework. The goal for this project is to come as close as possible to [twelve-factor app](https://12factor.net/) methodology.

#### Goals
* Declarative documentation and automation
* Maximum portability
* Support public clouds, [kubernetes](https://kubernetes.io/) and on-premise deployments
* Easy to scale
* Have UI ([gobs-react](https://github.com/stiks/gobs-react)) with functionality of typical CMS

#### Inspired by
* Domain Driven Design
* Clean Architecture

### Features
* High performance, minimalist web framework ([Echo v4](https://github.com/labstack/echo))
* Transactional emails ([Hermes v2](https://github.com/matcornic/hermes))
* JWT token authorisation
* Go modules

#### Structure

With Go, there's no real standard folder structure for the project.

```
├── app
├── lib
│   ├── controllers
│   ├── models
│   ├── providers
│   │   ├── appengine
│   │   ├── dummy
│   │   └── mock
│   ├── repositories
│   └── services
├── pkg
│   ├── auth
│   ├── env
│   ├── helpers
│   ├── parser
│   └── xlog
└── vendor
```

#### Auth

Auth service inspired by [go-oauth2-server](https://github.com/RichardKnop/go-oauth2-server)
