import {Application, Request, Response} from "express";
import expressSession from "express-session";
import {Issuer, Strategy} from 'openid-client';
import passport from 'passport';
import {ToXml} from "./util/XmlMetadataWriter";
import {MockModel} from "./model/MockModel";
import {GenerateTemperatures} from "./Testdata";

export let addMockService = async (app: Application) => {
  const memoryStore = new expressSession.MemoryStore();

  app.use(expressSession({
    secret: 'another_long_secret',
    resave: false,
    saveUninitialized: true,
    store: memoryStore
  }));

  app.use(passport.initialize());
  app.use(passport.authenticate('session'));

  const keycloakIssuer = await Issuer.discover('http://localhost:8080/realms/grafana')
  console.log('Discovered issuer %s %O', keycloakIssuer.issuer, keycloakIssuer.metadata);

  const client = new keycloakIssuer.Client({
    client_id: 'test-server',
    client_secret: '1fakeTestServerSecret00000000000',
    redirect_uris: ['http://localhost:4004/auth/callback'],
    post_logout_redirect_uris: ['http://localhost:4004/logout/callback'],
    response_types: ['code'],
  });

  let TokenSet: any
  passport.use('oidc', new Strategy({client},
    (tokenSet: any, userinfo: any, done: any) => {
      TokenSet = tokenSet;
      console.log('Token Set:', tokenSet);
      console.log('User Info:', userinfo);
      return done(null, tokenSet.claims());
    })
  )

  passport.serializeUser(function (user: Express.User, done: any) {
    done(null, user);
  });

  passport.deserializeUser(function (user: any, done: any) {
    done(null, user);
  });

  app.get('/test', (req, res, next) => {
    passport.authenticate('oidc')(req, res, next);
  });

  app.get('/auth/callback', (req, res, next) => {
    console.log('/auth/callback req.url:', req.url);
    passport.authenticate('oidc', {
      successRedirect: '/testauth',
      failureRedirect: '/failure',
      failureFlash: true
    })(req, res, next);
  });

  const checkAuthenticated = (req: any, res: any, next: any) => {
    // req.isAuthenticated is populated by password.js
    if (req.isAuthenticated()) {
      return next()
    }
    res.redirect("/test")
  }

  app.get('/testauth', checkAuthenticated, (req, res) => {
    res.send('You are connected (testauth)');
  });

  app.get('/logout', (req, res) => {
    res.redirect(client.endSessionUrl({
      id_token_hint: TokenSet.id_token
    }));
  });

  app.get('/failure', (req, res) => {
    const errorMessage = req.flash('error');
    res.send(`Failure! Error Message: ${errorMessage}`);
  });

  app.get('/logout/callback', (req, res, next) => {
    req.logout((err) => {
      if (err) {
        return next(err)
      }
      res.redirect('/mock');
    });
  });

  app.get('/mock/\\$metadata',
    async (_: Request, res: Response): Promise<Response> => {
      return res
        .contentType('application/xml')
        .status(200)
        .send(ToXml(MockModel));
    }
  );

  app.get('/mock/temperatures',
    async (req: Request, res: Response): Promise<Response> => {
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
      return res
        .contentType('application/json')
        .status(200).send(
          {
            '@odata.context': 'http://localhost:4004/odata/$metadata#Temperatures',
            value: values
          });
    }
  );

  app.get('/mock',
    async (_: Request, res: Response): Promise<Response> => {
      let entitySets = [];
      for (let schema of MockModel.schemas) {
        let entityContainer = schema.entityContainer;
        if (entityContainer != null) {
          for (let entitySet of entityContainer.entitySets) {
            entitySets.push({name: entitySet.name, kind: 'EntitySet', url: entitySet.name});
          }
        }
      }
      return res
        .contentType('application/json')
        .status(200)
        .send({
          '@odata.context': 'http://localhost:4004/odata/$metadata',
          value: entitySets
        });
    }
  );
}
