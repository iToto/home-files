# yield-mvp

Yield MVP service

## Overview

The `yield-mvp` application is designed to provide a minimal viable product (MVP) for yield-related services. It is a backend service that handles various operations related to yield management.

## Configurations

The application requires a configuration file for the appropriate environment to run. The supported environments are:

- Local
- Production

The configuration files are stored in the `configs/` directory in the format `env.{ENVIRONMENT}`. A starting environment file with all required keys can be found at `configs/env.example`.

## Running the Application Locally

To run the application locally, use the following command:

```sh
$ make run
```

This command will start the application using the local configuration.

# Features

- Yield Management: Core functionality for managing yield-related operations.
- Environment Configurations: Supports multiple environments with separate configuration files.
- Easy Setup: Simple commands to run the application locally.

# System Design

This application has multiple services packaged into one:

* Trading app
* Signal logger
* HTTP API

Here is a brief description of each:

## Trading App

The trading app runs on a continuous loop whose cadence is determined by a configurable timer. Every loop, it will query the DB for all the active signals and their paired active strategies. It will them loop through each signal, request the current signal, log it, then loop through each strategy and see if there is an action to take. Each strategy is configurable to run different heuristics. These are determined by the properties set on the `entities.Strategy` Object (persisted in the DB).

The app and it's dependencies are setup in `main.go`. It does quite a bit:

- Loads in secrets from environments
- Sets up DB connection
- Setup and connect web sockets for real time BTC and ETH price checks
- Sets the user (right now only setup for 1 user: Dan)
- Sets up Data Access Layer (DAL) which is what is used to retrieve data from dependencies
- 

# Dependencies

## API

You can find all HTTP API dependencies in: `pkg/`. The main one is `coinroutes` which uses API keys and secrets generated by Yield and used by the app.

There is also a dependency on the Signal APIs which is a generic API package that takes in a URL (signals set by user at runtime) makes a request and returns the signal response to be processed by the application.

## GCP

This application run on the Google Cloud Platform and is dependant on some services there. Most notably Cloud Run. 

# Directory Structure

### configs/: Contains configuration files for different environments

### handlers/: HTTP handlers for the application's endpoints

### models/: Data structures and database models

### services/: Business logic and core functionalities

### utils/: Utility functions and helpers

### middleware/: Middleware functions for request processing

### db/: Database interaction functions

### routes/: Routing configuration for the application

### pkg/: Dependencies and helpers for the application





Getting Started

Clone the repository:
```sh
git clone https://github.com/yourusername/yield-mvp.git
cd yield-mvp
```

Copy the example environment file and modify it as needed:

```sh
cp configs/env.example configs/env.local
```

Install dependencies:

```sh
make install
```

Run the application:

```sh
make run
```
