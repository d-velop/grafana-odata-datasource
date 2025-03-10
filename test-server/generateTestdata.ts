const fs = require('fs');
const {GenerateTemperatures, GenerateRooms} = require("./mock/Testdata");

function writeToCSV(name: string, dataGenerator: () => any[]) {
  const dataArray = dataGenerator();
  const header = Object.keys(dataArray[0]).join(',');
  const data = dataArray.map(value => Object.values(value).join(',')).join('\n');
  fs.writeFileSync(`db/data/test-${name}.csv`, header + '\n' + data + '\n');
}

writeToCSV('Temperatures', GenerateTemperatures);
writeToCSV('Rooms', GenerateRooms);
