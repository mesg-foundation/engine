FROM node:8

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY package.json ./
COPY package-lock.json ./
RUN npm install --production

COPY . ./

CMD [ "npm", "start" ]
