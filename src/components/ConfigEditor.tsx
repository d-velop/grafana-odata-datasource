import {
  DataSourcePluginOptionsEditorProps,
  SelectableValue
} from '@grafana/data';
import {DataSourceHttpSettings, FieldSet, InlineField, InlineFieldRow, Select} from '@grafana/ui';
import React, {ComponentType, useCallback} from 'react';
import {ODataOptions, URLSpaceEncoding, ODataVersion} from '../types';

type Props = DataSourcePluginOptionsEditorProps<ODataOptions>;

export const ConfigEditor: ComponentType<Props> = ({ options, onOptionsChange }) => {
  const onURLSpaceEncodingChange = useCallback((option: SelectableValue<URLSpaceEncoding>) => {
      const urlSpaceEncoding = option.value;
      onOptionsChange({
        ...options,
        jsonData: {
          ...options.jsonData,
          urlSpaceEncoding: urlSpaceEncoding || '+',
        },
      });
  }, [onOptionsChange, options]);

  const urlSpaceEncodings = Object.entries(URLSpaceEncoding)
    .map(([label, value]) => ({ label: `${label} (${value})`, value: value }));

  const onODataVersionChange = useCallback((option: SelectableValue<ODataVersion>) => {
    const odataVersion = option.value;
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        odataVersion: odataVersion || 'Auto',
      },
    });
  }, [onOptionsChange, options]);

  const odataVersions = Object.entries(ODataVersion)
    .map(([label, value]) => ({ label: label, value: value }));

  return (
    <>
      <DataSourceHttpSettings
        defaultUrl={'http://localhost:5000/odata'}
        dataSourceConfig={options}
        showAccessOptions={false}
        onChange={onOptionsChange}
      />
      <div className='gf-form-group'>
        <h3 className='page-heading'>Additional settings</h3>
        <FieldSet>
          <InlineFieldRow>
            <InlineField
              label='OData Version'
              labelWidth={26}
              tooltip={
                <p>
                  Select the OData version, currently V2, V3 and V4 are supported. The plugin currently only supports
                  XML format for metadata (`$metadata`) and JSON format for payload data for both OData versions.
                </p>
              }>
              <Select
                options={odataVersions}
                value={options.jsonData.odataVersion?.length > 0
                  ? odataVersions.find((type) => type.value === options.jsonData.odataVersion)
                  : ODataVersion.Auto}
                className='width-10'
                onChange={onODataVersionChange}
              />
            </InlineField>
          </InlineFieldRow>
          <InlineFieldRow>
            <InlineField
              label='URL space encoding'
              labelWidth={26}
              tooltip={
                <p>
                  Select the standard for encoding spaces in URLs. <i>Percent</i> uses <code>%20</code> (see RFC 3986),
                  while <i>Plus</i> uses <code>+</code> (used in form data). E.g. <code>$filter=value%20EQ%201</code>
                  (Percent) and <code>`$filter=value+EQ+1`</code> (Plus).
                </p>
              }>
              <Select
                options={urlSpaceEncodings}
                value={options.jsonData.urlSpaceEncoding?.length > 0
                  ? urlSpaceEncodings.find((type) => type.value === options.jsonData.urlSpaceEncoding)
                  : URLSpaceEncoding.Plus}
                className='width-10'
                onChange={onURLSpaceEncodingChange}
              />
            </InlineField>
          </InlineFieldRow>
        </FieldSet>
      </div>
      </>
  );
};
