const testServerUrl = process.env.TEST_SERVER_URL || 'http://localhost:4004';

async function get(path: string): Promise<Response> {
  const res = await fetch(`${testServerUrl}${path}`);
  if (!res.ok) {
    throw new Error(`test-server: GET ${path} returned ${res.status}`);
  }
  return res;
}

async function expectNonEmptyValue(path: string) {
  const body = await (await get(path)).json();
  if (!Array.isArray(body.value) || body.value.length === 0) {
    throw new Error(`test-server: GET ${path} returned no entities`);
  }
}

export default async function globalSetup() {
  const roomCount = await (await get('/odata/v4/test/Rooms/$count')).text();
  if (roomCount !== '100') {
    throw new Error(`test-server: expected 100 Rooms, got ${roomCount}`);
  }
  await expectNonEmptyValue('/odata/v4/test/Temperatures?$top=1');
  await get('/mock/$metadata');
  await expectNonEmptyValue('/mock/temperatures');
}
