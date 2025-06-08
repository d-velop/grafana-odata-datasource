import cds from "@sap/cds";
import morgan from "morgan";
import {Application} from "express";
// @ts-ignore
import {addMockService} from "../mock/MockService";

cds.on('bootstrap', async (app: Application) => {
  app.use(morgan('dev'));
  await addMockService(app);
})
