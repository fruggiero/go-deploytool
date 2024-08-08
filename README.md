# go-deploytool

This tool, written in Go, was created to be used within a continuous delivery pipeline in a CI/CD context. Specifically, it has been tested on a GitLab pipeline.
The purpose is to publish a nupkg package to a NuGet server, passing only an identifier name for the feed (not the URL, which will be taken from a configuration file called config.json). 
To resolve credentials, the HashiCorp Vault tool is used.
The goal is to make the publication of a NuGet package both simple and secure, without giving the user the ability to know the credentials used for publication, which will be resolved dynamically through Vault.
