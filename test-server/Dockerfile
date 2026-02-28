FROM node:24.14.0-bookworm

WORKDIR /app

COPY package.json yarn.lock generateTestdata.ts tsconfig.json /app/
COPY db /app/db
COPY srv /app/srv
COPY mock /app/mock

RUN yarn install

EXPOSE 4004/tcp

CMD [ "yarn", "start" ]
