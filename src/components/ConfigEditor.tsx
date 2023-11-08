import {
  DataSourcePluginOptionsEditorProps,
  SelectableValue
} from '@grafana/data';
import {DataSourceHttpSettings, FieldSet, InlineField, InlineFieldRow, Select} from '@grafana/ui';
import React, {ComponentType, useCallback} from 'react';
import {ODataOptions, URLSpaceEncoding} from '../types';

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
