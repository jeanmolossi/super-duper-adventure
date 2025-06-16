FROM node:14.21.3-alpine

WORKDIR /api

COPY . .

RUN yarn

CMD [ "yarn","nodemon", "index" ]