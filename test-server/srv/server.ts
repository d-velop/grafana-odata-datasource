import cds from "@sap/cds";
import morgan from "morgan";
import {Application} from "express";
import proxy from "@sap/cds-odata-v2-adapter-proxy";
// @ts-ignore
import {addMockService} from "../mock/MockService";

cds.on('bootstrap', async (app: Application) => {
  app.use(morgan('dev'));
  app.use(proxy());
  await addMockService(app);
})
