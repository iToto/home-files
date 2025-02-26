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

This command will start the application using the local configuration.

## Features

- Yield Management: Core functionality for managing yield-related operations.
- Environment Configurations: Supports multiple environments with separate configuration files.
- Easy Setup: Simple commands to run the application locally.


## Directory Structure

### configs/: Contains configuration files for different environments.

### handlers/: HTTP handlers for the application's endpoints.

### models/: Data structures and database models.

### services/: Business logic and core functionalities.

### utils/: Utility functions and helpers.

### middleware/: Middleware functions for request processing.

### db/: Database interaction functions.

### routes/: Routing configuration for the application.

Getting Started

Clone the repository:
```
git clone https://github.com/yourusername/yield-mvp.git
cd yield-mvp
```

Copy the example environment file and modify it as needed:

```
cp configs/env.example configs/env.local
```

Install dependencies:

```
make install
```

Run the application:

```
make run
```
