import type { PluginOptions } from '@grafana/plugin-e2e';
import { defineConfig, devices } from '@playwright/test';

export default defineConfig<PluginOptions>({
  testDir: './integration/keycloak/tests',
  retries: 0,
  workers: 1,
  reporter: 'list',
  use: {
    baseURL: 'http://dockerhost:3000',
    launchOptions: {
      args: ['--host-resolver-rules=MAP dockerhost 127.0.0.1'],
    },
    trace: 'retain-on-failure',
  },
  projects: [
    {
      name: 'keycloak-login',
      testMatch: /keycloak\.setup\.ts/,
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'chromium',
      testMatch: /\.spec\.ts/,
      use: {
        ...devices['Desktop Chrome'],
        storageState: 'playwright/.auth/keycloak-user.json',
      },
      dependencies: ['keycloak-login'],
    },
  ],
});
