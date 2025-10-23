[![Build Status](https://github.com/foomo/contentful/actions/workflows/test.yml/badge.svg?branch=main&event=push)](https://github.com/foomo/contentful/actions/workflows/pr.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/foomo/contentful)](https://goreportcard.com/report/github.com/foomo/contentful)
[![Coverage Status](https://coveralls.io/repos/github/foomo/contentful/badge.svg?branch=main&)](https://coveralls.io/github/foomo/contentful?branch=main)
[![GoDoc](https://godoc.org/github.com/foomo/contentful?status.svg)](https://godoc.org/github.com/foomo/contentful)

<p align="center">
  <img alt="sesamy" src=".github/assets/gocontentful.png"/>
</p>

# contentful-go

GoLang SDK for [Contentful's](https://www.contentful.com) Content Delivery, Preview and Management API's.

# About

[Contentful](https://www.contentful.com) provides a content infrastructure for digital teams to power content in websites, apps, and devices. Unlike a CMS, Contentful was built to integrate with the modern software stack. It offers a central hub for structured content, powerful management and delivery APIs, and a customizable web app that enable developers and content creators to ship digital products faster.

[Go](https://golang.org) is an open source programming language that makes it easy to build simple, reliable, and efficient software.

# Install

`go get github.com/foomo/contentful`

# Getting started

Import into your Go project or library

```go
import (
	"github.com/foomo/contentful"
)
```

Create a API client in order to interact with the Contentful's API endpoints.

```go
token := "your-cma-token" // observe your CMA token from Contentful's web page
cma := contentful.NewCMA(token)
```

#### Organization

If your Contentful account is part of an organization, you can setup your API client as so. When you set your organization id for the SDK client, every api request will have `X-Contentful-Organization: <your-organization-id>` header automatically.

```go
cma.SetOrganization("your-organization-id")
```

#### Debug mode

When debug mode is activated, sdk client starts to work in verbose mode and try to print as much informatin as possible. In debug mode, all outgoing http requests are printed nicely in the form of `curl` command so that you can easly drop into your command line to debug specific request.

```go
cma.Debug = true
```

# Using the SDK

## Working with resource services

Currently, SDK exposes the following resource services:

* Spaces
* APIKeys
* Assets
* ContentTypes
* Entries
* Locales
* Webhooks

Every resource service has at least the following interface:

```go
List(ctx) *Collection
Get(ctx, spaceID, resourceID string) <Resource>, error
Upsert(ctx, spaceID string, resourceID *Resource) error
Delete(ctx, spaceID string, resourceID *Resource) error
```

#### Example

```go
space, err := cma.Spaces.Get(ctx, "space-id")
if err != nil {
  log.Fatal(err)
}

collection := cma.ContentTypes.List(ctx, space.Sys.ID)
collection, err := collection.Next()
if err != nil {
  log.Fatal(err)
}

for _, contentType := range collection.ToContentType() {
  fmt.Println(contentType.Name, contentType.Description)
}
```

## Working with collections

All the endpoints which return an array of objects are wrapped around `Collection` struct. The main features of `Collection` are pagination and type assertion.

### Pagination
WIP

### Type assertion

`Collection` struct exposes the necessary converters (type assertion) such as `ToSpace()`. The following example gets all spaces for the given account:

### Example

```go
collection := cma.Spaces.List() // returns a collection
collection, err := collection.Next() // makes the actual api call
if err != nil {
  log.Fatal(err)
}

spaces := collection.ToSpace() // make the type assertion
for _, space := range spaces {
  fmt.Println(space.Name)
  fmt.Println(space.Sys.ID)
}

// In order to access collection metadata
fmt.Println(col.Total)
fmt.Println(col.Skip)
fmt.Println(col.Limit)
```

## Testing

```shell
$> go test
```

To enable higher verbose mode

```shell
$> go test -v -race
```

## Documentation/References

### Contentful
[Content Delivery API](https://www.contentful.com/developers/docs/references/content-delivery-api/)
[Content Management API](https://www.contentful.com/developers/docs/references/content-management-api/)
[Content Preview API](https://www.contentful.com/developers/docs/references/content-preview-api/)

### GoLang
[Effective Go](https://golang.org/doc/effective_go.html)

## Support

This is a project created for demo purposes and not officially supported, so if you find issues or have questions, you can let us know via the [issue](https://github.com/foomo/contentful/issues/new) page.

## How to Contribute

Please refer to the [CONTRIBUTING](.github/CONTRIBUTING.md) details and follow the [CODE_OF_CONDUCT](.github/CODE_OF_CONDUCT.md) and [SECURITY](.github/SECURITY.md) guidelines.

## License

Distributed under MIT License, please see license file within the code for more details.

_Made with ♥ [foomo](https://www.foomo.org) by [bestbytes](https://www.bestbytes.com)_
