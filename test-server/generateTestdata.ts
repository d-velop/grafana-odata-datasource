const fs = require('fs');
const {Generate} = require("./mock/Testdata");  // npm install uuid

const fileName = 'db/data/test-Temperatures.csv';
const header = "id,time,epoch,value1,value2,value3\n";
fs.writeFileSync(fileName, header);

let values = Generate();
for(let i = 0; i < values.length; i++)
{
  const value = values[i];
  const line = `${value.id},${value.time},${value.epoch},${value.value1},${value.value2},${value.value3}\n`;
  fs.appendFileSync(fileName, line);
}
