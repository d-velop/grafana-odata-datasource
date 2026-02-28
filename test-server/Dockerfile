FROM node:24.14.0-bookworm

WORKDIR /app

COPY package.json pnpm-lock.yaml .npmrc generateTestdata.ts tsconfig.json /app/
COPY db /app/db
COPY srv /app/srv
COPY mock /app/mock

RUN corepack enable && pnpm install --frozen-lockfile

EXPOSE 4004/tcp

CMD [ "pnpm", "start" ]
