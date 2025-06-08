# Start Developing
This plugin is a data source backend plugin. It uses [Grafana Plugin Tools](https://grafana.github.io/plugin-tools/) and
consists of both backend and frontend components.

## Prerequisites
For this project to work you need [Node.js](https://nodejs.org/en/) and [Go](https://go.dev) installed.

As with other Grafana data source plugin projects, [yarn](https://yarnpkg.com/) is used for managing and building the
frontend and [Mage](https://magefile.org) for the backend.

Additionally, [Docker](https://www.docker.com/) is used to simplify plugin development and integration testing.

## Getting started

### Clone and build
Clone this repository into your local environment
```bash
git clone https://github.com/d-velop/grafana-odata-datasource.git
```

#### Backend
Backend code is located in the `pkg` folder

Build backend plugin binaries by running the mage build (`-v` stands for verbose output)
```bash
mage -v
```

#### Frontend
Frontend code is located in the `src` folder

Install dependencies
```bash
yarn install
```

Build plugin in development mode
```bash
yarn dev
```

or build plugin in production mode
```bash
yarn build
```

### Using the OData test-server
If you don't have a full-fledged OData server at hand, you will find a test server based on
[Express](https://expressjs.com) and on the Core Data Services (CDS) of the
[SAP Cloud Application Programming Model](https://cap.cloud.sap/) in the [`test-server`](test-server) directory.

It can be started by using the command
```bash
cd test-server
yarn start 
```

In addition, the test server is automatically built and started using Docker Compose (see below).

For more information see [test-server/README.md](test-server/README.md).

### Try and test using Docker Compose
The project includes a [`docker-compose.yaml`](docker-compose.yaml) file. With this, Grafana can be started quickly for
development purposes. The local project directory is automatically mounted to the Grafana plugin directory.

Additionally, to keep development uncomplicated, anonymous authorization is enabled in Grafana. The project also comes
with predefined data source configurations (using the aforementioned test-server) and test dashboards that allow changes
to be tried out and tested directly. See folder [`provisioning`](provisioning) for details.

To start, simply run the following command
```bash
yarn server
```
in the projects root directory. Afterwards you can open `http://localhost:3000` in your browser and begin using Grafana
with the preconfigured OData Data Source.

> Note: If you want to access a locally running OData Service Root make sure you use the correct hostname. You can use
> docker's special DNS name `host.docker.internal` which resolves to the internal IP address used by the docker host.

## Testing

Run all backend test by executing the following command:

```bash
mage test
```

### Coverage

To evaluate the backend test coverage execute the following command:

```bash
mage coverage
```

The results are written to a `backend.html` file located in the [`./coverage`](./coverage) folder.

### Testing the `oauthPassThru` feature of the plugin

The plugin supports the `oauthPassThru` feature (Forward OAuth Identity), which allows Grafana to forward the user's
upstream OAuth identity to the data source. This enables the backend to make requests on behalf of the user, useful for
scenarios where per-user access control is needed.

The project includes a [`docker-compose.withKeycloak.yaml`](docker-compose.withKeycloak.yaml) file for running the
plugin together with a [Keycloak](https://www.keycloak.org/) instance.

Keycloak is used here as an example identity provider implementing OAuth2 and OpenID Connect.

The [integration/keycloak](integration/keycloak) folder contains configuration files for Keycloak and Grafana:

- `realm.json`: Keycloak realm configuration for testing purposes.
- `grafana.ini`: Grafana configuration enabling OAuth integration.

You can use the following command as a shortcut to run the Docker Compose setup:
```bash
yarn server:withKeycloak
```

To enable integrated testing across containers and from the host system, add the following entry to your `/etc/hosts`:
```
127.0.0.1 dockerhost
```

This ensures that all services are reachable consistently:
* Grafana: http://dockerhost:3000
* Keycloak: http://dockerhost:8080
* Test server: http://dockerhost:4004

This is required for the authentication flow to work correctly, both from the browser on the host system and for
internal communication between services within the Docker network.

## Update create-plugin versions
To update the plugin to use a newer version of the `create-plugin` tool, follow the instructions here:
<https://grafana.com/developers/plugin-tools/migration-guides/update-create-plugin-versions>.

The source code of the `create-plugin` tool can be found here:
<https://github.com/grafana/plugin-tools/tree/main/packages/create-plugin>.
