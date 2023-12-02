const fs = require('fs');
const {GenerateTemperatures, GenerateRooms} = require("./mock/Testdata");  // npm install uuid

// Temperatures
let fileName = 'db/data/test-Temperatures.csv';
let header = "id,noTimeProperty,epoch,value1,value2,value3\n";
fs.writeFileSync(fileName, header);

let values = GenerateTemperatures();
for(let i = 0; i < values.length; i++)
{
  const value = values[i];
  const line = `${value.id},${value.time},${value.epoch},${value.value1},${value.value2},${value.value3}\n`;
  fs.appendFileSync(fileName, line);
}

// Rooms
fileName = 'db/data/test-Rooms.csv';
header = "id,name\n";
fs.writeFileSync(fileName, header);

values = GenerateRooms();
for(let i = 0; i < values.length; i++)
{
  const value = values[i];
  const line = `${value.id},${value.name}\n`;
  fs.appendFileSync(fileName, line);
}
