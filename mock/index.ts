import express, { Application, Request, Response } from "express";
import { v4 as uuidv4 } from "uuid";
import morgan from "morgan";

// @ts-ignore
import {DataServices, Key} from "./Metamodel";
// @ts-ignore
import {MockModel} from "./model/MockModel";
// @ts-ignore
import {ToXml} from "./util/XmlMetadataWriter";

const app: Application = express();
const port = 5000;

app.use(morgan('dev'));

app.get('/odata/\\$metadata',
  async (req: Request, res: Response): Promise<Response> => {
    // TODO: req.accepts('application/xml')
    return res
      .contentType('application/xml')
      .status(200)
      .send(ToXml(MockModel));
  }
);

app.get('/odata/temperatures',
  async (req: Request, res: Response): Promise<Response> => {
    // console.log(req.query);
    const { $filter } = req.query;
    if (typeof $filter == "string") {
      const regex : RegExp = /Time ge ([0-9-TZ:.]+) and Time le ([0-9-TZ:.]+)/;
      const match = $filter.match(regex);
      if (match)
      {
        const min = match[1];
        const max = match[2];
        console.log(`${min} / ${max}`);
      }
    }
    const count = 1000;
    let values = [];
    let startTime = (new Date()).getTime() - (count * 60 * 1000);
    for(let i = 0; i < count; i++)
    {
      let epochMs = startTime + (i * 60 * 1000) + Math.floor(Math.random() * 10000);
      let time = (new Date(epochMs)).toISOString();
      values.push({
        Id: uuidv4(),
        Time: time,
        Epoch: epochMs,
        Value1: Math.sin(i) + Math.random(),
        Value2: Math.cos(i) + Math.random(),
        Value3: Math.log2(i) + Math.random()
      });
    }
    return res
      .contentType('application/json')
      .status(200).send(
      {
        '@odata.context': 'http://localhost:5000/odata/$metadata#Temperatures',
        value: values
      });
  }
);

app.get('/odata',
  async (req: Request, res: Response): Promise<Response> => {
    let entitySets = [];
    for (let schema of MockModel.schemas)
    {
      let entityContainer = schema.entityContainer;
      if (entityContainer != null)
      {
        for (let entitySet of entityContainer.entitySets)
        {
          entitySets.push({name: entitySet.name, kind: 'EntitySet', url: entitySet.name});
        }
      }
    }
    return res
      .contentType('application/json')
      .status(200)
      .send({
        '@odata.context': 'http://localhost:5000/odata/$metadata',
        value: entitySets
      });
  }
);

try {
  app.listen(port, (): void => {
    console.log(`Listening on port ${port}`);
  });
} catch (error: any) {
  console.error(`Error occurred: ${error.message}`);
}
