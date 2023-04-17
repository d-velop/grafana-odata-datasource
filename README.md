# Grafana OData Data Source
Visualize data from OData data sources with Grafana.

## About
This is a Grafana data source for showing data from OData V4 compliant data sources.

It was originally developed for internal purposes and is now made available to the open source community.

## Getting started
Open Grafana and go to Configuration / Data Sources. Click `Add data source` button.

Enter `OData` into the filter input field and select the OData Data Source.

![Add Data Source](https://raw.githubusercontent.com/d-velop/grafana-odata-datasource/master/src/img/AddDataSource.png)

Provide the URL of your OData Service Root and click `Save & test` to test the connection.

Add other connection settings, such as auth settings, as necessary.

To use the data source, create a new query and select the newly created OData data source.

![CreateQuery.png](https://raw.githubusercontent.com/d-velop/grafana-odata-datasource/master/src/img/CreateQuery.png)

Choose an entity set, an appropriate time property, and the metric you want to view.
Now you should be able to see data for the selected time frame.

## Related Links
* [Grafana](https://grafana.com) - the open source analytics & monitoring solution for many data sources
* [Build a Grafana data source plugin](https://grafana.com/tutorials/build-a-data-source-plugin/) - a tutorial that 
  explains how to develop your own data source plugin.
* [OData](https://www.odata.org) - the ISO/IEC approved, OASIS standard for building and using data-driven RESTful APIs

## Contributing
This project is maintained by d-velop but is looking for contributors. If you consider contributing to this project
please read [CONTRIBUTING](https://raw.githubusercontent.com/d-velop/grafana-odata-datasource/master/CONTRIBUTING.md)
and [DEVELOPING](https://raw.githubusercontent.com/d-velop/grafana-odata-datasource/master/DEVELOPING.md) for details on
how to get started.

## License
Please read [LICENSE](https://raw.githubusercontent.com/d-velop/grafana-odata-datasource/master/LICENSE) for licensing
information.
