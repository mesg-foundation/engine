FROM node:10.6.0-alpine
WORKDIR /app
COPY ./package* ./
RUN npm install
COPY . .
CMD [ "node", "index.js" ]