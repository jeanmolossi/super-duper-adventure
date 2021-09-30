FROM node:14.16.1-alpine

WORKDIR /api

COPY . .

RUN yarn

CMD [ "yarn","nodemon", "index" ]