import {Application, NextFunction, Request, Response} from 'express';
import expressSession from 'express-session';
import * as client from 'openid-client';
import {createRemoteJWKSet, jwtVerify} from 'jose';
import {ToXml} from './util/XmlMetadataWriter';
import {MockModel} from './model/MockModel';
import {GenerateTemperatures} from './Testdata';

export let addMockService = async (app: Application) => {
  const keycloakAddr = process.env.KEYCLOAK_ADDR;
  const authEnabled = !!keycloakAddr;
  let ensureAuth: (req: Request, res: Response, next: NextFunction) => Promise<void>;

  if (authEnabled) {
    const issuer = `http://${keycloakAddr}/realms/grafana`;
    const jwksUri = `${issuer}/protocol/openid-connect/certs`;
    const jwks = createRemoteJWKSet(new URL(jwksUri));

    app.use(expressSession({
      secret: 'another_long_secret',
      resave: false,
      saveUninitialized: false,
      cookie: {
        httpOnly: true,
        sameSite: 'lax',
        secure: false
      }
    }));

    const config = await client.discovery(
      new URL(issuer),
      'test-server',
      '1fakeTestServerSecret00000000000',
      null,
      {
        execute: [client.allowInsecureRequests],
      }
    );
    console.log('Discovered issuer %O', config.serverMetadata());

    app.get('/login',
      async (req: Request, res: Response): Promise<void> => {
        const returnTo = req.query.returnTo as string | undefined;
        const codeVerifier = client.randomPKCECodeVerifier();
        const codeChallenge = await client.calculatePKCECodeChallenge(codeVerifier);
        const state = client.randomState();
        req.session['pkceVerifier'] = codeVerifier;
        req.session['state'] = state;
        if (returnTo?.startsWith('/')) {
          req.session['returnTo'] = returnTo;
        }
        const redirectUrl = client.buildAuthorizationUrl(config, {
          redirect_uri: `${req.protocol}://${req.get('host')}/auth/callback`,
          scope: 'openid email profile offline_access roles',
          code_challenge: codeChallenge,
          code_challenge_method: 'S256',
          state
        });
        res.redirect(redirectUrl.href);
      }
    );

    app.get('/auth/callback',
      async (req: Request, res: Response): Promise<void> => {
        const url = new URL(`${req.protocol}://${req.get('host')}${req.originalUrl}`);
        const tokens = await client.authorizationCodeGrant(
          config,
          url,
          {
            pkceCodeVerifier: req.session['pkceVerifier']!,
            expectedState: req.session['state']!,
          }
        );
        const claims = tokens.claims();
        const user = await client.fetchUserInfo(config, tokens.access_token, claims.sub);
        req.session['tokens'] = tokens;
        req.session['user'] = user;
        res.redirect(req.session['returnTo']);
      }
    );

    app.get('/logout',
      async (req: Request, res: Response): Promise<void> => {
        const idToken = req.session['tokens']?.id_token;
        req.session.destroy(() => {
          if (idToken) {
            const endSessionUrl = client.buildEndSessionUrl(config, { id_token_hint: idToken });
            res.redirect(endSessionUrl.href);
          }
        });
      }
    );

    ensureAuth = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
      let rawAuthHeader = req.headers.authorization;
      if (rawAuthHeader) {
        let authHeader = Array.isArray(rawAuthHeader) ? rawAuthHeader[0] : rawAuthHeader;
        let jwt = authHeader.replace(/^Bearer /i, '');
        const { payload } = await jwtVerify(jwt, jwks, { issuer });
        console.log('user: ', { preferred_username: payload.preferred_username, email: payload.email });
        return next();
      }
      if (req.session?.['user']) return next();
      const nextUrl = encodeURIComponent(req.originalUrl);
      res.redirect(`/login?returnTo=${nextUrl}`);
    };
  } else {
    console.log('Keycloak authentication disabled - KEYCLOAK_ADDR not set');
    ensureAuth = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
      next();
    };
  }

  app.get('/mock/testauth', ensureAuth,
    async (req: Request, res: Response): Promise<void> => {
      if (!authEnabled) {
        res.send('Auth is disabled - no user information available');
      } else {
        const user = req.session['user'] as any;
        res.send(`Logged in as user: ${user.preferred_username} (${user.email})`);
      }
    }
  );

  app.get('/mock/\\$metadata', ensureAuth,
    async (res: Response): Promise<void> => {
      res
        .contentType('application/xml')
        .status(200)
        .send(ToXml(MockModel));
    }
  );

  app.get('/mock/temperatures', ensureAuth,
    async (req: Request, res: Response): Promise<void> => {
      const { $filter } = req.query;
      if (typeof $filter === 'string') {
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
        .status(200).send({
          '@odata.context': new URL(`${req.protocol}://${req.get('host')}/odata/$metadata#Temperatures`),
          value: values
        });
    }
  );

  app.get('/mock', ensureAuth,
    async (req: Request, res: Response): Promise<void> => {
      let entitySets = [];
      for (let schema of MockModel.schemas) {
        let entityContainer = schema.entityContainer;
        if (entityContainer != null) {
          for (let entitySet of entityContainer.entitySets) {
            entitySets.push({ name: entitySet.name, kind: 'EntitySet', url: entitySet.name });
          }
        }
      }
      res
        .contentType('application/json')
        .status(200)
        .send({
          '@odata.context': new URL(`${req.protocol}://${req.get('host')}/odata/$metadata`),
          value: entitySets
        });
    }
  );
}
