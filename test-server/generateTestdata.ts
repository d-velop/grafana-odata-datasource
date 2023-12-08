const fs = require('fs');
const {GenerateTemperatures, GenerateRooms} = require("./mock/Testdata");

function writeToCSV(fileName: string, header: string, dataGenerator: () => any[]) {
  let data = dataGenerator().map(value => Object.values(value).join(',')).join('\n');
  fs.writeFileSync(fileName, header + data + '\n');
}

writeToCSV('db/data/test-Temperatures.csv', 'id,time,epoch,value1,value2,value3\n', GenerateTemperatures);
writeToCSV('db/data/test-Rooms.csv', 'id,name\n', GenerateRooms);
