import {v4 as uuidv4} from "uuid";

export const PRIMITIVES_START = '2024-01-01T00:00:00.000Z';

export function GenerateTestPrimitives() {
  const count = 100;
  const startMs = new Date(PRIMITIVES_START).getTime();
  const values = [];
  for (let i = 0; i < count; i++) {
    const dt = new Date(startMs + i * 60_000);
    values.push({
      guid:            uuidv4(),
      dateTimeOffset:  dt.toISOString(),
      date:            dt.toISOString().substring(0, 10),
      int64:           i * 1_000_000_000,
      int32:           i,
      int16:           i,
      decimal:         (i * 0.5).toFixed(3),
      double:          Math.sin(i * Math.PI / 50),
      boolean:         i % 2 === 0,
      string:          `item-${String(i).padStart(3, '0')}`,
    });
  }
  return values;
}

const UNITS = ['C', 'F', 'K'];

export function GenerateTemperatures()
{
  const count = 1000;
  let values = [];
  let startTime = (new Date()).getTime() - (count * 60 * 1000);
  for(let i = 0; i < count; i++)
  {
    let epochMs = startTime + (i * 60 * 1000) + Math.floor(Math.random() * 10000);
    let time = (new Date(epochMs)).toISOString();
    let microseconds = String(Math.floor(Math.random() * 1000)).padStart(3, '0');
    let sampledAt = time.replace(/\.(\d{3})Z$/, `.$1${microseconds}Z`);
    let measurementDate = time.substring(0, 10);
    let sensorId = (i % 5) + 1;
    values.push({
      id: uuidv4(),
      time: time,
      sampledAt: sampledAt,
      measurementDate: measurementDate,
      epoch: epochMs,
      sensorId: sensorId,
      qualityCode: Math.random() < 0.9 ? 0 : (Math.random() < 0.5 ? 1 : 2),
      value1: Math.sin(i) + Math.random(),
      value2: Math.cos(i) + Math.random(),
      value3: Math.log2(i + 1) + Math.random(),
      pressure: Math.round((1013.25 + (Math.random() - 0.5) * 20) * 100) / 100,
      isOutdoor: sensorId <= 3,
      unit: UNITS[i % UNITS.length]
    });
  }
  return values;
}

export function GenerateRooms()
{
  const count = 100;
  let values = [];
  for(let i = 0; i < count; i++)
  {
    values.push({
      id: uuidv4(),
      name: `Room ${i}`,
    });
  }
  return values;
}
