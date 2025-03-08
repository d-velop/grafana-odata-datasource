import {Application, Request, Response} from "express";
import {ToXml} from "./util/XmlMetadataWriter";
import {MockModel} from "./model/MockModel";
import {GenerateTemperatures} from "./Testdata";

export let addMockService = async (app: Application) => {
  app.get('/mock/\\$metadata',
    async (_: Request, res: Response): Promise<void> => {
      res
        .contentType('application/xml')
        .status(200)
        .send(ToXml(MockModel));
    }
  );

  app.get('/mock/temperatures',
    async (req: Request, res: Response): Promise<void> => {
      const {$filter} = req.query;
      if (typeof $filter === "string") {
        const regex = /Time ge ([0-9-TZ:.]+) and Time le ([0-9-TZ:.]+)/;
        const match = $filter.match(regex);
        if (match) {
          const min = match[1];
          const max = match[2];
          console.log(`${min} / ${max}`);
        }
      }
      let values = GenerateTemperatures();
      res
        .contentType('application/json')
        .status(200).send(
          {
            '@odata.context': 'http://localhost:4004/odata/$metadata#Temperatures',
            value: values
          });
    }
  );

  app.get('/mock',
    async (_: Request, res: Response): Promise<void> => {
      let entitySets = [];
      for (let schema of MockModel.schemas) {
        let entityContainer = schema.entityContainer;
        if (entityContainer != null) {
          for (let entitySet of entityContainer.entitySets) {
            entitySets.push({name: entitySet.name, kind: 'EntitySet', url: entitySet.name});
          }
        }
      }
      res
        .contentType('application/json')
        .status(200)
        .send({
          '@odata.context': 'http://localhost:4004/odata/$metadata',
          value: entitySets
        });
    }
  );
}
