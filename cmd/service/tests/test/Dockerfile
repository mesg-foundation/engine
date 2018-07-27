FROM node:10.6.0-jessie
WORKDIR /app
COPY ./package* ./
RUN npm install
COPY . .
CMD [ "node", "index.js" ]