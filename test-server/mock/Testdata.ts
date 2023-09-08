import {v4 as uuidv4} from "uuid";

export function Generate()
{
  const count = 1000;
  let values = [];
  let startTime = (new Date()).getTime() - (count * 60 * 1000);
  for(let i = 0; i < count; i++)
  {
    let epochMs = startTime + (i * 60 * 1000) + Math.floor(Math.random() * 10000);
    let time = (new Date(epochMs)).toISOString();
    values.push({
      id: uuidv4(),
      time: time,
      epoch: epochMs,
      value1: Math.sin(i) + Math.random(),
      value2: Math.cos(i) + Math.random(),
      value3: Math.log2(i) + Math.random()
    });
  }
  return values;
}
