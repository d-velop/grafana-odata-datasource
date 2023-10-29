import cds from "@sap/cds";
import morgan from "morgan";
import flash from 'connect-flash';
import {Application} from "express";
// @ts-ignore
import {addMockService} from "../mock/MockService";

cds.on('bootstrap', async (app: Application) => {
  app.use(morgan('dev'));
  app.use(flash());
  await addMockService(app);
})
